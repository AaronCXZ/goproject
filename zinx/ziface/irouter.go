package ziface

/*
路由抽象接口：
理由里的数据都是IRequest
*/

type IRouter interface {
	//	在处理业务之前的方法Hook
	PreHandle(request IRequest)
	//  处理业务的主方法Hook
	Handle(request IRequest)
	//  助理业务之后的方法Hook
	PostHandle(request IRequest)
}
