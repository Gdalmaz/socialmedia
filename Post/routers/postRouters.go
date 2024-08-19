package routers

import (
	"post/controllers"
	"post/middleware"

	"github.com/gofiber/fiber/v2"
)

func PostRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	post := v1.Group("/post")

	post.Post("/create-post", middleware.TokenControl(), controllers.CreatePost)
	post.Put("/update-post/:id", middleware.TokenControl(), controllers.UpdatePost)
	post.Delete("/delete-post/:id", middleware.TokenControl(), controllers.DeletePost)
	post.Post("/like-post/:id", middleware.TokenControl(), controllers.LikePost)
	post.Get("/get-posts", controllers.GetAllPost)
	post.Get("/get-user-post", middleware.TokenControl(), controllers.GetUserPost)
}
