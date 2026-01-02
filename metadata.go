package go_salesforce_api_client

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MetadataDeployOptions configures deploy behavior
type MetadataDeployOptions struct {
	AllowMissingFiles bool
	AutoUpdatePackage bool
	CheckOnly         bool // Validation deploy only
	IgnoreWarnings    bool
	PerformRetrieve   bool
	PurgeOnDelete     bool
	RollbackOnError   bool
	RunTests          []string // Test classes to run
	SinglePackage     bool
	TestLevel         string // NoTestRun, RunSpecifiedTests, RunLocalTests, RunAllTestsInOrg
}

// MetadataRetrieveOptions configures retrieve behavior
type MetadataRetrieveOptions struct {
	ApiVersion        string
	PackageNames      []string
	SinglePackage     bool
	SpecificFiles     []string
	UnpackageManifest string // XML manifest content
}

// MetadataAsyncResult represents async operation status
type MetadataAsyncResult struct {
	ID              string
	Done            bool
	State           string
	StateDetail     string
	StatusCode      string
	ErrorMessage    string
	ErrorStatusCode string
}

// MetadataDeployResult represents deploy operation result
type MetadataDeployResult struct {
	CheckOnly                bool
	CompletedDate            string
	CreatedDate              string
	Details                  *DeployDetails
	Done                     bool
	ErrorMessage             string
	ID                       string
	IgnoreWarnings           bool
	NumberComponentErrors    int
	NumberComponentsDeployed int
	NumberComponentsTotal    int
	NumberTestErrors         int
	NumberTestsCompleted     int
	NumberTestsTotal         int
	RollbackOnError          bool
	RunTestsEnabled          bool
	StartDate                string
	Status                   string
	Success                  bool
}

// DeployDetails contains detailed deploy results
type DeployDetails struct {
	ComponentSuccesses []ComponentSuccess
	ComponentFailures  []ComponentFailure
	RunTestResult      *RunTestResult
}

// ComponentSuccess represents successful component
type ComponentSuccess struct {
	Changed       bool
	Created       bool
	Deleted       bool
	FileName      string
	FullName      string
	ComponentType string
	Success       bool
}

// ComponentFailure represents failed component
type ComponentFailure struct {
	Changed       bool
	Created       bool
	Deleted       bool
	FileName      string
	FullName      string
	ComponentType string
	Problem       string
	ProblemType   string
	LineNumber    int
	ColumnNumber  int
	Success       bool
}

// RunTestResult contains test execution results
type RunTestResult struct {
	NumFailures  int
	NumTestsRun  int
	TotalTime    float64
	Successes    []TestSuccess
	Failures     []TestFailure
	CodeCoverage []CodeCoverageResult
}

// TestSuccess represents passed test
type TestSuccess struct {
	ID         string
	MethodName string
	Name       string
	Namespace  string
	Time       float64
}

// TestFailure represents failed test
type TestFailure struct {
	ID         string
	Message    string
	MethodName string
	Name       string
	Namespace  string
	StackTrace string
	Time       float64
	Type       string
}

// CodeCoverageResult represents code coverage for a class
type CodeCoverageResult struct {
	ID                     string
	LocationsNotCovered    []CodeLocation
	Name                   string
	Namespace              string
	NumLocations           int
	NumLocationsNotCovered int
	Type                   string
}

// CodeLocation represents a code location
type CodeLocation struct {
	Column        int
	Line          int
	NumExecutions int
	Time          float64
}

// MetadataRetrieveResult represents retrieve operation result
type MetadataRetrieveResult struct {
	Done            bool
	ErrorMessage    string
	ErrorStatusCode string
	FileProperties  []FileProperty
	ID              string
	Messages        []RetrieveMessage
	State           string
	Status          string
	Success         bool
	ZipFileBase64   string
}

// FileProperty represents metadata about a retrieved file
type FileProperty struct {
	CreatedByID        string
	CreatedByName      string
	CreatedDate        string
	FileName           string
	FullName           string
	ID                 string
	LastModifiedByID   string
	LastModifiedByName string
	LastModifiedDate   string
	ManageableState    string
	NamespacePrefix    string
	Type               string
}

// RetrieveMessage represents a retrieve status message
type RetrieveMessage struct {
	FileName string
	Problem  string
}

// soapEnvelope represents SOAP request envelope
type soapEnvelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	SOAPENV string   `xml:"xmlns:soapenv,attr"`
	MET     string   `xml:"xmlns:met,attr"`
	Header  soapHeader
	Body    soapBody
}

