package main

import (
	"fmt"
	"log"

	go_salesforce_client "github.com/MASA-JAPAN/go-salesforce-client"
)

func main() {
	AuthenticatePasswordExample()
	AuthenticateClientCredentialsExample()
}

func AuthenticatePasswordExample() {
	auth := go_salesforce_client.Auth{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		Username:     "your_username",
		Password:     "your_password",
		TokenURL:     "https://login.salesforce.com/services/oauth2/token",
	}

	client, err := auth.AuthenticatePassword()
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	fmt.Printf("Authenticated! Access Token: %s\n", client.AccessToken)
}

func AuthenticateClientCredentialsExample() {
	auth := go_salesforce_client.Auth{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		TokenURL:     "https://yourdomain/services/oauth2/token",
	}

	client, err := auth.AuthenticateClientCredentials()
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	fmt.Printf("Authenticated! Access Token: %s\n", client.AccessToken)
}
