//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"net/http"
	userController "sirkelin/backend/app/user/controller"
	userRepository "sirkelin/backend/app/user/repository"
	userService "sirkelin/backend/app/user/service"
	"sirkelin/backend/middlewares"

	roomController "sirkelin/backend/app/room/controller"
	roomRepository "sirkelin/backend/app/room/repository"
	roomService "sirkelin/backend/app/room/service"
	"sirkelin/backend/initializers"
	"sirkelin/backend/router"
)

var userSet = wire.NewSet(
	userRepository.NewUserRepository,
	wire.Bind(new(userRepository.IUserRepository), new(*userRepository.UserRepository)),
	userService.NewUserService,
	wire.Bind(new(userService.IUserService), new(*userService.UserService)),
	userController.NewUserController,
	wire.Bind(new(userController.IUserController), new(*userController.UserController)),
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
		userSet,
		roomSet,
	))
	return nil
}
