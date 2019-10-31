package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

//基于Zinx框架来开发的服务器端应用程序

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Test PreHandle
//func (this *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping....\n")); err != nil {
//		fmt.Println("call back before ping error: ", err)
//	}
//}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	//	读取客户端的数据，再回写
	data := request.GetData()
	msgId := request.GetMsgID()
	msgLen := request.GetMsgLen()
	fmt.Printf("recv from client: MsgID = %d, Len = %d, data = %s\n", msgId, msgLen, string(data))
	if err := request.GetConnection().SendMsg(1, []byte("ping...")); err != nil {
		fmt.Println(err)
	}
}

//Test PostHandle
//func (this *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandle...")
//	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping....\n")); err != nil {
//		fmt.Println("call back after ping error: ", err)
//	}
//}

func main() {
	//	创建一个server句柄使用Zinx的api
	s := znet.NewServer()
	//给当前zinx框架添加一个自定义router
	s.AddRouter(&PingRouter{})
	//启动server
	s.Server()
}
