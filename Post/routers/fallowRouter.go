package routers

import (
	"post/controllers"
	"post/middleware"

	"github.com/gofiber/fiber/v2"
)

func FollowRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	follow := v1.Group("/follow")

	follow.Post("/follow", middleware.TokenControl(), controllers.Follow)
	follow.Delete("/unfollow", middleware.TokenControl(), controllers.UnFollow)
	follow.Get("/get-follower", middleware.TokenControl(), controllers.GetAllFollower)
	follow.Get("/get-following", middleware.TokenControl(), controllers.GetAllFollowing)
}
