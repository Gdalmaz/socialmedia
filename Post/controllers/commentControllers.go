package controllers

import (
	"post/config"
	"post/database"
	"post/models"

	"github.com/gofiber/fiber/v2"
)

func AddComment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-1"})
	}

	var post models.Post
	var comment models.Comment

	err := c.BodyParser(&comment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-2"})
	}

	id := c.Params("id")
	err = database.DB.Db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-3"})
	}

	comment.UserID = user.ID
	comment.PostID = post.ID

	commentText := c.FormValue("commenttext")
	if commentText == "" {
		return c.Status(400).JSON(fiber.Map{"Status": "Error", "Message": "Comment text cannot be empty"})
	}
	comment.CommentText = commentText

	file, err := c.FormFile("commentphoto")
	if err == nil {
		fileBytes, err := file.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-4"})
		}
		defer fileBytes.Close()

		imageBytes := make([]byte, file.Size)
		_, err = fileBytes.Read(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-5"})
		}

		id, url, err := config.CloudConnect(imageBytes)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-6"})
		}

		comment.CommentPhoto = &id
		comment.CommentPhotoURL = &url
	}

	err = database.DB.Db.Create(&comment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-7"})
	}

	post.CommentCount += 1
	err = database.DB.Db.Updates(&post).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-8"})
	}

	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Comment added successfully"})
}

func DeleteComment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(404).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : D-C-1"})
	}
	var comment models.Comment
	if user.ID != comment.UserID {
		return c.Status(404).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : D-C-2"})
	}
	id := c.Params("id")
	err := database.DB.Db.Where("id=?", id).First(&comment).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : D-C-3"})
	}
	comment.IsActive = false
	err = database.DB.Db.Updates(&comment).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : D-C-4"})
	}

	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success"})
}

func AnswerComment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-1"})
	}
	var comment models.Comment
	var answerComment models.AnswerComment
	err := c.BodyParser(&answerComment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-2"})
	}
	id := c.Params("id")
	err = database.DB.Db.Where("id=?", id).First(&comment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-3"})
	}
	answerComment.CommentID = comment.ID
	answerComment.UserID = user.ID
	if len(answerComment.Answer) == 0 {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "ERROR : A-C-4"})
	}
	err = database.DB.Db.Create(&answerComment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : A-C-5"})
	}

	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success"})

}

func GetComment(c *fiber.Ctx) error {
	_, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : G-C-1"})
	}
	postID := c.Params("id")
	var comments []models.Comment
	err := database.DB.Db.Preload("User").Where("post_id = ? AND is_active = ?", postID, true).Find(&comments).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Status": "Error", "Message": "G-C-2"})
	}
	return c.Status(200).JSON(fiber.Map{"Status": "Success", "Message": "Success", "data": comments})
}

func LikeComment(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : L-C-1"})
	}
	var comment models.Comment
	var likeComment models.LikeComment
	id := c.Params("id")
	// err := c.BodyParser(&likeComment)
	// if err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : L-C-2"})
	// }
	err := database.DB.Db.Where("id=?", id).First(&comment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : L-C-3"})
	}
	likeComment.UserID = user.ID
	likeComment.CommentID = comment.ID
	err = database.DB.Db.Create(&likeComment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : L-C-4"})
	}
	comment.LikeCount += 1
	err = database.DB.Db.Save(&comment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "Error", "Message": "ERROR : L-C-5"})
	}
	return c.Status(200).JSON(fiber.Map{"status": "Success", "Message": "Success"})

}
