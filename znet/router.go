package znet

import "zinx/ziface"

// 实现Router时，先嵌入这个基类，然后根据这个基类的方法重写
type BaseRouter struct{}

// 这里之所以BaseRouter的方法为空
// 是因为有的Router不希望有PreHandle和PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化
func (br *BaseRouter) PreHandle(req ziface.IRequest)  {}
func (br *BaseRouter) Handle(req ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
