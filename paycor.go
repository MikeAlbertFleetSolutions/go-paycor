package paycor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client is our type
type Client struct {
	host       string
	httpclient *http.Client
}

// NewClient creates new paycor client
func NewClient(publicKey, privateKey, paycorHost string) (client *Client) {
	client = &Client{
		host: paycorHost,
		httpclient: &http.Client{
			Transport: &transport{
				transport: http.DefaultTransport,
				mac: mac{
					PublicKey:  publicKey,
					PrivateKey: []byte(privateKey),
				},
			},
		},
	}

	return
}

// makeRequest is a helper function to wrap making REST calls to paycor
func (client *Client) makeRequest(method, url string, parameters map[string]string) (data []byte, err error) {
	// create request
	var request *http.Request
	request, err = http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	// add query parameters
	if parameters != nil {
		q := request.URL.Query()
		for k, v := range parameters {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}

	// make request, get response
	var response *http.Response
	response, err = client.httpclient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// error?
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s call to %s returned status code %d ", method, url, response.StatusCode)
		return
	}

	// get body for caller, if there is something
	if response.ContentLength != 0 {
		data, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
	}

	return
}

// GetDocumentTypes returns the complete list of document types that the caller can request
func (client *Client) GetDocumentTypes() (results []string, err error) {
	var data []byte
	data, err = client.makeRequest("GET", fmt.Sprintf("https://%s/documents/api/documenttypes", client.host), nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return
	}

	return
}

// GetDocumentListing returns the list of available documents of the spcified DocumentType
func (client *Client) GetDocumentListing(documentType string) (results map[string]interface{}, err error) {
	var data []byte
	data, err = client.makeRequest("GET", fmt.Sprintf("https://%s/documents/api/documents/%s", client.host, documentType), nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return
	}

	return
}

// GetLiveOrSavedReport retrieves the report identified by documentType & documentID
func (client *Client) GetLiveOrSavedReport(documentType, documentID string) (results []byte, err error) {
	results, err = client.makeRequest("GET", fmt.Sprintf("https://%s/documents/api/documents/%s/%s", client.host, documentType, documentID), nil)

	return
}

// GetReportByName retrieves the report identified by reportName
func (client *Client) GetReportByName(reportName string) (results []byte, err error) {
	reportType := "customreport"

	// get all documents of type reportType
	var docs map[string]interface{}
	docs, err = client.GetDocumentListing(reportType)
	if err != nil {
		return
	}

	// find report caller wants
	var reportID string
	items := docs["Items"].([]interface{})
	for _, item := range items {
		i := item.(map[string]interface{})
		if i["DocumentName"] == reportName {
			reportID = i["Id"].(string)
			break
		}
	}

	// nope
	if reportID == "" {
		err = fmt.Errorf("could not find report %s", reportName)
		return
	}

	// retrieve
	results, err = client.GetLiveOrSavedReport(reportType, reportID)

	return
}
