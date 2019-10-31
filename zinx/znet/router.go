package znet

import (
	"zinx/ziface"
)

//实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类方法进行重写
type BaseRouter struct {
}

func (r *BaseRouter) PreHandle(request ziface.IRequest) {
}

func (r *BaseRouter) Handle(request ziface.IRequest) {
}

func (r *BaseRouter) PostHandle(request ziface.IRequest) {
}