// soapHeader contains SOAP header with session
type soapHeader struct {
	XMLName   xml.Name `xml:"soapenv:Header"`
	SessionID string   `xml:"met:SessionHeader>met:sessionId"`
}

// soapBody contains SOAP operation
type soapBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Content string   `xml:",innerxml"`
}

// soapFault represents SOAP error response
type soapFault struct {
	XMLName xml.Name      `xml:"Envelope"`
	Body    soapFaultBody `xml:"Body"`
}

// soapFaultBody contains fault details
type soapFaultBody struct {
	Fault fault `xml:"Fault"`
}

// fault represents SOAP fault
type fault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
}

// Error variables
var (
	ErrSOAPFault       = errors.New("SOAP fault occurred")
	ErrInvalidSession  = errors.New("invalid session ID")
	ErrDeployFailed    = errors.New("deployment failed")
	ErrRetrieveFailed  = errors.New("retrieve failed")
	ErrInvalidZip      = errors.New("invalid ZIP file")
	ErrInvalidManifest = errors.New("invalid package manifest")
)

// getMetadataAPIVersion returns the Metadata API version
func (c *Client) getMetadataAPIVersion() string {
	return "58.0"
}

// buildSOAPEnvelope constructs a SOAP envelope for Metadata API requests
func (c *Client) buildSOAPEnvelope(bodyContent string) *soapEnvelope {
	return &soapEnvelope{
		SOAPENV: "http://schemas.xmlsoap.org/soap/envelope/",
		MET:     "http://soap.sforce.com/2006/04/metadata",
		Header: soapHeader{
			SessionID: c.AccessToken,
		},
		Body: soapBody{
			Content: bodyContent,
		},
	}
}

// sendSOAPRequest sends a SOAP request and returns the response body
func (c *Client) sendSOAPRequest(endpoint string, envelope *soapEnvelope) ([]byte, error) {
	// Marshal envelope to XML
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, err
	}

	// Add XML header
	xmlRequest := []byte(xml.Header + string(xmlData))

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(xmlRequest))
	if err != nil {
		return nil, err
	}

	// Set SOAP headers
	req.Header.Set("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Set("SOAPAction", "\"\"")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for SOAP faults
	if resp.StatusCode >= 400 || bytes.Contains(body, []byte("<faultcode>")) {
		return nil, c.parseSOAPFault(body)
	}

	return body, nil
}

// parseSOAPFault parses SOAP fault responses
func (c *Client) parseSOAPFault(responseBody []byte) error {
	var soapFault soapFault
	if err := xml.Unmarshal(responseBody, &soapFault); err != nil {
		return fmt.Errorf("%w: failed to parse SOAP fault: %w", ErrSOAPFault, err)
	}

	faultString := soapFault.Body.Fault.FaultString

	// Check for specific error types
	if strings.Contains(faultString, "INVALID_SESSION_ID") {
		return fmt.Errorf("%w: %s", ErrInvalidSession, faultString)
	}

	return fmt.Errorf("%w: [%s] %s",
		ErrSOAPFault,
		soapFault.Body.Fault.FaultCode,
		faultString)
}

// DeployMetadata initiates an asynchronous metadata deployment
func (c *Client) DeployMetadata(zipFileBase64 string, options MetadataDeployOptions) (*MetadataAsyncResult, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	// Build SOAP endpoint
	apiVersion := c.getMetadataAPIVersion()
	endpoint := fmt.Sprintf("%s/services/Soap/m/%s", c.InstanceURL, apiVersion)

	// Build deploy request body XML
	bodyContent := fmt.Sprintf(`<met:deploy>
      <met:ZipFile>%s</met:ZipFile>
      <met:DeployOptions>
        <met:allowMissingFiles>%t</met:allowMissingFiles>
        <met:autoUpdatePackage>%t</met:autoUpdatePackage>
        <met:checkOnly>%t</met:checkOnly>
        <met:ignoreWarnings>%t</met:ignoreWarnings>
        <met:performRetrieve>%t</met:performRetrieve>
        <met:purgeOnDelete>%t</met:purgeOnDelete>
        <met:rollbackOnError>%t</met:rollbackOnError>
        <met:singlePackage>%t</met:singlePackage>
        <met:testLevel>%s</met:testLevel>
      </met:DeployOptions>
    </met:deploy>`,
		zipFileBase64,
		options.AllowMissingFiles,
		options.AutoUpdatePackage,
		options.CheckOnly,
		options.IgnoreWarnings,
		options.PerformRetrieve,
		options.PurgeOnDelete,
		options.RollbackOnError,
		options.SinglePackage,
		options.TestLevel,
	)

	// Build SOAP envelope
	envelope := c.buildSOAPEnvelope(bodyContent)

	// Send SOAP request
	responseBody, err := c.sendSOAPRequest(endpoint, envelope)
	if err != nil {
		return nil, err
	}

	// Parse response to extract async process ID
	var result struct {
		ID string `xml:"Body>deployResponse>result>id"`
	}

	if err := xml.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse deploy response: %w", err)
	}

	return &MetadataAsyncResult{
		ID:    result.ID,
		Done:  false,
		State: "Queued",
	}, nil
}

