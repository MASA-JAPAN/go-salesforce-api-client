package main

import (
	"fmt"
	"log"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func main() {
	queryExample()
}

func queryExample() {
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

	// Define the SOQL query
	soql := "SELECT Id, Name FROM Account LIMIT 10"

	// Execute the query
	queryResponse, err := client.Query(soql)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	// Print query results
	fmt.Println("Query Results:")
	for _, record := range queryResponse.Records {
		fmt.Printf("ID: %s, Name: %s\n", record["Id"], record["Name"])
	}
}
