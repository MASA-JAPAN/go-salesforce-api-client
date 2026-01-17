package go_salesforce_api_client_test

import (
	"testing"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
	"github.com/MASA-JAPAN/go-salesforce-emulator/pkg/auth"
	sfemulator "github.com/MASA-JAPAN/go-salesforce-emulator/pkg/emulator"
)

func TestAuthenticatePassword_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New(
		sfemulator.WithCredentials(auth.Credential{
			ClientID:     "test_client_id",
			ClientSecret: "test_client_secret",
			Username:     "test@example.com",
			Password:     "test_password",
		}),
	)
	emu.Start()
	defer emu.Stop()

	auth := go_salesforce_api_client.Auth{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		Username:     "test@example.com",
		Password:     "test_password",
		TokenURL:     emu.URL() + "/services/oauth2/token",
	}

	client, err := auth.AuthenticatePassword()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if client.AccessToken == "" {
		t.Error("Expected access token to be set")
	}

	if client.InstanceURL == "" {
		t.Error("Expected instance URL to be set")
	}
}

func TestAuthenticatePassword_InvalidCredentials(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New(
		sfemulator.WithCredentials(auth.Credential{
			ClientID:     "valid_client_id",
			ClientSecret: "valid_client_secret",
			Username:     "valid@example.com",
			Password:     "valid_password",
		}),
	)
	emu.Start()
	defer emu.Stop()

	auth := go_salesforce_api_client.Auth{
		ClientID:     "wrong_client_id",
		ClientSecret: "wrong_client_secret",
		Username:     "wrong@example.com",
		Password:     "wrong_password",
		TokenURL:     emu.URL() + "/services/oauth2/token",
	}

	_, err := auth.AuthenticatePassword()
	if err == nil {
		t.Error("Expected error for invalid credentials, got nil")
	}
}

func TestAuthenticateClientCredentials_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New(
		sfemulator.WithCredentials(auth.Credential{
			ClientID:     "cc_client_id",
			ClientSecret: "cc_client_secret",
		}),
	)
	emu.Start()
	defer emu.Stop()

	auth := go_salesforce_api_client.Auth{
		ClientID:     "cc_client_id",
		ClientSecret: "cc_client_secret",
		TokenURL:     emu.URL() + "/services/oauth2/token",
	}

	client, err := auth.AuthenticateClientCredentials()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if client.AccessToken == "" {
		t.Error("Expected access token to be set")
	}
}

func TestQuery_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	_, _ = store.CreateRecord("Account", map[string]interface{}{
		"Name":     "Acme Corporation",
		"Industry": "Technology",
	})
	_, _ = store.CreateRecord("Account", map[string]interface{}{
		"Name":     "Global Industries",
		"Industry": "Manufacturing",
	})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	resp, err := client.Query("SELECT Id, Name, Industry FROM Account")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.TotalSize != 2 {
		t.Errorf("Expected TotalSize 2, got: %d", resp.TotalSize)
	}

	if !resp.Done {
		t.Error("Expected Done to be true")
	}

	if len(resp.Records) != 2 {
		t.Errorf("Expected 2 records, got: %d", len(resp.Records))
	}
}

func TestQuery_WithWhereClause(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	_, _ = store.CreateRecord("Account", map[string]interface{}{
		"Name":     "Tech Corp",
		"Industry": "Technology",
	})
	_, _ = store.CreateRecord("Account", map[string]interface{}{
		"Name":     "Manufacturing Inc",
		"Industry": "Manufacturing",
	})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	resp, err := client.Query("SELECT Id, Name FROM Account WHERE Industry = 'Technology'")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.TotalSize != 1 {
		t.Errorf("Expected TotalSize 1, got: %d", resp.TotalSize)
	}

	if len(resp.Records) != 1 {
		t.Errorf("Expected 1 record, got: %d", len(resp.Records))
	}

	if resp.Records[0]["Name"] != "Tech Corp" {
		t.Errorf("Expected Name 'Tech Corp', got: %v", resp.Records[0]["Name"])
	}
}

