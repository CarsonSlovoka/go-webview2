# Example

本項目提供兩個例子

1. [ExampleHelloWorld](https://github.com/CarsonSlovoka/go-webview2/blob/26f92e09b0347b6452f5ef9f986d17badc289d53/examples/helloWorld/main.go#L53-L67)
2. [ExampleWithNotifyIcon](https://github.com/CarsonSlovoka/go-webview2/blob/26f92e09b0347b6452f5ef9f986d17badc289d53/examples/helloWorld/main.go#L69-L194)
   - 此範例會運用到[notifyIcon.go](./notifyicon.go)

而其中[server.go](./server.go)只是一個在本機運行的簡單server，模擬使用者自定義的頁面

> 注意！ 在使用的時候都要先進行[webviewloader.Install("要安裝到哪一個資料夾去", 是否允許使用本機的webviewloader.dll)](https://github.com/CarsonSlovoka/go-webview2/blob/26f92e09b0347b6452f5ef9f986d17badc289d53/examples/helloWorld/main.go#L25-L27)
>
> 目的是用來安裝[webviewloader.dll](../../v1/webviewloader/sdk/)