// Helper types for CheckDeployStatus XML parsing
type deployStatusResponse struct {
	Result struct {
		CheckOnly                bool             `xml:"checkOnly"`
		CompletedDate            string           `xml:"completedDate"`
		CreatedDate              string           `xml:"createdDate"`
		Done                     bool             `xml:"done"`
		ErrorMessage             string           `xml:"errorMessage"`
		ID                       string           `xml:"id"`
		IgnoreWarnings           bool             `xml:"ignoreWarnings"`
		NumberComponentErrors    int              `xml:"numberComponentErrors"`
		NumberComponentsDeployed int              `xml:"numberComponentsDeployed"`
		NumberComponentsTotal    int              `xml:"numberComponentsTotal"`
		NumberTestErrors         int              `xml:"numberTestErrors"`
		NumberTestsCompleted     int              `xml:"numberTestsCompleted"`
		NumberTestsTotal         int              `xml:"numberTestsTotal"`
		RollbackOnError          bool             `xml:"rollbackOnError"`
		RunTestsEnabled          bool             `xml:"runTestsEnabled"`
		StartDate                string           `xml:"startDate"`
		Status                   string           `xml:"status"`
		Success                  bool             `xml:"success"`
		Details                  deployDetailsXML `xml:"details"`
	} `xml:"Body>checkDeployStatusResponse>result"`
}

type deployDetailsXML struct {
	ComponentSuccesses []componentSuccessXML `xml:"componentSuccesses"`
	ComponentFailures  []componentFailureXML `xml:"componentFailures"`
	RunTestResult      runTestResultXML      `xml:"runTestResult"`
}

type componentSuccessXML struct {
	Changed       bool   `xml:"changed"`
	Created       bool   `xml:"created"`
	Deleted       bool   `xml:"deleted"`
	FileName      string `xml:"fileName"`
	FullName      string `xml:"fullName"`
	ComponentType string `xml:"componentType"`
	Success       bool   `xml:"success"`
}

type componentFailureXML struct {
	Changed       bool   `xml:"changed"`
	Created       bool   `xml:"created"`
	Deleted       bool   `xml:"deleted"`
	FileName      string `xml:"fileName"`
	FullName      string `xml:"fullName"`
	ComponentType string `xml:"componentType"`
	Problem       string `xml:"problem"`
	ProblemType   string `xml:"problemType"`
	LineNumber    int    `xml:"lineNumber"`
	ColumnNumber  int    `xml:"columnNumber"`
	Success       bool   `xml:"success"`
}

type runTestResultXML struct {
	NumFailures  int               `xml:"numFailures"`
	NumTestsRun  int               `xml:"numTestsRun"`
	TotalTime    float64           `xml:"totalTime"`
	Successes    []testSuccessXML  `xml:"successes"`
	Failures     []testFailureXML  `xml:"failures"`
	CodeCoverage []codeCoverageXML `xml:"codeCoverage"`
}

type testSuccessXML struct {
	ID         string  `xml:"id"`
	MethodName string  `xml:"methodName"`
	Name       string  `xml:"name"`
	Namespace  string  `xml:"namespace"`
	Time       float64 `xml:"time"`
}

type testFailureXML struct {
	ID         string  `xml:"id"`
	Message    string  `xml:"message"`
	MethodName string  `xml:"methodName"`
	Name       string  `xml:"name"`
	Namespace  string  `xml:"namespace"`
	StackTrace string  `xml:"stackTrace"`
	Time       float64 `xml:"time"`
	Type       string  `xml:"type"`
}

