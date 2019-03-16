package paycor

import (
	"crypto/hmac"
	"crypto/sha1" // #nosec G505
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// mac holds keys & handles data signing
type mac struct {
	PublicKey  string
	PrivateKey []byte
}

// Sign generates signature for data
func (mac *mac) Sign(data []byte) (token string, err error) {
	h := hmac.New(sha1.New, mac.PrivateKey)
	_, err = h.Write(data)
	if err != nil {
		return
	}

	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	token = fmt.Sprintf("%s:%s", mac.PublicKey, sign)

	return
}

// SignRequest generates signature for a http.Request
// creates close to an RFC 2104-compliant HMAC signature
// sign these fields:
//   HTTP Method (UPPERCASE)
//   Content MD5 (optional)
//   Content Type (optional)
//   UTC Datetime (RFC 2616)
//   Request Path + Query
// seperated by "\r\n"
func (mac *mac) SignRequest(req *http.Request, httpDate string) (token string, err error) {
	// get to fields that are signed
	method := strings.ToUpper(req.Method)
	u := req.URL
	path := u.Path
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	// create string to sign (need an empty entry at the end to have a trailing \r\n on the data to sign)
	s := strings.Join([]string{method, "", "", httpDate, path, ""}, "\r\n")

	token, err = mac.Sign([]byte(s))
	if err != nil {
		return
	}

	return
}
