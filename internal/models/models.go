package models

type Response struct {
	Status  int         `json:"Status"`
	Payload interface{} `json:"Payload"`
}

type UserToken struct {
	UserId       string
	RefreshToken string
	AccessToken  string
}
