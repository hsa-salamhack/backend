package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hsa-salamhack/backend/routes"
	_ "github.com/hsa-salamhack/backend/routes"
)

func main() {
	app := fiber.New()

	for _, route := range routes.AllRoutes {
		switch route.Method {
		case "GET":
			app.Get(route.Name, route.Run)
		case "POST":
			app.Post(route.Name, route.Run)
		case "PUT":
			app.Put(route.Name, route.Run)
		case "DELETE":
			app.Delete(route.Name, route.Run)
		default:
			log.Printf("❌ %s %s\n", route.Method, route.Name)
		}

		fmt.Printf("✅ %s %s\n", route.Method, route.Name)
	}

	app.Listen(":3000")
}