type codeCoverageXML struct {
	ID                     string            `xml:"id"`
	Name                   string            `xml:"name"`
	Namespace              string            `xml:"namespace"`
	NumLocations           int               `xml:"numLocations"`
	NumLocationsNotCovered int               `xml:"numLocationsNotCovered"`
	Type                   string            `xml:"type"`
	LocationsNotCovered    []codeLocationXML `xml:"locationsNotCovered"`
}

type codeLocationXML struct {
	Column        int     `xml:"column"`
	Line          int     `xml:"line"`
	NumExecutions int     `xml:"numExecutions"`
	Time          float64 `xml:"time"`
}

func convertDeployDetails(xmlDetails deployDetailsXML) *DeployDetails {
	details := &DeployDetails{}

	// Convert component successes
	//nolint:staticcheck // Cannot use type conversion - XML struct has different tags
	for _, s := range xmlDetails.ComponentSuccesses {
		details.ComponentSuccesses = append(details.ComponentSuccesses, ComponentSuccess{
			Changed:       s.Changed,
			Created:       s.Created,
			Deleted:       s.Deleted,
			FileName:      s.FileName,
			FullName:      s.FullName,
			ComponentType: s.ComponentType,
			Success:       s.Success,
		})
	}

	// Convert component failures
	//nolint:staticcheck // Cannot use type conversion - XML struct has different tags
	for _, f := range xmlDetails.ComponentFailures {
		details.ComponentFailures = append(details.ComponentFailures, ComponentFailure{
			Changed:       f.Changed,
			Created:       f.Created,
			Deleted:       f.Deleted,
			FileName:      f.FileName,
			FullName:      f.FullName,
			ComponentType: f.ComponentType,
			Problem:       f.Problem,
			ProblemType:   f.ProblemType,
			LineNumber:    f.LineNumber,
			ColumnNumber:  f.ColumnNumber,
			Success:       f.Success,
		})
	}

	// Convert test results if present
	if xmlDetails.RunTestResult.NumTestsRun > 0 {
		details.RunTestResult = convertTestResult(xmlDetails.RunTestResult)
	}

	return details
}

func convertTestResult(xmlTest runTestResultXML) *RunTestResult {
	testResult := &RunTestResult{
		NumFailures: xmlTest.NumFailures,
		NumTestsRun: xmlTest.NumTestsRun,
		TotalTime:   xmlTest.TotalTime,
	}

	// Convert test successes
	//nolint:staticcheck // Cannot use type conversion - XML struct has different tags
	for _, ts := range xmlTest.Successes {
		testResult.Successes = append(testResult.Successes, TestSuccess{
			ID:         ts.ID,
			MethodName: ts.MethodName,
			Name:       ts.Name,
			Namespace:  ts.Namespace,
			Time:       ts.Time,
		})
	}

	// Convert test failures
	//nolint:staticcheck // Cannot use type conversion - XML struct has different tags
	for _, tf := range xmlTest.Failures {
		testResult.Failures = append(testResult.Failures, TestFailure{
			ID:         tf.ID,
			Message:    tf.Message,
			MethodName: tf.MethodName,
			Name:       tf.Name,
			Namespace:  tf.Namespace,
			StackTrace: tf.StackTrace,
			Time:       tf.Time,
			Type:       tf.Type,
		})
	}

	// Convert code coverage
	for _, cc := range xmlTest.CodeCoverage {
		coverage := CodeCoverageResult{
			ID:                     cc.ID,
			Name:                   cc.Name,
			Namespace:              cc.Namespace,
			NumLocations:           cc.NumLocations,
			NumLocationsNotCovered: cc.NumLocationsNotCovered,
			Type:                   cc.Type,
		}

		//nolint:staticcheck // Cannot use type conversion - XML struct has different tags
		for _, loc := range cc.LocationsNotCovered {
			coverage.LocationsNotCovered = append(coverage.LocationsNotCovered, CodeLocation{
				Column:        loc.Column,
				Line:          loc.Line,
				NumExecutions: loc.NumExecutions,
				Time:          loc.Time,
			})
		}

		testResult.CodeCoverage = append(testResult.CodeCoverage, coverage)
	}

	return testResult
}

