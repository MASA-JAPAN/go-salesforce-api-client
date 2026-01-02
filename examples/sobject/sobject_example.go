package main

import (
	"fmt"
	"log"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func main() {
	CrudExample()
	DescribeExample()
}

func DescribeExample() {
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

	// Describe an SObject
	sObjectType := "Account"
	describe, err := client.DescribeSObject(sObjectType)
	if err != nil {
		fmt.Println("Error retrieving SObject description:", err)
		return
	}

	// Print describe details
	fmt.Println("SObject Description:", describe)
}

func CrudExample() {
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

	// Create a record
	record := map[string]interface{}{
		"Name": "Example Account",
	}
	createdRecord, err := client.CreateRecord("Account", record)
	if err != nil {
		fmt.Println("Error creating record:", err)
		return
	}
	fmt.Println("Record Created! ID:", createdRecord.ID)

	// Retrieve the record
	retrievedRecord, err := client.GetRecord("Account", createdRecord.ID)
	if err != nil {
		fmt.Println("Error retrieving record:", err)
		return
	}
	fmt.Println("Retrieved Record:", retrievedRecord)

	// Update the record
	updates := map[string]interface{}{
		"Name": "Updated Account Name",
	}
	err = client.UpdateRecord("Account", createdRecord.ID, updates)
	if err != nil {
		fmt.Println("Error updating record:", err)
		return
	}
	fmt.Println("Record Updated Successfully!")

	// Delete the record
	err = client.DeleteRecord("Account", createdRecord.ID)
	if err != nil {
		fmt.Println("Error deleting record:", err)
		return
	}
	fmt.Println("Record Deleted Successfully!")
}
