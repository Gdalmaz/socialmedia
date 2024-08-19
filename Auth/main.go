package main

import (
	"auth/database"
	"auth/middleware"
	"auth/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	routers.UserRoute(app)
	middleware.InitMetrics(app)
	app.Listen(":9090")
}
