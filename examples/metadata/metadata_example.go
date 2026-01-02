package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func main() {
	// Initialize authentication
	auth := go_salesforce_api_client.Auth{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		TokenURL:     "https://your-domain.my.salesforce.com/services/oauth2/token",
	}

	client, err := auth.AuthenticateClientCredentials()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Example 1: Deploy Metadata
	fmt.Println("=== Deploy Metadata Example ===")
	deployExample(client)

	// Example 2: Retrieve Metadata
	fmt.Println("\n=== Retrieve Metadata Example ===")
	retrieveExample(client)
}

func deployExample(client *go_salesforce_api_client.Client) {
	// Create a deployment ZIP
	zipBase64, err := createDeploymentPackage()
	if err != nil {
		log.Fatalf("Failed to create deployment package: %v", err)
	}

	// Configure deploy options
	options := go_salesforce_api_client.MetadataDeployOptions{
		CheckOnly:       false, // Set to true for validation-only deploy
		RollbackOnError: true,
		SinglePackage:   true,
		TestLevel:       "NoTestRun", // or "RunSpecifiedTests", "RunLocalTests", "RunAllTestsInOrg"
		IgnoreWarnings:  false,
	}

	// Initiate deployment
	asyncResult, err := client.DeployMetadata(zipBase64, options)
	if err != nil {
		log.Fatalf("Failed to deploy metadata: %v", err)
	}

	fmt.Printf("Deploy initiated. Async Process ID: %s\n", asyncResult.ID)

	// Poll for completion
	for {
		status, err := client.CheckDeployStatus(asyncResult.ID)
		if err != nil {
			log.Fatalf("Failed to check deploy status: %v", err)
		}

		fmt.Printf("Deploy Status: %s (%d/%d components)\n",
			status.Status,
			status.NumberComponentsDeployed,
			status.NumberComponentsTotal)

		if status.Done {
			if status.Success {
				fmt.Println("Deployment completed successfully!")

				// Print component successes
				if status.Details != nil && len(status.Details.ComponentSuccesses) > 0 {
					fmt.Println("\nDeployed Components:")
					for _, success := range status.Details.ComponentSuccesses {
						action := "Modified"
						if success.Created {
							action = "Created"
						} else if success.Deleted {
							action = "Deleted"
						}
						fmt.Printf("  - %s: %s (%s)\n", action, success.FullName, success.ComponentType)
					}
				}
			} else {
				fmt.Println("Deployment failed!")

				// Print component failures
				if status.Details != nil && len(status.Details.ComponentFailures) > 0 {
					fmt.Println("\nFailures:")
					for _, failure := range status.Details.ComponentFailures {
						fmt.Printf("  - %s (%s): %s [Line %d, Col %d]\n",
							failure.FileName,
							failure.ComponentType,
							failure.Problem,
							failure.LineNumber,
							failure.ColumnNumber)
					}
				}

				// Print test failures if any
				if status.Details != nil && status.Details.RunTestResult != nil && len(status.Details.RunTestResult.Failures) > 0 {
					fmt.Println("\nTest Failures:")
					for _, testFailure := range status.Details.RunTestResult.Failures {
						fmt.Printf("  - %s.%s: %s\n",
							testFailure.Name,
							testFailure.MethodName,
							testFailure.Message)
						if testFailure.StackTrace != "" {
							fmt.Printf("    Stack Trace: %s\n", testFailure.StackTrace)
						}
					}
				}
			}
			break
		}

		// Wait before next poll
		time.Sleep(5 * time.Second)
	}
}

