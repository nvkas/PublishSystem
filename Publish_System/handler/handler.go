package handler

import "context"

type ReceiveHookService struct {
	HookService
}

type HookService interface {
	ReceiveRequest(ctx context.Context, args *interface{}, result *interface{}) error
}


//空实现,因为注册微服务需要服务的方法.
//可用于扩展功能
func(self *ReceiveHookService)ReceiveRequest(ctx context.Context, args *interface{}, result *interface{}) error {
	return nil
}
