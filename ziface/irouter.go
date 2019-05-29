package ziface

type IRouter interface {
	PreHandle(request IRequest, response IResponse)

	Handle(request IRequest, response IResponse)

	PostHandle(request IRequest, response IResponse)
}
