package api

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

type api_request struct {
	url string
	secret string
}

// Create an object for making http get api requests to the BBB server.
// The requests are described here: https://bigbluebutton.org/api-mate/ and
// https://docs.bigbluebutton.org/dev/api.html
func NewRequest(url string, secret string) (api_request, error) {
	
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		return api_request{}, errors.New("url has the wrong format. It should look like this: https://example.com/api/")
	}


	if !strings.HasSuffix(url, "/") {
		//Add / to the end of the url
		url = url + string("/")
	}
	if !strings.HasSuffix(url, "api/") {
		//Add api/ to the end of the url
		url = url + string("api/")
	}

	if len(secret) != 40 {
		return api_request{}, errors.New("secret must be 40 characters")
	}
	
	return api_request{
		url: url,
		secret: secret,
	}, nil
}

type params struct {
	name string
	value string
}
func buildParams(params ...params) string {
	var param string
	for count, p := range params {
		if count == 0 {
			param = p.name + string("=") + p.value
			continue
		}
		param = param + string("&") + p.name + string("=") + p.value
	}
	return param
}



// Generate the SHA256 checksum for a api request.
func (api api_request) generateChecksumSHA256(action string, params string) string {
	//Generate sha256 and sha1 checksum
	checksum := sha256.New()
	checksum.Write([]byte(action + params + api.secret))
	return hex.EncodeToString(checksum.Sum(nil))
}

// Generate the SHA1 checksum for a api request.
func (api api_request) generateChecksumSHA1(action string, params string) string {
	//Generate sha256 and sha1 checksum
	checksum := sha1.New()
	checksum.Write([]byte(action + params + api.secret))
	return hex.EncodeToString(checksum.Sum(nil))
}