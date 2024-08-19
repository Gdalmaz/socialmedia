package routers

import (
	"post/controllers"
	"post/middleware"

	"github.com/gofiber/fiber/v2"
)

func CommentRouters(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	comment := v1.Group("/comment")

	comment.Post("/add-comment/:id", middleware.TokenControl(), controllers.AddComment)
	comment.Delete("/delete-comment", middleware.TokenControl(), controllers.DeleteComment)
	comment.Post("/answer-comment/:id", middleware.TokenControl(), controllers.AnswerComment)
	comment.Post("/get-comment", middleware.TokenControl(), controllers.GetComment)
	comment.Post("/like-comment/:id", middleware.TokenControl(), controllers.LikeComment)
}
