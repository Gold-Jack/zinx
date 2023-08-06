package ziface

type IRouter interface {
	PreHandle(request IRequest)  // 处理conn业务之前的钩子方法
	Handle(request IRequest)     // 处理conn业务
	PostHandle(request IRequest) // 处理conn业务之后的钩子方法
}
