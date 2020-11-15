package alert

import "net/http"

func Main() {

	//注册回调函数
	http.HandleFunc(PhoneAlertURL, PhoneAlertHandler)

	//绑定tcp监听地址，并开始接受请求，然后调用服务端处理程序来处理传入的连接请求。
	//params：第一个参数 addr 即监听地址；第二个参数表示服务端处理程序，通常为nil
	//当参2为nil时，服务端调用 http.DefaultServeMux 进行处理
	http.ListenAndServe("0.0.0.0:8080", nil)
}
