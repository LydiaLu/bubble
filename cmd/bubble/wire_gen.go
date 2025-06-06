// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"bubble/internal/biz"
	"bubble/internal/conf"
	"bubble/internal/data"
	"bubble/internal/server"
	"bubble/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	db, err := data.NewDB(confData)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(db, logger)
	if err != nil {
		return nil, nil, err
	}
	todoRepo := data.NewTodoRepo(dataData, logger)
	todoUsecase := biz.NewTodoUsecase(todoRepo, logger)
	todoService := service.NewTodoService(todoUsecase)
	healthCheckService := service.NewHealthCheckService()
	grpcServer := server.NewGRPCServer(confServer, todoService, healthCheckService, logger)
	httpServer := server.NewHTTPServer(confServer, todoService, healthCheckService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
