//go:build windows

/*
與Bind實作等相關函數，都在這個檔案之中
*/

package webview2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	rpc2 "github.com/CarsonSlovoka/go-webview2/v1/pkg/rpc"
	"log"
	"reflect"
	"strconv"
)

// AddDispatch 會把函數發送到消息循環之中等待執行。等待執行的函數只會執行一次之後就會從等待隊列中移除，即下次觸發WM_APP將不再執行
// 此函數被messageCallback所應用
func (w *webView) AddDispatch(f func()) {
	w.m.Lock()
	w.dispatchQueue = append(w.dispatchQueue, f)
	w.m.Unlock()
	_ = dll.User.PostThreadMessage(w.threadID, w32.WM_APP, 0, 0)
}

func (w *webView) AddScriptOnDocumentCreated(script string) {
	w.Browser.AddScriptOnDocumentCreated(script)
}

// SetBind 提供javascript觸發window.chrome.webview.postMessage的途徑，並且觸發後能執行所設定的函數
// 過程:
// 主要是透過AddScriptOnDocumentCreated注入函數
// 此函數會向window新增一個name的成員，這個成員為一個函數，回傳型態為Promise
// 裡面會做兩件事情
// 1. 定義回傳Promise的內容: 這裡面定義了window._rpc[seq]的內容，所以可以透過其他的javascript來呼叫此內容，把resolve或者reject來完成
// 2. 觸發window.external.postMessage
// 當透過javascript主動呼叫此window.name時，就會開始執行此Promise，並執行以上兩點內容
// 其中第二點，其識別該執行哪一個我們在go定義的函數，主要是透過name來區分
// window.external.postMessage 此函數是在webview註冊完成時，也是透過AddScriptOnDocumentCreated所加入
// 其定義為: c.AddScriptOnDocumentCreated(`window.external={postMessage:jsonStr=>window.chrome.webview.postMessage(jsonStr)}`)
func (w *webView) SetBind(name string, f any) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return fmt.Errorf("'f' is not the function")
	}

	if n := v.Type().NumOut(); n > 2 {
		return fmt.Errorf("function may only return a value or (value, error)")
	}

	w.m.Lock()
	if _, exists := w.bindFuncMap[name]; exists {
		return fmt.Errorf("the function has bonded already")
	}
	w.bindFuncMap[name] = f
	w.m.Unlock()

	// 此嵌入的函數打開inspect還是有辦法debug，只要在window[name]下斷點，進入函數內部即可對此內容進行debug
	w.AddScriptOnDocumentCreated(fmt.Sprintf(`(() => {
  const name = %q
  if (window._rpc === undefined) {
    window._rpc = {nextSeq: 1}
  }
  window[name] = (arguments) => {
    const seq =  window._rpc.nextSeq++
    const promise = new Promise((resolve, reject) => {
      window._rpc[seq] = {
        resolve: resolve,
        reject: reject,
      }
    })
    window.external.postMessage(JSON.stringify({
      id: seq,
      method: name,
      params: Array.prototype.slice.call(arguments),
    }))
    return promise
  }
})()`, name))
	// 其中params: Array.prototype.slice.call() 可以使參數都是傳入一個array。(當arguments為字串，傳入[<str>]，為array則不影響還是array
	return nil
}

func (w *webView) ExecuteScript(javascript string) {
	w.Browser.ExecuteScript(javascript)
}

// WebMessageReceived時會需要使用到此函數
func (w *webView) messageCallback(msg string) {
	request := rpc2.Request{}
	if err := json.Unmarshal([]byte(msg), &request); err != nil {
		log.Printf("invalid RPC message: %v", err)
		return
	}

	id := strconv.Itoa(request.ID)
	var err error
	defer func() {
		if err != nil {
			jsBytes, _ := json.Marshal(err.Error())
			// window._rpc是我們在SetBind的時候所定義的內容，如果有發生錯誤，我們呼叫他的Promise的reject函數，去觸發錯誤，並且再完成之後重新將此window._rpc[id]設為undefined
			w.ExecuteScript("window._rpc[" + id + "].reject(" + string(jsBytes) + "); window._rpc[" + id + "] = undefined")
		}
	}()

	var response any
	if response, err = w.callBinding(&request); err != nil {
		return
	}

	var b []byte
	if b, err = json.Marshal(response); err != nil { // 因為我們不曉得使用者定義的go函數的回傳值真正型別是什麼，所以統一轉換成json字串
		return
	}

	// 如果轉換成功
	w.AddDispatch(func() { // 如果字串轉換成功，我們會把這個字串b當成參數，傳給resolve來完成
		// window._rpc[id]是我們在SetBind的時候所定義的內容，他是一個object，有resolve與reject的成員，該兩成員的內容都是一個函數，為某一個promise的resolve與reject函數，所以要透過呼叫此成員來完成該promise
		// 注意我們回傳的其實不是一個字串，因為並沒有用"把它括號起來，所以如果是物件，就會直接是物件，前端不需要在透過JSON.parse來轉成物件
		w.ExecuteScript("window._rpc[" + id + "].resolve(" + string(b) + "); window._rpc[" + id + "] = undefined") // 執行resolve方法，使得Promise完成，且再完成之後，我們將此id重設，其實能被再利用
	})
}

// 依據Request去向bindFuncMap，找到相對應的函數，執行後，將結果回傳
func (w *webView) callBinding(r *rpc2.Request) (any, error) {
	w.m.Lock()
	f, ok := w.bindFuncMap[r.Method] // 取得我們定義的go函數
	w.m.Unlock()
	if !ok {
		return nil, nil
	}

	v := reflect.ValueOf(f)
	isVariadic := v.Type().IsVariadic() // 是否存在...Type
	numIn := v.Type().NumIn()           // 看我們的function有幾個參數
	// 確認我們輸入的function函數要與參數做匹配，其中r.Params為javascript所傳入的參數
	// 例如:
	//	Variadic: func(a, b string, opt ...int) => 則params至少要>=2
	//  非Variadic => Params 要與參數個數相同
	if (isVariadic && len(r.Params) < numIn-1) || (!isVariadic && len(r.Params) != numIn) {
		return nil, errors.New("function arguments mismatch")
	}
	args := make([]reflect.Value, 0)
	for i := range r.Params {
		var arg reflect.Value
		if isVariadic && i >= numIn-1 { // numIn-1才是真的下標值(其從0開始)
			arg = reflect.New(v.Type().In(numIn - 1).Elem())
		} else {
			arg = reflect.New(v.Type().In(i))
		}
		if err := json.Unmarshal(r.Params[i], arg.Interface()); err != nil {
			return nil, err
		}
		args = append(args, arg.Elem())
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	res := v.Call(args) // 執行bindFuncMap的該函數
	switch len(res) {
	case 0:
		// No results from the function, just return nil
		return nil, nil
	case 1:
		// One result may be a value, or an error
		if res[0].Type().Implements(errorType) { // 利用此方法可以得知是否回傳的型別為error
			if err := res[0].Interface(); err != nil { // 確認該錯誤是否為nil
				return nil, err.(error)
			}
			return nil, nil
		}
		return res[0].Interface(), nil
	case 2:
		// Two results: first one is value, second is error
		if !res[1].Type().Implements(errorType) {
			return nil, errors.New("second return value must be an error")
		}
		if res[1].Interface() == nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)

	default:
		return nil, errors.New(fmt.Sprintf("unexpected number of return values. number of return can be {0, 1, 2} not %d", len(res)))
	}
}
