//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"net/http"
	"sirkelin/backend/app/auth/controller"
	"sirkelin/backend/app/auth/repository"
	"sirkelin/backend/app/auth/service"
	"sirkelin/backend/router"
)

var authSet = wire.NewSet(
	repository.NewAuthRepository,
	wire.Bind(new(repository.IAuthRepository), new(*repository.AuthRepository)),
	service.NewAuthService,
	wire.Bind(new(service.IAuthService), new(*service.AuthService)),
	controller.NewAuthController,
	wire.Bind(new(controller.IAuthController), new(*controller.AuthController)),
)

func CreateHTTPServer() *http.Server {
	panic(wire.Build(
		NewServer,
		router.NewRouter,
		authSet,
	))
	return nil
}
