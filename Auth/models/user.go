package models

type User struct {
	ID       int     `json:"id"`
	FullName string  `json:"fullname"`
	Password string  `json:"password"`
	Mail     string  `json:"mail"`
	PP       *string `json:"profilimage"`
	PPURL    *string `json:"profilimageurl"`
}

type LogIn struct {
	Password string `json:"password"`
	Mail     string `json:"mail"`
}

type UpdatePassword struct {
	OldPass  string `json:"oldpass"`
	NewPass1 string `json:"newpass1"`
	NewPass2 string `json:"newpass2"`
}