func TestQuery_EmptyResult(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	resp, err := client.Query("SELECT Id, Name FROM Account")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.TotalSize != 0 {
		t.Errorf("Expected TotalSize 0, got: %d", resp.TotalSize)
	}

	if len(resp.Records) != 0 {
		t.Errorf("Expected 0 records, got: %d", len(resp.Records))
	}
}

func TestCreateRecord_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	record := map[string]interface{}{
		"Name":     "New Account",
		"Industry": "Healthcare",
	}

	resp, err := client.CreateRecord("Account", record)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.ID == "" {
		t.Error("Expected ID to be set")
	}

	if !resp.Success {
		t.Error("Expected Success to be true")
	}
}

func TestGetRecord_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	id, _ := store.CreateRecord("Account", map[string]interface{}{
		"Name":     "Test Account",
		"Industry": "Finance",
	})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	record, err := client.GetRecord("Account", id)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if record["Name"] != "Test Account" {
		t.Errorf("Expected Name 'Test Account', got: %v", record["Name"])
	}

	if record["Industry"] != "Finance" {
		t.Errorf("Expected Industry 'Finance', got: %v", record["Industry"])
	}
}

func TestGetRecord_NotFound(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	_, err := client.GetRecord("Account", "001NONEXISTENT")
	if err == nil {
		t.Error("Expected error for non-existent record, got nil")
	}
}

func TestUpdateRecord_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	id, _ := store.CreateRecord("Account", map[string]interface{}{
		"Name":     "Original Name",
		"Industry": "Tech",
	})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	updates := map[string]interface{}{
		"Name": "Updated Name",
	}

	err := client.UpdateRecord("Account", id, updates)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	record, _ := client.GetRecord("Account", id)
	if record["Name"] != "Updated Name" {
		t.Errorf("Expected Name 'Updated Name', got: %v", record["Name"])
	}

	if record["Industry"] != "Tech" {
		t.Errorf("Expected Industry 'Tech' to remain unchanged, got: %v", record["Industry"])
	}
}

func TestDeleteRecord_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	id, _ := store.CreateRecord("Account", map[string]interface{}{
		"Name": "To Be Deleted",
	})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	err := client.DeleteRecord("Account", id)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	_, err = client.GetRecord("Account", id)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestDescribeSObject_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	describe, err := client.DescribeSObject("Account")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if describe["name"] != "Account" {
		t.Errorf("Expected name 'Account', got: %v", describe["name"])
	}
}

func TestCreateRecords_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	records := []map[string]interface{}{
		{"Name": "Bulk Account 1"},
		{"Name": "Bulk Account 2"},
		{"Name": "Bulk Account 3"},
	}

	resp, err := client.CreateRecords("Account", records)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(resp) != 3 {
		t.Errorf("Expected 3 responses, got: %d", len(resp))
	}

	for i, r := range resp {
		if !r.Success {
			t.Errorf("Expected record %d to be successful", i)
		}
		if r.ID == "" {
			t.Errorf("Expected record %d to have an ID", i)
		}
	}

	queryResp, _ := client.Query("SELECT Id, Name FROM Account")
	if queryResp.TotalSize != 3 {
		t.Errorf("Expected 3 records in database, got: %d", queryResp.TotalSize)
	}
}

func TestUpdateRecords_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	id1, _ := store.CreateRecord("Account", map[string]interface{}{"Name": "Account 1"})
	id2, _ := store.CreateRecord("Account", map[string]interface{}{"Name": "Account 2"})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	updates := []map[string]interface{}{
		{"Id": id1, "Name": "Updated Account 1"},
		{"Id": id2, "Name": "Updated Account 2"},
	}

	err := client.UpdateRecords("Account", updates)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	record1, _ := client.GetRecord("Account", id1)
	if record1["Name"] != "Updated Account 1" {
		t.Errorf("Expected Name 'Updated Account 1', got: %v", record1["Name"])
	}

	record2, _ := client.GetRecord("Account", id2)
	if record2["Name"] != "Updated Account 2" {
		t.Errorf("Expected Name 'Updated Account 2', got: %v", record2["Name"])
	}
}

