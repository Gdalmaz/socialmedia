package routers

import (
	"mail/controllers"
	"mail/middleware"

	"github.com/gofiber/fiber/v2"
)

func MailRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	mail := v1.Group("/mail")

	//Admin için
	mail.Post("/send-nofication", controllers.SendNofication)
	//User için
	mail.Get("/take-nofication", middleware.TokenControl(), controllers.TakeNotification)
}
