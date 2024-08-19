package helpers

import (
	"mail/database"
	"mail/models"
)

func GetUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
