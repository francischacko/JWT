package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func Loadvariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