// CheckDeployStatus checks the status of an asynchronous deployment
func (c *Client) CheckDeployStatus(asyncProcessID string) (*MetadataDeployResult, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	apiVersion := c.getMetadataAPIVersion()
	endpoint := fmt.Sprintf("%s/services/Soap/m/%s", c.InstanceURL, apiVersion)

	// Build checkDeployStatus request
	bodyContent := fmt.Sprintf(`<met:checkDeployStatus>
      <met:asyncProcessId>%s</met:asyncProcessId>
      <met:includeDetails>true</met:includeDetails>
    </met:checkDeployStatus>`, asyncProcessID)

	envelope := c.buildSOAPEnvelope(bodyContent)

	responseBody, err := c.sendSOAPRequest(endpoint, envelope)
	if err != nil {
		return nil, err
	}

	// Parse deploy result using helper types
	var response deployStatusResponse
	if err := xml.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse deploy status: %w", err)
	}

	// Convert to MetadataDeployResult
	result := &MetadataDeployResult{
		CheckOnly:                response.Result.CheckOnly,
		CompletedDate:            response.Result.CompletedDate,
		CreatedDate:              response.Result.CreatedDate,
		Done:                     response.Result.Done,
		ErrorMessage:             response.Result.ErrorMessage,
		ID:                       response.Result.ID,
		IgnoreWarnings:           response.Result.IgnoreWarnings,
		NumberComponentErrors:    response.Result.NumberComponentErrors,
		NumberComponentsDeployed: response.Result.NumberComponentsDeployed,
		NumberComponentsTotal:    response.Result.NumberComponentsTotal,
		NumberTestErrors:         response.Result.NumberTestErrors,
		NumberTestsCompleted:     response.Result.NumberTestsCompleted,
		NumberTestsTotal:         response.Result.NumberTestsTotal,
		RollbackOnError:          response.Result.RollbackOnError,
		RunTestsEnabled:          response.Result.RunTestsEnabled,
		StartDate:                response.Result.StartDate,
		Status:                   response.Result.Status,
		Success:                  response.Result.Success,
	}

	// Convert details if present using helper function
	if len(response.Result.Details.ComponentSuccesses) > 0 ||
		len(response.Result.Details.ComponentFailures) > 0 ||
		response.Result.Details.RunTestResult.NumTestsRun > 0 {
		result.Details = convertDeployDetails(response.Result.Details)
	}

	return result, nil
}

// CancelDeploy cancels an in-progress deployment
func (c *Client) CancelDeploy(asyncProcessID string) (*MetadataAsyncResult, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	apiVersion := c.getMetadataAPIVersion()
	endpoint := fmt.Sprintf("%s/services/Soap/m/%s", c.InstanceURL, apiVersion)

	// Build cancelDeploy request
	bodyContent := fmt.Sprintf(`<met:cancelDeploy>
      <met:asyncProcessId>%s</met:asyncProcessId>
    </met:cancelDeploy>`, asyncProcessID)

	envelope := c.buildSOAPEnvelope(bodyContent)

	responseBody, err := c.sendSOAPRequest(endpoint, envelope)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result struct {
		ID   string `xml:"Body>cancelDeployResponse>result>id"`
		Done bool   `xml:"Body>cancelDeployResponse>result>done"`
	}

	if err := xml.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse cancel deploy response: %w", err)
	}

	return &MetadataAsyncResult{
		ID:    result.ID,
		Done:  result.Done,
		State: "Canceled",
	}, nil
}

// RetrieveMetadata initiates an asynchronous metadata retrieval
func (c *Client) RetrieveMetadata(options MetadataRetrieveOptions) (*MetadataAsyncResult, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	apiVersion := c.getMetadataAPIVersion()
	if options.ApiVersion != "" {
		apiVersion = options.ApiVersion
	}

	endpoint := fmt.Sprintf("%s/services/Soap/m/%s", c.InstanceURL, apiVersion)

	// Strip XML declaration from manifest if present
	manifest := options.UnpackageManifest
	manifest = strings.TrimSpace(manifest)
	if strings.HasPrefix(manifest, "<?xml") {
		// Find the end of the XML declaration and remove it
		if endIdx := strings.Index(manifest, "?>"); endIdx != -1 {
			manifest = strings.TrimSpace(manifest[endIdx+2:])
		}
	}

	// Extract content from Package element (remove <Package> wrapper)
	// The API expects the inner content, not the Package element itself
	packageContent := manifest
	if strings.Contains(manifest, "<Package") {
		// Find the start of Package content
		if startIdx := strings.Index(manifest, ">"); startIdx != -1 {
			content := manifest[startIdx+1:]
			// Find the end tag and remove it
			if endIdx := strings.LastIndex(content, "</Package>"); endIdx != -1 {
				packageContent = strings.TrimSpace(content[:endIdx])
			}
		}
	}

	// Build retrieve request
	bodyContent := fmt.Sprintf(`<met:retrieve>
      <met:retrieveRequest>
        <met:apiVersion>%s</met:apiVersion>
        <met:singlePackage>%t</met:singlePackage>
        <met:unpackaged>%s</met:unpackaged>
      </met:retrieveRequest>
    </met:retrieve>`,
		apiVersion,
		options.SinglePackage,
		packageContent,
	)

	envelope := c.buildSOAPEnvelope(bodyContent)

	responseBody, err := c.sendSOAPRequest(endpoint, envelope)
	if err != nil {
		return nil, err
	}

	var result struct {
		ID string `xml:"Body>retrieveResponse>result>id"`
	}

	if err := xml.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse retrieve response: %w", err)
	}

	return &MetadataAsyncResult{
		ID:    result.ID,
		Done:  false,
		State: "Queued",
	}, nil
}

