package main

import (
	"fmt"
	"log"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func main() {
	CreateRecordsExample()
	UpdateRecordsExample()
	DeleteRecordsExample()
}

func CreateRecordsExample() {
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

	// Bulk Create Records
	records := []map[string]interface{}{
		{"Name": "Sample Corp A"},
		{"Name": "Sample Corp B"},
		{"Name": "Sample Corp C"},
	}

	response, err := client.CreateRecords("Account", records)
	if err != nil {
		fmt.Println("Error creating records:", err)
		return
	}

	fmt.Println("Record Creation Response:", response)
}

func UpdateRecordsExample() {
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

	records := []map[string]interface{}{
		{"Id": "001IR00001ulZ5TYAU", "Name": "Updated Sample Corp A"},
		{"Id": "001IR00001ulZ5UYAU", "Name": "Updated Sample Corp B"},
		{"Id": "001IR00001ulZ5VYAU", "Name": "Updated Sample Corp C"},
	}

	err = client.UpdateRecords("Account", records)
	if err != nil {
		fmt.Println("Error updating records:", err)
		return
	}
}

func DeleteRecordsExample() {
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

	ids := []string{"001IR00001ulZ5YYAU", "001IR00001ulZ5ZYAU", "001IR00001ulZ5aYAE"}

	err = client.DeleteRecords("Account", ids)
	if err != nil {
		fmt.Println("Error updating records:", err)
		return
	}
}
