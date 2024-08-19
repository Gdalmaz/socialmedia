package models

type Post struct {
	ID           int    `json:"id"`
	User         User   `gorm:"foreignKey:UserID"`
	UserID       int    `json:"userid"`
	PostPhoto    string `json:"postphoto"`
	PostPhotoURL string `json:"postphotourl"`
	PostText     string `json:"posttext"`
	LikeCount    int    `json:"likecount" gorm:"default:0"`
	CommentCount int    `json:"commentcount" gorm:"default:0"`
	IsActive     bool   `gorm:"default:true"`
}

type LikePost struct {
	Post   Post `json:"PostID"`
	User   User `json:"UserID"`
	UserID int  `json:"userid"`
	PostID int  `json:"postid"`
}
