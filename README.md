# 🚀 Go Salesforce API Client

A **lightweight, fast, and developer-friendly** Go client for interacting with **Salesforce APIs**. This library provides easy access to **CRUD operations, SOQL queries, Tooling API, and authentication**. 

## 🎯 Features
✅ **Easy Authentication**: Supports OAuth2 Password Flow and Client Credentials Flow.
✅ **SOQL Query Support**: Execute complex SOQL queries with ease.
✅ **CRUD Operations**: Perform create, read, update, delete on any Salesforce object.
✅ **Tooling API Access**: Interact with metadata and developer tooling API.
✅ **Bulk Query API Support**: Efficiently fetch large datasets using Salesforce Bulk Query Jobs with automatic pagination.
✅ **Well-Tested**: 90%+ test coverage with `httptest`-based mocks.
✅ **Lightweight & Fast**: Minimal dependencies for blazing-fast performance.

## 📦 Installation
```sh
go install github.com/MASA-JAPAN/go-salesforce-api-client
```

## 🚀 Quick Start
### 1️⃣ Authenticate with Salesforce
```go
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
```

### 2️⃣ Query Salesforce Data
```go
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
```

### 3️⃣ Create a New Records
```go
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
```

### 4️⃣ Update a Records
```go
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
```

### 5️⃣ Delete a Record
```go
ids := []string{"001IR00001ulZ5YYAU", "001IR00001ulZ5ZYAU", "001IR00001ulZ5aYAE"}

err = client.DeleteRecords("Account", ids)
if err != nil {
    fmt.Println("Error updating records:", err)
    return
}
```

### 6️⃣ Deploy Metadata
```go
// Create deployment ZIP (package.xml + metadata files)
zipBase64 := createDeploymentZip() // Your function to create base64-encoded ZIP

// Configure deploy options
options := go_salesforce_api_client.MetadataDeployOptions{
    CheckOnly:       false, // Set true for validation-only
    RollbackOnError: true,
    SinglePackage:   true,
    TestLevel:       "NoTestRun",
}

// Initiate deployment
asyncResult, err := client.DeployMetadata(zipBase64, options)
if err != nil {
    fmt.Println("Deploy failed:", err)
    return
}

// Poll for completion
for {
    status, err := client.CheckDeployStatus(asyncResult.ID)
    if err != nil {
        fmt.Println("Error checking status:", err)
        return
    }

    fmt.Printf("Status: %s (%d/%d components)\n",
        status.Status,
        status.NumberComponentsDeployed,
        status.NumberComponentsTotal)

    if status.Done {
        if status.Success {
            fmt.Println("Deploy succeeded!")
        } else {
            // Print failures
            for _, failure := range status.Details.ComponentFailures {
                fmt.Printf("%s: %s [Line %d]\n",
                    failure.FileName,
                    failure.Problem,
                    failure.LineNumber)
            }
        }
        break
    }

    time.Sleep(5 * time.Second)
}
```

### 7️⃣ Retrieve Metadata
```go
// Define package manifest
manifest := `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <types>
        <members>*</members>
        <name>ApexClass</name>
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
    fmt.Println("Retrieve failed:", err)
    return
}

// Poll for completion
for {
    status, err := client.CheckRetrieveStatus(asyncResult.ID)
    if err != nil {
        fmt.Println("Error checking status:", err)
        return
    }

    if status.Done {
        if status.Success {
            // Save ZIP file
            zipData, _ := base64.StdEncoding.DecodeString(status.ZipFileBase64)
            os.WriteFile("metadata.zip", zipData, 0644)
            fmt.Println("Metadata retrieved successfully!")
        }
        break
    }

    time.Sleep(5 * time.Second)
}
```

## 📌 Supported APIs
- **Authentication** (OAuth2)
- **SOQL Queries**
- **CRUD Operations**
- **Tooling API**
- **Bulk Query API**
- **Composite Requests**
- **Limits API** (Monitor API usage)
- **Metadata API** (Deploy & Retrieve metadata packages)

## 📜 License
MIT License © 2025 MASA-JAPAN

## ⭐ Show Your Support
If you found this useful, please **star this repository** ⭐ and share it with others!

---
🚀 **Go Salesforce API Client** - Making Salesforce Development Easier for Go Developers!

