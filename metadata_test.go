package go_salesforce_api_client_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func TestDeployMetadata_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify SOAP headers
		if r.Header.Get("Content-Type") != "text/xml; charset=UTF-8" {
			t.Errorf("Expected Content-Type text/xml; charset=UTF-8, got: %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("SOAPAction") != "\"\"" {
			t.Errorf("Expected SOAPAction \"\", got: %s", r.Header.Get("SOAPAction"))
		}

		// Read and verify SOAP request
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)

		if !strings.Contains(bodyStr, "met:deploy") {
			t.Error("Expected deploy operation in SOAP body")
		}

		if !strings.Contains(bodyStr, "met:sessionId") {
			t.Error("Expected session ID in SOAP header")
		}

		// Return mock SOAP response
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <deployResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <id>0Af1X00000XXXXXQAQ</id>
            </result>
        </deployResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.DeployMetadata("base64zipdata", go_salesforce_api_client.MetadataDeployOptions{
		CheckOnly:       true,
		RollbackOnError: true,
		SinglePackage:   true,
		TestLevel:       "NoTestRun",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.ID != "0Af1X00000XXXXXQAQ" {
		t.Errorf("Expected ID 0Af1X00000XXXXXQAQ, got: %s", result.ID)
	}

	if result.Done {
		t.Error("Expected Done to be false for newly initiated deploy")
	}

	if result.State != "Queued" {
		t.Errorf("Expected State Queued, got: %s", result.State)
	}
}

func TestDeployMetadata_MissingAuth(t *testing.T) {
	t.Parallel()

	client := &go_salesforce_api_client.Client{}

	_, err := client.DeployMetadata("base64zip", go_salesforce_api_client.MetadataDeployOptions{})

	if err == nil {
		t.Fatal("Expected error for missing authentication")
	}

	if !strings.Contains(err.Error(), "missing authentication") {
		t.Errorf("Expected missing authentication error, got: %v", err)
	}
}

func TestDeployMetadata_SOAPFault(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapFault := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <soapenv:Fault>
            <faultcode>INVALID_SESSION_ID</faultcode>
            <faultstring>Session expired or invalid</faultstring>
        </soapenv:Fault>
    </soapenv:Body>
</soapenv:Envelope>`

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(soapFault))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "invalid_token",
		InstanceURL: server.URL,
	}

	_, err := client.DeployMetadata("base64zipdata", go_salesforce_api_client.MetadataDeployOptions{})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "INVALID_SESSION_ID") {
		t.Errorf("Expected INVALID_SESSION_ID error, got: %v", err)
	}
}

func TestCheckDeployStatus_InProgress(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <checkDeployStatusResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <done>false</done>
                <id>0Af1X00000XXXXXQAQ</id>
                <status>InProgress</status>
                <numberComponentsDeployed>5</numberComponentsDeployed>
                <numberComponentsTotal>10</numberComponentsTotal>
                <success>false</success>
                <checkOnly>false</checkOnly>
            </result>
        </checkDeployStatusResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.CheckDeployStatus("0Af1X00000XXXXXQAQ")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Done {
		t.Error("Expected Done to be false")
	}

	if result.Status != "InProgress" {
		t.Errorf("Expected status InProgress, got: %s", result.Status)
	}

	if result.NumberComponentsDeployed != 5 {
		t.Errorf("Expected 5 components deployed, got: %d", result.NumberComponentsDeployed)
	}

	if result.NumberComponentsTotal != 10 {
		t.Errorf("Expected 10 total components, got: %d", result.NumberComponentsTotal)
	}
}

func TestCheckDeployStatus_SuccessWithDetails(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <checkDeployStatusResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <done>true</done>
                <id>0Af1X00000XXXXXQAQ</id>
                <status>Succeeded</status>
                <success>true</success>
                <numberComponentsDeployed>10</numberComponentsDeployed>
                <numberComponentsTotal>10</numberComponentsTotal>
                <checkOnly>false</checkOnly>
                <details>
                    <componentSuccesses>
                        <fileName>classes/MyClass.cls</fileName>
                        <fullName>MyClass</fullName>
                        <componentType>ApexClass</componentType>
                        <success>true</success>
                        <created>true</created>
                        <changed>false</changed>
                        <deleted>false</deleted>
                    </componentSuccesses>
                    <componentSuccesses>
                        <fileName>classes/MyClass.cls-meta.xml</fileName>
                        <fullName>MyClass</fullName>
                        <componentType>ApexClass</componentType>
                        <success>true</success>
                        <created>true</created>
                        <changed>false</changed>
                        <deleted>false</deleted>
                    </componentSuccesses>
                </details>
            </result>
        </checkDeployStatusResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.CheckDeployStatus("0Af1X00000XXXXXQAQ")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Done {
		t.Error("Expected Done to be true")
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}

	if result.Status != "Succeeded" {
		t.Errorf("Expected status Succeeded, got: %s", result.Status)
	}

	if result.Details == nil {
		t.Fatal("Expected details to be present")
	}

	if len(result.Details.ComponentSuccesses) != 2 {
		t.Errorf("Expected 2 component successes, got: %d", len(result.Details.ComponentSuccesses))
	}

	if result.Details.ComponentSuccesses[0].FileName != "classes/MyClass.cls" {
		t.Errorf("Expected fileName classes/MyClass.cls, got: %s", result.Details.ComponentSuccesses[0].FileName)
	}

	if result.Details.ComponentSuccesses[0].ComponentType != "ApexClass" {
		t.Errorf("Expected componentType ApexClass, got: %s", result.Details.ComponentSuccesses[0].ComponentType)
	}

	if !result.Details.ComponentSuccesses[0].Created {
		t.Error("Expected component to be created")
	}
}

func TestCheckDeployStatus_FailureWithDetails(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <checkDeployStatusResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <done>true</done>
                <id>0Af1X00000XXXXXQAQ</id>
                <status>Failed</status>
                <success>false</success>
                <numberComponentErrors>1</numberComponentErrors>
                <numberComponentsDeployed>0</numberComponentsDeployed>
                <numberComponentsTotal>1</numberComponentsTotal>
                <checkOnly>false</checkOnly>
                <details>
                    <componentFailures>
                        <fileName>classes/MyClass.cls</fileName>
                        <fullName>MyClass</fullName>
                        <componentType>ApexClass</componentType>
                        <success>false</success>
                        <created>false</created>
                        <changed>false</changed>
                        <deleted>false</deleted>
                        <problem>Invalid syntax at line 5</problem>
                        <problemType>Error</problemType>
                        <lineNumber>5</lineNumber>
                        <columnNumber>12</columnNumber>
                    </componentFailures>
                </details>
            </result>
        </checkDeployStatusResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.CheckDeployStatus("0Af1X00000XXXXXQAQ")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Done {
		t.Error("Expected Done to be true")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}

	if result.Status != "Failed" {
		t.Errorf("Expected status Failed, got: %s", result.Status)
	}

	if result.Details == nil {
		t.Fatal("Expected details to be present")
	}

	if len(result.Details.ComponentFailures) != 1 {
		t.Fatalf("Expected 1 component failure, got: %d", len(result.Details.ComponentFailures))
	}

	failure := result.Details.ComponentFailures[0]
	if failure.FileName != "classes/MyClass.cls" {
		t.Errorf("Expected fileName classes/MyClass.cls, got: %s", failure.FileName)
	}

	if failure.Problem != "Invalid syntax at line 5" {
		t.Errorf("Expected problem 'Invalid syntax at line 5', got: %s", failure.Problem)
	}

	if failure.LineNumber != 5 {
		t.Errorf("Expected lineNumber 5, got: %d", failure.LineNumber)
	}

	if failure.ColumnNumber != 12 {
		t.Errorf("Expected columnNumber 12, got: %d", failure.ColumnNumber)
	}
}

func TestCheckDeployStatus_WithTestResults(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <checkDeployStatusResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <done>true</done>
                <id>0Af1X00000XXXXXQAQ</id>
                <status>Succeeded</status>
                <success>true</success>
                <numberComponentsDeployed>2</numberComponentsDeployed>
                <numberComponentsTotal>2</numberComponentsTotal>
                <numberTestsCompleted>2</numberTestsCompleted>
                <numberTestsTotal>2</numberTestsTotal>
                <numberTestErrors>1</numberTestErrors>
                <checkOnly>false</checkOnly>
                <runTestsEnabled>true</runTestsEnabled>
                <details>
                    <runTestResult>
                        <numFailures>1</numFailures>
                        <numTestsRun>2</numTestsRun>
                        <totalTime>150.5</totalTime>
                        <successes>
                            <id>01p1X000000AAAAQAA</id>
                            <methodName>testMethod1</methodName>
                            <name>MyTestClass</name>
                            <namespace>myns</namespace>
                            <time>50.0</time>
                        </successes>
                        <failures>
                            <id>01p1X000000BBBAQAA</id>
                            <message>Assertion failed</message>
                            <methodName>testMethod2</methodName>
                            <name>MyTestClass</name>
                            <namespace>myns</namespace>
                            <stackTrace>Class.MyTestClass.testMethod2: line 10</stackTrace>
                            <time>100.5</time>
                            <type>AssertionError</type>
                        </failures>
                    </runTestResult>
                </details>
            </result>
        </checkDeployStatusResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.CheckDeployStatus("0Af1X00000XXXXXQAQ")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.RunTestsEnabled {
		t.Error("Expected RunTestsEnabled to be true")
	}

	if result.NumberTestsCompleted != 2 {
		t.Errorf("Expected 2 tests completed, got: %d", result.NumberTestsCompleted)
	}

	if result.Details == nil || result.Details.RunTestResult == nil {
		t.Fatal("Expected run test result to be present")
	}

	testResult := result.Details.RunTestResult
	if testResult.NumTestsRun != 2 {
		t.Errorf("Expected 2 tests run, got: %d", testResult.NumTestsRun)
	}

	if testResult.NumFailures != 1 {
		t.Errorf("Expected 1 test failure, got: %d", testResult.NumFailures)
	}

	if testResult.TotalTime != 150.5 {
		t.Errorf("Expected total time 150.5, got: %f", testResult.TotalTime)
	}

	if len(testResult.Successes) != 1 {
		t.Fatalf("Expected 1 test success, got: %d", len(testResult.Successes))
	}

	if testResult.Successes[0].MethodName != "testMethod1" {
		t.Errorf("Expected method testMethod1, got: %s", testResult.Successes[0].MethodName)
	}

	if len(testResult.Failures) != 1 {
		t.Fatalf("Expected 1 test failure, got: %d", len(testResult.Failures))
	}

	if testResult.Failures[0].MethodName != "testMethod2" {
		t.Errorf("Expected method testMethod2, got: %s", testResult.Failures[0].MethodName)
	}

	if testResult.Failures[0].Message != "Assertion failed" {
		t.Errorf("Expected message 'Assertion failed', got: %s", testResult.Failures[0].Message)
	}
}

func TestCancelDeploy_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <cancelDeployResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <id>0Af1X00000XXXXXQAQ</id>
                <done>true</done>
            </result>
        </cancelDeployResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.CancelDeploy("0Af1X00000XXXXXQAQ")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.ID != "0Af1X00000XXXXXQAQ" {
		t.Errorf("Expected ID 0Af1X00000XXXXXQAQ, got: %s", result.ID)
	}

	if !result.Done {
		t.Error("Expected Done to be true")
	}

	if result.State != "Canceled" {
		t.Errorf("Expected State Canceled, got: %s", result.State)
	}
}

func TestRetrieveMetadata_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify SOAP request
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)

		if !strings.Contains(bodyStr, "met:retrieve") {
			t.Error("Expected retrieve operation in SOAP body")
		}

		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <retrieveResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <id>09S1X00000XXXXXQAQ</id>
            </result>
        </retrieveResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	manifest := `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <types>
        <members>*</members>
        <name>ApexClass</name>
    </types>
    <version>58.0</version>
</Package>`

	result, err := client.RetrieveMetadata(go_salesforce_api_client.MetadataRetrieveOptions{
		ApiVersion:        "58.0",
		SinglePackage:     true,
		UnpackageManifest: manifest,
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.ID != "09S1X00000XXXXXQAQ" {
		t.Errorf("Expected ID 09S1X00000XXXXXQAQ, got: %s", result.ID)
	}

	if result.Done {
		t.Error("Expected Done to be false for newly initiated retrieve")
	}

	if result.State != "Queued" {
		t.Errorf("Expected State Queued, got: %s", result.State)
	}
}

func TestRetrieveMetadata_MissingAuth(t *testing.T) {
	t.Parallel()

	client := &go_salesforce_api_client.Client{}

	_, err := client.RetrieveMetadata(go_salesforce_api_client.MetadataRetrieveOptions{})

	if err == nil {
		t.Fatal("Expected error for missing authentication")
	}

	if !strings.Contains(err.Error(), "missing authentication") {
		t.Errorf("Expected missing authentication error, got: %v", err)
	}
}

func TestCheckRetrieveStatus_InProgress(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <checkRetrieveStatusResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <done>false</done>
                <id>09S1X00000XXXXXQAQ</id>
                <state>InProgress</state>
                <status>InProgress</status>
                <success>false</success>
            </result>
        </checkRetrieveStatusResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.CheckRetrieveStatus("09S1X00000XXXXXQAQ")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Done {
		t.Error("Expected Done to be false")
	}

	if result.State != "InProgress" {
		t.Errorf("Expected state InProgress, got: %s", result.State)
	}
}

func TestCheckRetrieveStatus_SuccessWithZip(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <checkRetrieveStatusResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <done>true</done>
                <id>09S1X00000XXXXXQAQ</id>
                <state>Succeeded</state>
                <status>Succeeded</status>
                <success>true</success>
                <zipFile>UEsDBAoAAAAAAA==</zipFile>
                <fileProperties>
                    <fileName>classes/MyClass.cls</fileName>
                    <fullName>MyClass</fullName>
                    <type>ApexClass</type>
                    <createdById>0051X000000AAAAQAA</createdById>
                    <createdByName>Test User</createdByName>
                    <createdDate>2024-01-01T10:00:00.000Z</createdDate>
                    <id>01p1X000000AAAAQAA</id>
                    <lastModifiedById>0051X000000AAAAQAA</lastModifiedById>
                    <lastModifiedByName>Test User</lastModifiedByName>
                    <lastModifiedDate>2024-01-01T10:00:00.000Z</lastModifiedDate>
                    <manageableState>unmanaged</manageableState>
                </fileProperties>
                <fileProperties>
                    <fileName>classes/MyClass.cls-meta.xml</fileName>
                    <fullName>MyClass</fullName>
                    <type>ApexClass</type>
                    <createdById>0051X000000AAAAQAA</createdById>
                    <createdByName>Test User</createdByName>
                    <createdDate>2024-01-01T10:00:00.000Z</createdDate>
                    <id>01p1X000000BBBAQAA</id>
                    <lastModifiedById>0051X000000AAAAQAA</lastModifiedById>
                    <lastModifiedByName>Test User</lastModifiedByName>
                    <lastModifiedDate>2024-01-01T10:00:00.000Z</lastModifiedDate>
                    <manageableState>unmanaged</manageableState>
                </fileProperties>
            </result>
        </checkRetrieveStatusResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.CheckRetrieveStatus("09S1X00000XXXXXQAQ")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Done {
		t.Error("Expected Done to be true")
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}

	if result.ZipFileBase64 == "" {
		t.Error("Expected ZIP file data to be present")
	}

	if result.ZipFileBase64 != "UEsDBAoAAAAAAA==" {
		t.Errorf("Expected ZIP file UEsDBAoAAAAAAA==, got: %s", result.ZipFileBase64)
	}

	if len(result.FileProperties) != 2 {
		t.Fatalf("Expected 2 file properties, got: %d", len(result.FileProperties))
	}

	if result.FileProperties[0].FileName != "classes/MyClass.cls" {
		t.Errorf("Expected fileName classes/MyClass.cls, got: %s", result.FileProperties[0].FileName)
	}

	if result.FileProperties[0].Type != "ApexClass" {
		t.Errorf("Expected type ApexClass, got: %s", result.FileProperties[0].Type)
	}

	if result.FileProperties[0].CreatedByName != "Test User" {
		t.Errorf("Expected createdByName 'Test User', got: %s", result.FileProperties[0].CreatedByName)
	}
}

func TestCheckRetrieveStatus_WithMessages(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <checkRetrieveStatusResponse xmlns="http://soap.sforce.com/2006/04/metadata">
            <result>
                <done>true</done>
                <id>09S1X00000XXXXXQAQ</id>
                <state>Failed</state>
                <status>Failed</status>
                <success>false</success>
                <errorMessage>Retrieve failed</errorMessage>
                <errorStatusCode>INVALID_TYPE</errorStatusCode>
                <messages>
                    <fileName>classes/InvalidClass.cls</fileName>
                    <problem>Type not found in target org</problem>
                </messages>
            </result>
        </checkRetrieveStatusResponse>
    </soapenv:Body>
</soapenv:Envelope>`

		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(soapResponse))
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	result, err := client.CheckRetrieveStatus("09S1X00000XXXXXQAQ")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Done {
		t.Error("Expected Done to be true")
	}

	if result.Success {
		t.Error("Expected Success to be false")
	}

	if result.ErrorMessage != "Retrieve failed" {
		t.Errorf("Expected error message 'Retrieve failed', got: %s", result.ErrorMessage)
	}

	if result.ErrorStatusCode != "INVALID_TYPE" {
		t.Errorf("Expected error status code INVALID_TYPE, got: %s", result.ErrorStatusCode)
	}

	if len(result.Messages) != 1 {
		t.Fatalf("Expected 1 message, got: %d", len(result.Messages))
	}

	if result.Messages[0].FileName != "classes/InvalidClass.cls" {
		t.Errorf("Expected fileName classes/InvalidClass.cls, got: %s", result.Messages[0].FileName)
	}

	if result.Messages[0].Problem != "Type not found in target org" {
		t.Errorf("Expected problem 'Type not found in target org', got: %s", result.Messages[0].Problem)
	}
}
