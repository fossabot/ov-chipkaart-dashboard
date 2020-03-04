package main

import (
	"github.com/joho/godotenv"
	"log"
)

type AuthRequestPayload struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}



func getAccessToken(payload AuthRequestPayload, endpoint string) (accessToken string) {
	const grantTypePassword = "password"
	const scope = "openid"


	return accessToken
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
