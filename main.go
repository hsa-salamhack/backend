package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"              // Import Swagger middleware
	_ "github.com/hsa-salamhack/backend/docs" // Import generated docs
	"github.com/hsa-salamhack/backend/routes"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title Wiki? API
// @version 0.5
// @description Wiki? API
// @host 9.141.41.77:8080
// @BasePath /
func main() {
	app := fiber.New()
app.Use(cors.New())

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

	app.Listen(":8080")
}
