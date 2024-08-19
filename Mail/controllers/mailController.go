package controllers

import (
	"log"
	"mail/config"
	"mail/database"
	"mail/helpers"
	"mail/models"
	"mail/utils"

	"github.com/gofiber/fiber/v2"
)

func SendNofication(c *fiber.Ctx) error {
	users, err := helpers.GetUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-N-1"})
	}
	var sendNoficaion models.SendNoficitaion
	err = c.BodyParser(&sendNoficaion)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-N-2"})
	}
	writed := sendNoficaion.Text
	for _, user := range users {
		nofication := models.SendNoficitaion{
			SendingPerson: user.Mail,
			Text:          writed,
			OurMail:       "gokhandalmzz@gmail.com",
		}
		err := config.RabbitMqPublish([]byte(nofication.Text), nofication.SendingPerson)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-N-3"})
		}
		err = config.RabbitMqConsume(nofication.SendingPerson, nofication.OurMail)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "S-N-4"})
		}
		err = database.DB.Db.Create(&nofication).Error
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "S-N-5"})
		}
		err = utils.WriteNotificationToFile(nofication)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-N-6"})
		}
		log.Println(nofication)
	}
	var nofication []models.SendNoficitaion
	return c.Status(200).JSON(fiber.Map{"status": "Success", "Message": "Success", "data": nofication})
}

func TakeNotification(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "T-N-1"})
	}
	var notification []models.SendNoficitaion
	err := database.DB.Db.Where("sending_person=?", user.Mail).First(&notification).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR : T-N-2", "details": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success", "data": notification})
}
