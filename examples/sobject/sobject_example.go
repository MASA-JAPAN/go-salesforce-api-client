package main

import (
	"fmt"

	go_salesforce_client "github.com/MASA-JAPAN/go-salesforce-client"
)

func main() {
	auth := go_salesforce_client.Auth{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		Username:     "your_username",
		Password:     "your_password",
		TokenURL:     "https://login.salesforce.com/services/oauth2/token",
	}

	// Authenticate
	client, err := auth.AuthenticatePassword()
	if err != nil {
		fmt.Println("Error authenticating:", err)
		return
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
