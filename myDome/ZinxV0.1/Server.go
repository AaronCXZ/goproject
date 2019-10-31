package main

import "zinx/znet"

//基于Zinx框架来开发的服务器端应用程序

func main() {
	//	创建一个server句柄使用Zinx的api
	s := znet.NewServer("Zinx V0_1")
	//启动server
	s.Server()
}
