package paycor

import (
	"net/http"
	"net/url"
	"testing"
)

func Test_mac_Sign(t *testing.T) {
	type fields struct {
		PublicKey  string
		PrivateKey []byte
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			"Validate that data signing is correct",
			fields{
				"public-key",
				[]byte("private-key"),
			},
			args{[]byte("data-to_sign")},
			"public-key:zSiy832GGzdcnhDwfmmGmY7N1v0=",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mac := &mac{
				PublicKey:  tt.fields.PublicKey,
				PrivateKey: tt.fields.PrivateKey,
			}
			gotToken, err := mac.Sign(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("mac.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("mac.Sign() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func Test_mac_SignRequest(t *testing.T) {
	type fields struct {
		PublicKey  string
		PrivateKey []byte
	}
	type args struct {
		req      *http.Request
		httpDate string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			"Validate that requests are signing correctly",
			fields{
				"public-key",
				[]byte("private-key"),
			},
			args{
				&http.Request{
					Method: "GET",
					URL: &url.URL{
						Path:     "/testendpoint",
						RawQuery: "testing=true",
					},
				},
				"Mon, 02 Jan 2006 15:04:05 GMT",
			},
			"public-key:xMCNkZz+EPsUh5TOCyCdZq2BVb8=",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mac := &mac{
				PublicKey:  tt.fields.PublicKey,
				PrivateKey: tt.fields.PrivateKey,
			}
			gotToken, err := mac.SignRequest(tt.args.req, tt.args.httpDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("mac.SignRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("mac.SignRequest() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
