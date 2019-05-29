package znet

import "zinx/ziface"

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request ziface.IRequest, response ziface.IResponse) {}

func (br *BaseRouter) Handle(request ziface.IRequest, response ziface.IResponse) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest, response ziface.IResponse) {}
