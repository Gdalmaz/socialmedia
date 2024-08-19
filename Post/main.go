package main

import (
	"post/database"
	"post/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	routers.CommentRouters(app)
	routers.FollowRouter(app)
	routers.PostRouter(app)
	app.Listen(":9091")
}
