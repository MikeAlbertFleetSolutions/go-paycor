package paycor

import (
	"net/http"
	"reflect"
	"testing"
)

var (
	validPaycorHost = "secure.paycor.com"
	validPublicKey  = "REPLACE"
	validPrivateKey = "REPLACE"
	validReportID   = "REPLACE"
	validReportName = "REPLACE"
)

func TestNewClient(t *testing.T) {
	type args struct {
		publicKey  string
		privateKey string
		paycorHost string
	}
	tests := []struct {
		name       string
		args       args
		wantClient *Client
	}{
		{
			"test persisting host & keys",
			args{"test-public-key", "test-private-key", "test.paycor.com"},
			&Client{
				host: "test.paycor.com",
				httpclient: &http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  "test-public-key",
							PrivateKey: []byte("test-private-key"),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotClient := NewClient(tt.args.publicKey, tt.args.privateKey, tt.args.paycorHost); !reflect.DeepEqual(gotClient, tt.wantClient) {
				t.Errorf("NewClient() = %v, want %v", gotClient, tt.wantClient)
			}
		})
	}
}

func TestClient_GetDocumentTypes(t *testing.T) {
	type fields struct {
		host       string
		httpclient *http.Client
	}
	tests := []struct {
		name        string
		fields      fields
		wantResults bool
		wantErr     bool
	}{
		{
			"Retrieve document types from paycor with valid public & private keys",
			fields{
				validPaycorHost,
				&http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  validPublicKey,
							PrivateKey: []byte(validPrivateKey),
						},
					},
				},
			},
			true,
			false,
		},
		// Apparantly paycor doesn't care about securing this call?
		// {
		// 	"Retrieve document types from paycor with invalid public & private keys",
		// 	fields{
		// 		validPaycorHost,
		// 		&http.Client{
		// 			Transport: &transport{
		// 				transport: http.DefaultTransport,
		// 				mac: mac{
		// 					PublicKey:  "badkey",
		// 					PrivateKey: []byte("badkey"),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	false,
		// 	true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				host:       tt.fields.host,
				httpclient: tt.fields.httpclient,
			}
			gotResults, err := client.GetDocumentTypes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetDocumentTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (len(gotResults) > 0) != tt.wantResults {
				t.Errorf("Client.GetDocumentTypes() = %v, wantResults %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestClient_GetDocumentListing(t *testing.T) {
	type fields struct {
		host       string
		httpclient *http.Client
	}
	type args struct {
		documentType string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults bool
		wantErr     bool
	}{
		{
			"Retrieve document list from paycor with valid document type",
			fields{
				validPaycorHost,
				&http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  validPublicKey,
							PrivateKey: []byte(validPrivateKey),
						},
					},
				},
			},
			args{"customreport"},
			true,
			false,
		},
		{
			"Retrieve document list from paycor with invalid document type",
			fields{
				validPaycorHost,
				&http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  validPublicKey,
							PrivateKey: []byte(validPrivateKey),
						},
					},
				},
			},
			args{"bad-type"},
			false,
			true,
		},
		{
			"Retrieve document list from paycor with invalid public & private keys",
			fields{
				validPaycorHost,
				&http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  "badkey",
							PrivateKey: []byte("badkey"),
						},
					},
				},
			},
			args{"bad-type"},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				host:       tt.fields.host,
				httpclient: tt.fields.httpclient,
			}
			gotResults, err := client.GetDocumentListing(tt.args.documentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetDocumentListing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (len(gotResults) > 0) != tt.wantResults {
				t.Errorf("Client.GetDocumentListing() = %v, wantResults %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestClient_GetLiveOrSavedReport(t *testing.T) {
	type fields struct {
		host       string
		httpclient *http.Client
	}
	type args struct {
		documentType string
		documentID   string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults bool
		wantErr     bool
	}{
		{
			"Retrieve report from paycor with valid document type & ID",
			fields{
				validPaycorHost,
				&http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  validPublicKey,
							PrivateKey: []byte(validPrivateKey),
						},
					},
				},
			},
			args{"customreport", validReportID},
			true,
			false,
		},
		{
			"Retrieve report from paycor with invalid document type & ID",
			fields{
				validPaycorHost,
				&http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  validPublicKey,
							PrivateKey: []byte(validPrivateKey),
						},
					},
				},
			},
			args{"bad-report", "bad-id"},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				host:       tt.fields.host,
				httpclient: tt.fields.httpclient,
			}
			gotResults, err := client.GetLiveOrSavedReport(tt.args.documentType, tt.args.documentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLiveOrSavedReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (len(gotResults) > 0) != tt.wantResults {
				t.Errorf("Client.GetDocumentListing() = %v, wantResults %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestClient_GetReportByName(t *testing.T) {
	type fields struct {
		host       string
		httpclient *http.Client
	}
	type args struct {
		documentName string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults bool
		wantErr     bool
	}{
		{
			"Retrieve report from paycor with valid name",
			fields{
				validPaycorHost,
				&http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  validPublicKey,
							PrivateKey: []byte(validPrivateKey),
						},
					},
				},
			},
			args{validReportName},
			true,
			false,
		},
		{
			"Retrieve report from paycor with invalid name",
			fields{
				validPaycorHost,
				&http.Client{
					Transport: &transport{
						transport: http.DefaultTransport,
						mac: mac{
							PublicKey:  validPublicKey,
							PrivateKey: []byte(validPrivateKey),
						},
					},
				},
			},
			args{"bad-name"},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				host:       tt.fields.host,
				httpclient: tt.fields.httpclient,
			}
			gotResults, err := client.GetReportByName(tt.args.documentName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetReportByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (len(gotResults) > 0) != tt.wantResults {
				t.Errorf("Client.GetDocumentListing() = %v, wantResults %v", gotResults, tt.wantResults)
			}
		})
	}
}
