package paycor

import (
	"net/http"
	"time"
)

// Transport holds what we need to customize the http client
type transport struct {
	mac       mac
	transport http.RoundTripper
}

// RoundTrip sets Authorization & Date header in request
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	httpDate := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	token, err := t.mac.SignRequest(req, httpDate)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "paycorapi "+token)
	req.Header.Set("Date", httpDate)

	resp, err = t.transport.RoundTrip(req)
	if err != nil {
		return
	}

	return
}
