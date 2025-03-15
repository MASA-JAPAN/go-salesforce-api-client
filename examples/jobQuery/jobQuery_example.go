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

	// Example: Create a Job Query
	query := "SELECT Id, Name FROM Account"
	jobResponse, err := client.CreateJobQuery(query)
	if err != nil {
		log.Fatalf("Failed to create job query: %v", err)
	}
	fmt.Printf("Job Created: %+v\n", jobResponse)

	// Example: Retrieve Job Query Status
	jobID := jobResponse.ID
	jobDetails, err := client.GetJobQuery(jobID)
	if err != nil {
		log.Fatalf("Failed to get job query details: %v", err)
	}
	fmt.Printf("Job Details: %+v\n", jobDetails)

	// Example: Retrieve Filtered Job Queries
	filteredJobs, err := client.GetFilteredJobQueries(false, "V2Query", "Parallel", "")
	if err != nil {
		log.Fatalf("Failed to get filtered job queries: %v", err)
	}
	fmt.Printf("Filtered Jobs: %+v\n", filteredJobs)

	// Example: Retrieve Job Query Result Pages
	resultPages, err := client.GetJobQueryResultPages(jobID)
	if err != nil {
		log.Fatalf("Failed to get job query result pages: %v", err)
	}
	fmt.Printf("Job Query Result Pages: %+v\n", resultPages)

	// Example: Abort Job Query
	err = client.AbortJobQuery(jobID)
	if err != nil {
		log.Fatalf("Failed to abort job query: %v", err)
	}
	fmt.Println("Job query aborted successfully.")

	// Example: Delete Job Query
	err = client.DeleteJobQuery(jobID)
	if err != nil {
		log.Fatalf("Failed to delete job query: %v", err)
	}
	fmt.Println("Job query deleted successfully.")
}
