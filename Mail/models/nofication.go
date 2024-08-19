package models

type SendNoficitaion struct {
	ID            int    `json:"id"`
	SendingPerson string `json:"sendingperson"`
	Text          string `json:"text"`
	OurMail       string `json:"ourmail"`
}

