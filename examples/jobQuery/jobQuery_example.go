package main

import (
	"fmt"
	"log"
	"time"

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

	// Example: Create a Job Query
	query := "SELECT Id, Name FROM Account"
	jobResponse, err := client.CreateJobQuery(query)
	if err != nil {
		log.Fatalf("Failed to create job query: %v", err)
	}
	fmt.Printf("Job Created: %+v\n", jobResponse)

	// Example: Wait for Job Query Completion
	jobID := jobResponse.ID
	for {
		jobDetails, err := client.GetJobQuery(jobID)
		if err != nil {
			log.Fatalf("Failed to get job query details: %v", err)
		}
		fmt.Printf("Job Details: %+v\n", jobDetails)

		if jobDetails.State == "JobComplete" {
			break
		}
		fmt.Println("Waiting for job to complete...")
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	// Example: Get All Job Query Results using Locator and maxRecords
	allResults := ""
	queryLocator := ""
	maxRecords := 10000
	for {
		partialResults, nextLocator, err := client.GetJobQueryResults(jobID, queryLocator, maxRecords)
		if err != nil {
			log.Fatalf("Failed to get job query results: %v", err)
		}
		allResults += partialResults
		fmt.Printf("Fetched %d records\n", len(partialResults))

		fmt.Printf("Next Locator: %s\n", nextLocator)

		if nextLocator == "null" {
			break
		}

		queryLocator = nextLocator
	}

	fmt.Printf("Total Job Results:\n%s\n", allResults)

	err = client.DeleteJobQuery(jobID)
	if err != nil {
		log.Fatalf("Failed to delete job query: %v", err)
	}
	fmt.Println("Job query deleted successfully.")
}
