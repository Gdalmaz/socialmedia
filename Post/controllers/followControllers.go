package controllers

import (
	"log"
	"post/database"
	"post/models"

	"github.com/gofiber/fiber/v2"
)

func Follow(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR : F-1"})
	}
	log.Println(user)
	var folow models.Follower
	var data models.Data
	c.BodyParser(data)
	//KENDİNİ TAKİP ETMESİN
	if user.ID == data.ID {
		return c.Status(400).JSON(fiber.Map{"status": "Error", "message": "ERROR : F-2"})
	}
	folow.UserID = user.ID
	folow.FollowerID = data.ID
	err := database.DB.Db.Save(&folow).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "message": "ERROR : F-3"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success"})
}

func UnFollow(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "Message": "ERROR : U-F-1"})
	}
	userID := user.ID
	var follow []models.Follower
	var data models.Data
	c.BodyParser(&data)
	followingId := data.ID
	err := database.DB.Db.Where("user_id =?", userID).Find(&follow).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "error": err})
	}
	for _, i := range follow {
		if i.FollowerID == followingId {
			err := database.DB.Db.Delete(&i).Error
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "Message": "ERROR : U-F-2"})
			}
			return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success"})
		}
	}
	return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "you allready unfollowing this user"})
}

func GetAllFollower(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR : G-A-F-1"})
	}
	follow := new([]models.Follower)
	err := database.DB.Db.Where("follower_id =?", user.ID).Find(&follow).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR : G-A-F-2"})
	}
	var users []models.User
	for _, user1 := range *follow {
		user := new(models.User)
		err := database.DB.Db.Where("id =?", user1.FollowerID).Find(&user).Error
		if err != nil {
			return err
		}
		users = append(users, *user)
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": users})
}

func GetAllFollowing(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR : G-A-F-1"})
	}
	follow := new([]models.Follower)
	err := database.DB.Db.Where("user_id=?", user.ID).Find(&follow).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR : G-A-F-2"})
	}
	var user2 []models.User
	for _, users := range *follow {
		user := new(models.User)
		err := database.DB.Db.Where("id=?", users.FollowerID).Find(&user).Error
		if err != nil {
			return err
		}

		user2 = append(user2, *user)
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user2})
}
