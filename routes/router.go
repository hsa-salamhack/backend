package routes

import "github.com/gofiber/fiber/v2"

type Route struct {
	Name   string
	Method string
	Run    fiber.Handler
}

var AllRoutes []Route

func Register(route Route) {
	AllRoutes = append(AllRoutes, route)
}
