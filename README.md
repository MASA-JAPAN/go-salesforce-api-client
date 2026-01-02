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

## 📌 Supported APIs
- **Authentication** (OAuth2)
- **SOQL Queries**
- **CRUD Operations**
- **Tooling API**
- **Bulk Query API**
- **Composite Requests**
- **Limits API** (Monitor API usage)

## 📜 License
MIT License © 2025 MASA-JAPAN

## ⭐ Show Your Support
If you found this useful, please **star this repository** ⭐ and share it with others!

---
🚀 **Go Salesforce API Client** - Making Salesforce Development Easier for Go Developers!