func TestDeleteRecords_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	id1, _ := store.CreateRecord("Account", map[string]interface{}{"Name": "Delete Me 1"})
	id2, _ := store.CreateRecord("Account", map[string]interface{}{"Name": "Delete Me 2"})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	err := client.DeleteRecords("Account", []string{id1, id2})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	queryResp, _ := client.Query("SELECT Id FROM Account")
	if queryResp.TotalSize != 0 {
		t.Errorf("Expected 0 records after deletion, got: %d", queryResp.TotalSize)
	}
}

func TestCreateJobQuery_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	_, _ = store.CreateRecord("Account", map[string]interface{}{"Name": "Bulk Test Account"})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	resp, err := client.CreateJobQuery("SELECT Id, Name FROM Account")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.ID == "" {
		t.Error("Expected job ID to be set")
	}

	if resp.Object != "Account" {
		t.Errorf("Expected Object 'Account', got: %s", resp.Object)
	}
}

func TestGetJobQuery_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	createResp, _ := client.CreateJobQuery("SELECT Id FROM Account")

	resp, err := client.GetJobQuery(createResp.ID)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.ID != createResp.ID {
		t.Errorf("Expected ID %s, got: %s", createResp.ID, resp.ID)
	}
}

func TestAbortJobQuery_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	createResp, _ := client.CreateJobQuery("SELECT Id FROM Account")

	err := client.AbortJobQuery(createResp.ID)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	jobResp, _ := client.GetJobQuery(createResp.ID)
	if jobResp.State != "Aborted" {
		t.Errorf("Expected State 'Aborted', got: %s", jobResp.State)
	}
}

func TestDeleteJobQuery_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	createResp, _ := client.CreateJobQuery("SELECT Id FROM Account")

	err := client.DeleteJobQuery(createResp.ID)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestQueryToolingAPI_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	resp, err := client.QueryToolingAPI("SELECT Id, Name FROM ApexClass")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Done {
		t.Error("Expected Done to be true")
	}
}

func TestGetLimits_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	limits, err := client.GetLimits()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if limits == nil {
		t.Error("Expected limits to be non-nil")
	}

	dailyRequests, ok := limits["DailyApiRequests"]
	if !ok {
		t.Error("Expected DailyApiRequests to be present")
	}

	dailyMap, ok := dailyRequests.(map[string]interface{})
	if !ok {
		t.Error("Expected DailyApiRequests to be a map")
	}

	if _, ok := dailyMap["Max"]; !ok {
		t.Error("Expected Max field in DailyApiRequests")
	}

	if _, ok := dailyMap["Remaining"]; !ok {
		t.Error("Expected Remaining field in DailyApiRequests")
	}
}