// CheckRetrieveStatus checks the status of an asynchronous retrieval
func (c *Client) CheckRetrieveStatus(asyncProcessID string) (*MetadataRetrieveResult, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	apiVersion := c.getMetadataAPIVersion()
	endpoint := fmt.Sprintf("%s/services/Soap/m/%s", c.InstanceURL, apiVersion)

	bodyContent := fmt.Sprintf(`<met:checkRetrieveStatus>
      <met:asyncProcessId>%s</met:asyncProcessId>
      <met:includeZip>true</met:includeZip>
    </met:checkRetrieveStatus>`, asyncProcessID)

	envelope := c.buildSOAPEnvelope(bodyContent)

	responseBody, err := c.sendSOAPRequest(endpoint, envelope)
	if err != nil {
		return nil, err
	}

	// Parse retrieve result
	var response struct {
		Result struct {
			Done            bool   `xml:"done"`
			ErrorMessage    string `xml:"errorMessage"`
			ErrorStatusCode string `xml:"errorStatusCode"`
			ID              string `xml:"id"`
			State           string `xml:"state"`
			Status          string `xml:"status"`
			Success         bool   `xml:"success"`
			ZipFile         string `xml:"zipFile"`
			FileProperties  []struct {
				CreatedByID        string `xml:"createdById"`
				CreatedByName      string `xml:"createdByName"`
				CreatedDate        string `xml:"createdDate"`
				FileName           string `xml:"fileName"`
				FullName           string `xml:"fullName"`
				ID                 string `xml:"id"`
				LastModifiedByID   string `xml:"lastModifiedById"`
				LastModifiedByName string `xml:"lastModifiedByName"`
				LastModifiedDate   string `xml:"lastModifiedDate"`
				ManageableState    string `xml:"manageableState"`
				NamespacePrefix    string `xml:"namespacePrefix"`
				Type               string `xml:"type"`
			} `xml:"fileProperties"`
			Messages []struct {
				FileName string `xml:"fileName"`
				Problem  string `xml:"problem"`
			} `xml:"messages"`
		} `xml:"Body>checkRetrieveStatusResponse>result"`
	}

	if err := xml.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse retrieve status: %w", err)
	}

	// Convert to MetadataRetrieveResult
	result := &MetadataRetrieveResult{
		Done:            response.Result.Done,
		ErrorMessage:    response.Result.ErrorMessage,
		ErrorStatusCode: response.Result.ErrorStatusCode,
		ID:              response.Result.ID,
		State:           response.Result.State,
		Status:          response.Result.Status,
		Success:         response.Result.Success,
		ZipFileBase64:   response.Result.ZipFile,
	}

	// Convert file properties
	for _, fp := range response.Result.FileProperties {
		result.FileProperties = append(result.FileProperties, FileProperty{
			CreatedByID:        fp.CreatedByID,
			CreatedByName:      fp.CreatedByName,
			CreatedDate:        fp.CreatedDate,
			FileName:           fp.FileName,
			FullName:           fp.FullName,
			ID:                 fp.ID,
			LastModifiedByID:   fp.LastModifiedByID,
			LastModifiedByName: fp.LastModifiedByName,
			LastModifiedDate:   fp.LastModifiedDate,
			ManageableState:    fp.ManageableState,
			NamespacePrefix:    fp.NamespacePrefix,
			Type:               fp.Type,
		})
	}

	// Convert messages
	for _, msg := range response.Result.Messages {
		result.Messages = append(result.Messages, RetrieveMessage{
			FileName: msg.FileName,
			Problem:  msg.Problem,
		})
	}

	return result, nil
}
