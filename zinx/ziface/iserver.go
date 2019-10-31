package ziface

//定义一个服务器接口
type IServer interface {
	//	启动服务器
	Start()
	//	停止服务器
	Stop()
	//运行服务器
	Server()
	//	添加路由
	AddRouter(msgID uint32, router IRouter)
	//  获取链接管理器
	GetConnManager() IConnManager
	//注册OnConnStart钩子函数的方法
	SetOnConnStart(func(conn IConnection))
	//注册OnConnStop钩子函数的方法
	SetOnConnStop(func(conn IConnection))
	//调用OnConnStart钩子函数的方法
	CallOnConnStart(conn IConnection)
	//调用OnConnStart钩子函数的方法
	CallOnConnStop(conn IConnection)
}
