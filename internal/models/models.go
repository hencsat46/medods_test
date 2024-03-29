package models

type Response struct {
	Status  int         `json:"Status"`
	Payload interface{} `json:"Payload"`
}
