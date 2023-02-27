package edge

// ICoreWebView2WebMessageReceivedEventArgsVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2webmessagereceivedeventargsvtbl
type ICoreWebView2WebMessageReceivedEventArgsVTbl struct {
	iUnknownVTbl

	// https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2webmessagereceivedeventargs?view=webview2-1.0.1518.46

	getSource                uintptr // The URI of the document that sent this web message.
	getWebMessageAsJson      uintptr // The message posted from the WebView content to the host converted to a JSON string.
	tryGetWebMessageAsString uintptr // If the message posted from the WebView content to the host is a string type, this method returns the value of that string.
}

type ICoreWebView2WebMessageReceivedEventArgs struct {
	vTbl *ICoreWebView2WebMessageReceivedEventArgsVTbl
}
