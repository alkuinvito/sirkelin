//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"net/http"
	authController "sirkelin/backend/app/auth/controller"
	authRepository "sirkelin/backend/app/auth/repository"
	authService "sirkelin/backend/app/auth/service"
	"sirkelin/backend/middlewares"

	roomController "sirkelin/backend/app/room/controller"
	roomRepository "sirkelin/backend/app/room/repository"
	roomService "sirkelin/backend/app/room/service"
	"sirkelin/backend/initializers"
	"sirkelin/backend/router"
)

var authSet = wire.NewSet(
	authRepository.NewAuthRepository,
	wire.Bind(new(authRepository.IAuthRepository), new(*authRepository.AuthRepository)),
	authService.NewAuthService,
	wire.Bind(new(authService.IAuthService), new(*authService.AuthService)),
	authController.NewAuthController,
	wire.Bind(new(authController.IAuthController), new(*authController.AuthController)),
)

var roomSet = wire.NewSet(
	roomRepository.NewRoomRepository,
	wire.Bind(new(roomRepository.IRoomRepository), new(*roomRepository.RoomRepository)),
	roomService.NewRoomService,
	wire.Bind(new(roomService.IRoomService), new(*roomService.RoomService)),
	roomController.NewRoomController,
	wire.Bind(new(roomController.IRoomController), new(*roomController.RoomController)),
)

func CreateHTTPServer() *http.Server {
	panic(wire.Build(
		initializers.NewDB,
		NewServer,
		middlewares.NewMiddleware,
		router.NewRouter,
		authSet,
		roomSet,
	))
	return nil
}
