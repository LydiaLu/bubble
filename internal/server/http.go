package server

import (
	v1 "bubble/api/bubble/v1"
	v2 "bubble/api/healthcheck/v1"
	"bubble/internal/conf"
	"bubble/internal/service"
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Middleware自定义中间件
// type middleware func(Handler) Handler
// type Handler func(ctx context.Context, req interface{}) (interface{}, error)
func Middleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		// 传handler，原来就要执行handler函数，现在可以在前后做一些其他步骤
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 执行之前做点事(比如JWT鉴权，必须有token，否则拒绝掉请求等等)
			fmt.Println("Middleware: 执行handler之前")
			// 做token校验
			if tr, ok := transport.FromServerContext(ctx); ok {
				token := tr.RequestHeader().Get("token")
				fmt.Printf("token: %v\n", token)
			}
			defer func() {
				fmt.Println("Middleware: 执行handler之后")
			}()
			return handler(ctx, req) //  执行目标handler
		}
	}
}

// Middleware1 自定义中间件
// type middleware func(Handler) Handler
// 闭包形式：定制化不同需求，更加通用，代码复用
// func Middleware1(opts ...string) middleware.Middleware {
// 	return func(handler middleware.Handler) middleware.Handler {
// 		//opts
// 		return nil
// 	}
// }

// Middleware2 自定义中间件2，相比Middleware1失去了灵活性
// 一般框架不这么做！调用时无法定制化需求
// func Middleware2(middleware.Handler) middleware.Handler{
// 	return nil
// }

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, todo *service.TodoService, healthcheck *service.HealthCheckService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(), //全局中间件
			validate.Validator(),
			// Middleware(),   //全局中间件
			// Middleware1("a"),
			// Middleware1("b"),
			selector.Server( //特定Path才执行的中间件
				Middleware(),
			).Path("/api.bubble.v1.Todo/GetTodo").
				Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterTodoHTTPServer(srv, todo)
	v2.RegisterHealthCheckServiceHTTPServer(srv, healthcheck)
	return srv
}
