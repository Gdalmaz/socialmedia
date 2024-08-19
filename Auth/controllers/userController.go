package controllers

import (
	"auth/config"
	"auth/database"
	"auth/helpers"
	"auth/middleware"
	"auth/models"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {
	var user models.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-I-1"})
	}

	if len(user.Mail) < 10 {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-I-2"})
	}

	if len(user.Password) < 7 {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-I-3"})
	}
	//En küçük haneli isim örneği Can Ak
	if len(user.FullName) < 5 {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-I-4"})
	}
	mailControl, _ := helpers.MailControl(user.Mail)
	if mailControl == true {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "S-U-2"})
	}

	user.Password = helpers.HashPass(user.Password)
	err = database.DB.Db.Create(&user).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-I-5"})
	}
	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success", "data": user})

}

func SignIn(c *fiber.Ctx) error {
	var user models.User
	var successUser models.LogIn
	var session models.Session
	err := c.BodyParser(&successUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : S-I-1"})
	}
	successUser.Password = helpers.HashPass(successUser.Password)
	err = database.DB.Db.Where("mail=? and password=?", successUser.Mail, successUser.Password).First(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : S-I-2"})
	}
	token, err := middleware.GenerateToken(user.Mail)
	session.UserID = user.ID
	session.Token = token
	err = database.DB.Db.Create(&session).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : S-I-3"})
	}
	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success"})

}

func UpdatePass(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : U-P-1"})
	}
	var changePass models.UpdatePassword
	err := c.BodyParser(&changePass)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : U-P-2"})
	}

	changePass.OldPass = helpers.HashPass(changePass.OldPass)

	if changePass.OldPass != user.Password {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : U-P-3"})
	}
	if len(changePass.NewPass1) < 7 {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : U-P-4"})
	}
	if changePass.NewPass1 != changePass.NewPass2 {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : U-P-5"})
	}
	changePass.NewPass1 = helpers.HashPass(changePass.NewPass1)
	err = database.DB.Db.Updates(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : U-P-6"})
	}
	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success"})

}

func UpdateAccount(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : U-A-1"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-2"})
	}

	fileBytes, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-3"})
	}
	defer fileBytes.Close()

	imageBytes := make([]byte, file.Size)
	_, err = fileBytes.Read(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-4"})
	}

	id, url, err := config.CloudConnect(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-5"})
	}

	user.PP = &id
	user.PPURL = &url

	fullname := c.FormValue("fullname")
	mail := c.FormValue("mail")

	if len(fullname) > 5 {
		user.FullName = fullname
	}

	if len(mail) > 10 {
		user.Mail = mail
	}

	err = database.DB.Db.Updates(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "U-A-6"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": "Success"})

}

func DeleteAccount(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : D-A-1"})
	}
	id := c.Params("id")
	err := database.DB.Db.Where("id=?", id).First(&user).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : D-A-2"})
	}
	err = database.DB.Db.Delete(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : D-A-3"})
	}
	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success"})
}

func LogOut(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : L-O-1"})
	}
	var session models.Session
	err := database.DB.Db.Raw("DELETE FROM sessions WHERE user_id= ?", user.ID).Scan(&session).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "message": "ERROR : L-O-2"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Success"})

}
