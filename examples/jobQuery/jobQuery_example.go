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
		TokenURL:     "https://your-domain.my.salesforce.com/services/oauth2/token",
	}

	client, err := auth.AuthenticateClientCredentials()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Create a Job Query
	query := "SELECT Id, Name FROM Account"
	jobResponse, err := client.CreateJobQuery(query)
	if err != nil {
		log.Fatalf("Failed to create job query: %v", err)
	}
	fmt.Printf("Job Created: %+v\n", jobResponse)

	// Wait for Job Completion
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
		time.Sleep(5 * time.Second)
	}

	// Fetch and Print Parsed Results
	var allResults []go_salesforce_api_client.JobQueryResult
	queryLocator := ""
	maxRecords := 10000
	for {
		results, nextLocator, err := client.GetJobQueryResultsParsed(jobID, queryLocator, maxRecords)
		if err != nil {
			log.Fatalf("Failed to get job query results: %v", err)
		}

		allResults = append(allResults, results...)
		fmt.Printf("Fetched %d records\n", len(results))

		if nextLocator == "" || nextLocator == "null" {
			break
		}

		queryLocator = nextLocator
	}

	// Display Data
	fmt.Println("Total Job Results:")
	for _, row := range allResults {
		fmt.Println(row)
	}

	// Delete the Job Query
	err = client.DeleteJobQuery(jobID)
	if err != nil {
		log.Fatalf("Failed to delete job query: %v", err)
	}
	fmt.Println("Job query deleted successfully.")
}
