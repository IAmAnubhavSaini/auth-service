package config

import (
	"auth-service/helpers"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var JwtSecret []byte
var JwtExpiration int

func LoadJWTConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(JwtSecret) == 0 {
		fmt.Println("JWT_SECRET not set in .env file")
		return
	}

	JwtExpiration := helpers.StringToNumber(os.Getenv("JWT_EXPIRATION"))
	if JwtExpiration == 0 {
		JwtExpiration = 72
	}
}
