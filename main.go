package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"              // Import Swagger middleware
	_ "github.com/hsa-salamhack/backend/docs" // Import generated docs
	"github.com/hsa-salamhack/backend/routes"
)

// @title Wiki? API
// @version 0.5
// @description Wiki? API
// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

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
			fmt.Printf("❌ %s %s\n", route.Method, route.Name)
		}

		fmt.Printf("✅ %s %s\n", route.Method, route.Name)
	}

	app.Listen(":3000")
}