func TestSObjectCRUD_FullWorkflow(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	createResp, err := client.CreateRecord("Account", map[string]interface{}{
		"Name":     "Workflow Test Account",
		"Industry": "Technology",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	accountID := createResp.ID

	record, err := client.GetRecord("Account", accountID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if record["Name"] != "Workflow Test Account" {
		t.Errorf("Expected Name 'Workflow Test Account', got: %v", record["Name"])
	}

	err = client.UpdateRecord("Account", accountID, map[string]interface{}{
		"Industry": "Finance",
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	record, _ = client.GetRecord("Account", accountID)
	if record["Industry"] != "Finance" {
		t.Errorf("Expected Industry 'Finance', got: %v", record["Industry"])
	}

	queryResp, err := client.Query("SELECT Id, Name, Industry FROM Account WHERE Id = '" + accountID + "'")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if queryResp.TotalSize != 1 {
		t.Errorf("Expected 1 record, got: %d", queryResp.TotalSize)
	}

	err = client.DeleteRecord("Account", accountID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = client.GetRecord("Account", accountID)
	if err == nil {
		t.Error("Expected error after deletion")
	}
}

func TestContactRelationship_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	accountResp, err := client.CreateRecord("Account", map[string]interface{}{
		"Name": "Parent Account",
	})
	if err != nil {
		t.Fatalf("Failed to create Account: %v", err)
	}

	contactResp, err := client.CreateRecord("Contact", map[string]interface{}{
		"FirstName": "John",
		"LastName":  "Doe",
		"AccountId": accountResp.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create Contact: %v", err)
	}

	contact, err := client.GetRecord("Contact", contactResp.ID)
	if err != nil {
		t.Fatalf("Failed to get Contact: %v", err)
	}

	if contact["AccountId"] != accountResp.ID {
		t.Errorf("Expected AccountId %s, got: %v", accountResp.ID, contact["AccountId"])
	}
}

func TestQuery_WithOrderBy(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	_, _ = store.CreateRecord("Account", map[string]interface{}{"Name": "Zebra Corp"})
	_, _ = store.CreateRecord("Account", map[string]interface{}{"Name": "Alpha Inc"})
	_, _ = store.CreateRecord("Account", map[string]interface{}{"Name": "Beta LLC"})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	resp, err := client.Query("SELECT Id, Name FROM Account ORDER BY Name ASC")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(resp.Records) != 3 {
		t.Fatalf("Expected 3 records, got: %d", len(resp.Records))
	}

	if resp.Records[0]["Name"] != "Alpha Inc" {
		t.Errorf("Expected first record to be 'Alpha Inc', got: %v", resp.Records[0]["Name"])
	}

	if resp.Records[2]["Name"] != "Zebra Corp" {
		t.Errorf("Expected last record to be 'Zebra Corp', got: %v", resp.Records[2]["Name"])
	}
}

func TestQuery_WithLimit(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	for i := 0; i < 10; i++ {
		_, _ = store.CreateRecord("Account", map[string]interface{}{"Name": "Account"})
	}

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	resp, err := client.Query("SELECT Id, Name FROM Account LIMIT 5")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(resp.Records) != 5 {
		t.Errorf("Expected 5 records, got: %d", len(resp.Records))
	}
}

func TestQuery_WithLikeOperator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	store := emu.Store()
	_, _ = store.CreateRecord("Account", map[string]interface{}{"Name": "Acme Corporation"})
	_, _ = store.CreateRecord("Account", map[string]interface{}{"Name": "Acme Industries"})
	_, _ = store.CreateRecord("Account", map[string]interface{}{"Name": "Global Tech"})

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	resp, err := client.Query("SELECT Id, Name FROM Account WHERE Name LIKE 'Acme%'")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(resp.Records) != 2 {
		t.Errorf("Expected 2 records matching 'Acme%%', got: %d", len(resp.Records))
	}
}

func TestMultipleObjectTypes_WithEmulator(t *testing.T) {
	t.Parallel()

	emu := sfemulator.New()
	emu.Start()
	defer emu.Stop()

	client := &go_salesforce_api_client.Client{
		AccessToken: emu.CreateTestSession(),
		InstanceURL: emu.URL(),
	}

	_, err := client.CreateRecord("Account", map[string]interface{}{"Name": "Test Account"})
	if err != nil {
		t.Fatalf("Failed to create Account: %v", err)
	}

	_, err = client.CreateRecord("Contact", map[string]interface{}{
		"FirstName": "Jane",
		"LastName":  "Smith",
	})
	if err != nil {
		t.Fatalf("Failed to create Contact: %v", err)
	}

	_, err = client.CreateRecord("Lead", map[string]interface{}{
		"FirstName": "Bob",
		"LastName":  "Johnson",
		"Company":   "Test Company",
	})
	if err != nil {
		t.Fatalf("Failed to create Lead: %v", err)
	}

	accountResp, _ := client.Query("SELECT Id FROM Account")
	if accountResp.TotalSize != 1 {
		t.Errorf("Expected 1 Account, got: %d", accountResp.TotalSize)
	}

	contactResp, _ := client.Query("SELECT Id FROM Contact")
	if contactResp.TotalSize != 1 {
		t.Errorf("Expected 1 Contact, got: %d", contactResp.TotalSize)
	}

	leadResp, _ := client.Query("SELECT Id FROM Lead")
	if leadResp.TotalSize != 1 {
		t.Errorf("Expected 1 Lead, got: %d", leadResp.TotalSize)
	}
}

