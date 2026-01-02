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
		TokenURL:     "https://your-domain.my.salesforce.com/services/oauth2/token",
	}
	// Authenticate and retrieve an access token
	client, err := auth.AuthenticateClientCredentials()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Get record counts for specified objects
	objects := []string{"Account", "Contact", "Lead"}
	recordCounts, err := client.GetRecordCounts(objects)
	if err != nil {
		fmt.Println("Error retrieving record counts:", err)
		return
	}

	// Print record counts
	fmt.Println("Salesforce Record Counts:", recordCounts)
}