func retrieveExample(client *go_salesforce_api_client.Client) {
	// Define package manifest
	manifest := `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <types>
        <members>*</members>
        <name>ApexClass</name>
    </types>
    <types>
        <members>*</members>
        <name>CustomObject</name>
    </types>
    <version>58.0</version>
</Package>`

	// Configure retrieve options
	options := go_salesforce_api_client.MetadataRetrieveOptions{
		ApiVersion:        "58.0",
		SinglePackage:     true,
		UnpackageManifest: manifest,
	}

	// Initiate retrieve
	asyncResult, err := client.RetrieveMetadata(options)
	if err != nil {
		log.Fatalf("Failed to retrieve metadata: %v", err)
	}

	fmt.Printf("Retrieve initiated. Async Process ID: %s\n", asyncResult.ID)

	// Poll for completion
	for {
		status, err := client.CheckRetrieveStatus(asyncResult.ID)
		if err != nil {
			log.Fatalf("Failed to check retrieve status: %v", err)
		}

		fmt.Printf("Retrieve Status: %s\n", status.State)

		if status.Done {
			if status.Success {
				fmt.Println("Retrieve completed successfully!")

				// Print retrieved files
				fmt.Printf("\nRetrieved %d files:\n", len(status.FileProperties))
				for _, file := range status.FileProperties {
					fmt.Printf("  - %s (%s)\n", file.FileName, file.Type)
				}

				// Save ZIP file
				if err := saveZipFile(status.ZipFileBase64, "retrieved_metadata.zip"); err != nil {
					log.Printf("Failed to save ZIP file: %v", err)
				} else {
					fmt.Println("\nMetadata saved to: retrieved_metadata.zip")
				}
			} else {
				fmt.Println("Retrieve failed!")
				if status.ErrorMessage != "" {
					fmt.Printf("Error: %s\n", status.ErrorMessage)
				}
				for _, msg := range status.Messages {
					fmt.Printf("  - %s: %s\n", msg.FileName, msg.Problem)
				}
			}
			break
		}

		// Wait before next poll
		time.Sleep(5 * time.Second)
	}
}

// createDeploymentPackage creates a sample deployment ZIP package
func createDeploymentPackage() (string, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Create package.xml
	packageXML := `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <types>
        <members>MyApexClass</members>
        <name>ApexClass</name>
    </types>
    <version>58.0</version>
</Package>`

	packageFile, err := zipWriter.Create("package.xml")
	if err != nil {
		return "", err
	}
	if _, err := packageFile.Write([]byte(packageXML)); err != nil {
		return "", err
	}

	// Create Apex class
	apexClass := `public class MyApexClass {
    public static String getMessage() {
        return 'Hello from MyApexClass!';
    }

    public static Integer add(Integer a, Integer b) {
        return a + b;
    }
}`

	apexFile, err := zipWriter.Create("classes/MyApexClass.cls")
	if err != nil {
		return "", err
	}
	if _, err := apexFile.Write([]byte(apexClass)); err != nil {
		return "", err
	}

	// Create Apex class metadata
	apexMeta := `<?xml version="1.0" encoding="UTF-8"?>
<ApexClass xmlns="http://soap.sforce.com/2006/04/metadata">
    <apiVersion>58.0</apiVersion>
    <status>Active</status>
</ApexClass>`

	metaFile, err := zipWriter.Create("classes/MyApexClass.cls-meta.xml")
	if err != nil {
		return "", err
	}
	if _, err := metaFile.Write([]byte(apexMeta)); err != nil {
		return "", err
	}

	if err := zipWriter.Close(); err != nil {
		return "", err
	}

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}

// saveZipFile saves base64-encoded ZIP data to a file
func saveZipFile(base64Data, filename string) error {
	// Decode base64
	zipData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, zipData, 0644)
}

// extractZipFile extracts a ZIP file to a directory.
// This is an example utility function to help users understand how to extract
// the retrieved metadata ZIP files to a local directory.
//
//nolint:unused // Example utility function for user reference
func extractZipFile(base64Data, outputDir string) error {
	// Decode base64
	zipData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	// Read ZIP
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}

	// Extract files
	for _, file := range zipReader.File {
		// Create output path
		outputPath := outputDir + "/" + file.Name

		if file.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return err
			}
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(outputPath[:len(outputPath)-len(file.Name)], 0755); err != nil {
			return err
		}

		// Extract file
		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		fileData, err := os.Create(outputPath)
		if err != nil {
			fileReader.Close()
			return err
		}

		if _, err := fileData.ReadFrom(fileReader); err != nil {
			fileData.Close()
			fileReader.Close()
			return err
		}

		fileData.Close()
		fileReader.Close()
	}

	return nil
}
