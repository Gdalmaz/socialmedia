package controllers

import (
	"post/config"
	"post/database"
	"post/models"

	"github.com/gofiber/fiber/v2"
)

//Postları asla silmeyeceğiz sadece false ve true durumuna göre getireceğiz

func CreatePost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: C-P-1"})
	}
	var post models.Post
	post.UserID = user.ID
	file, err := c.FormFile("postphoto")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR: C-P-2"})
	}

	fileBytes, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR: C-P-3"})
	}
	defer fileBytes.Close()

	imageBytes := make([]byte, file.Size)
	_, err = fileBytes.Read(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR: C-P-4"})
	}

	id, url, err := config.CloudConnect(imageBytes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR: C-P-5"})
	}
	post.PostPhoto = id
	post.PostPhotoURL = url

	posttext := c.FormValue("posttext")
	if len(posttext) == 0 {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: C-P-6"})
	}

	err = database.DB.Db.Preload("User").Create(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: C-P-6"})
	}
	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success", "data": post})
}

func UpdatePost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: U-P-1"})
	}
	postID := c.Params("id")
	var post models.Post
	if err := database.DB.Db.First(&post, postID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: U-P-2"})
	}
	if post.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: U-P-3"})
	}
	posttext := c.FormValue("posttext")
	if len(posttext) != 0 {
		post.PostText = posttext
	}
	file, err := c.FormFile("postphoto")
	if err == nil {
		fileBytes, err := file.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR: U-P-4"})
		}
		defer fileBytes.Close()

		imageBytes := make([]byte, file.Size)
		_, err = fileBytes.Read(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR: U-P-5"})
		}

		id, url, err := config.CloudConnect(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "ERROR: U-P-6"})
		}
		post.PostPhoto = id
		post.PostPhotoURL = url
	}
	if err := database.DB.Db.Save(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: U-P-7"})
	}

	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success", "data": post})
}

func DeletePost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: D-P-1"})
	}
	var post models.Post

	if user.ID != post.ID {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: D-P-2"})
	}
	id := c.Params("id")
	err := database.DB.Db.Where("id=?", id).First(&post).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: D-P-3"})
	}
	post.IsActive = false

	err = database.DB.Db.Updates(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: D-P-4"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "Message": "Success"})
}

func LikePost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: L-P-1"})
	}
	var post models.Post
	var likePost models.LikePost

	id := c.Params("id")
	err := database.DB.Db.Where("id=?", id).Find(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: L-P-2"})
	}

	likePost.UserID = user.ID
	likePost.PostID = post.ID

	err = database.DB.Db.Create(&likePost).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: L-P-3"})
	}

	post.LikeCount += 1

	err = database.DB.Db.Save(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR: L-P-4"})
	}
	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success"})

}

func GetAllPost(c *fiber.Ctx) error {
	var posts []models.Post
	err := database.DB.Db.Preload("User").Where("is_active = ?", true).Order("id DESC").Find(&posts).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "G-A-P-1"})
	}
	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success", "data": posts})
}

func GetUserPost(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : G-U-P-1"})
	}

	var posts []models.Post
	err := database.DB.Db.Preload("User").Where("user_id = ? AND is_active = ?", user.ID, true).Find(&posts).Error
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : G-U-P-2"})
	}

	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success", "data": posts})
}
