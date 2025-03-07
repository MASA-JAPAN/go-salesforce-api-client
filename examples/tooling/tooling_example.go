package main

import (
	"fmt"
	"log"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func main() {
	// Initialize authentication details
	auth := go_salesforce_api_client.Auth{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		TokenURL:     "https://yourdomain/services/oauth2/token",
	}

	// Authenticate and retrieve an access token
	client, err := auth.AuthenticateClientCredentials()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Query the Tooling API
	soql := "SELECT Id, Name, ApiVersion, Status FROM ApexClass"
	queryResponse, err := client.QueryToolingAPI(soql)
	if err != nil {
		fmt.Println("Error executing tooling query:", err)
		return
	}

	// Print query results
	fmt.Println("Tooling API Query Results:", queryResponse.Records)
}
