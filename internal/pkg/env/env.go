package env

import (
	"log"

	dotenv "github.com/joho/godotenv"
)

func Init() {
	if err := dotenv.Load("../.env"); err != nil {
		log.Fatalln("Cannot find .env file")
	}
}
