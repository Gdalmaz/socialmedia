package models

type Follower struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	UserID     int `json:"folowing"` // Takip eden kullanıcı
	FollowerID int `json:"folowed"`  // Takip edilen kullanıcı
}

type Data struct {
	ID int `json:"id"`
}