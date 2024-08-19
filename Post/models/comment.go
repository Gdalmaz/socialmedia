package models

type Comment struct {
	ID              int    `json"id"`
	Post            Post   `json:"foreignKey:PostId"`
	User            User   `json:"foreignKey:UserID"`
	PostID          int    `json:"postid"`
	UserID          int    `json:"userid"`
	CommentText     string `json:"commenttext"`
	LikeCount       int    `gorm:"default:0"`
	IsActive        bool   `gorm:"default:true"`
	CommentPhoto    *string `json:"commentphoto"`
	CommentPhotoURL *string `json:"commentphotourl"`
}

type LikeComment struct {
	User      User    `gorm:"foreigenKey:UserID"`
	Comment   Comment `json:"foreignKey:CommentID"`
	CommentID int     `json:"commentid"`
	UserID    int     `json:"userid"`
}

type AnswerComment struct {
	User      User    `gorm:"foreigenKey:UserID"`
	Comment   Comment `gorm:"foreigenKey:CommentID"`
	CommentID int     `json:"commentid"`
	UserID    int     `json:"userid"`
	Answer    string  `json:"answer"`
}
