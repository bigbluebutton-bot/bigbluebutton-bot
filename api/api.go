package api

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

type sha int64

const (
	SHA1   sha = 0
	SHA256 sha = 1
)

type api_request struct {
	url     string
	secret  string
	shatype sha
}

// Create an object for making http get api requests to the BBB server.
// The requests are described here: https://bigbluebutton.org/api-mate/ and
// https://docs.bigbluebutton.org/dev/api.html
func NewRequest(url string, secret string, shatype sha) (api_request, error) {

	switch shatype {
	case SHA1:
		break
	case SHA256:
		break
	default:
		shatype = SHA256
	}

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
		url:     url,
		secret:  secret,
		shatype: shatype,
	}, nil
}

type params struct {
	name  string
	value string
}

func buildParams(params ...params) string {
	var param string
	for count, p := range params {
		name := strings.ReplaceAll(p.name, " ", "+")
		value := strings.ReplaceAll(p.value, " ", "+")

		if count == 0 {
			param = name + string("=") + value
			continue
		}
		param = param + string("&") + name + string("=") + value
	}
	return param
}

// Generate the checksum for a api request.
// The checksum is generated with the sha1 or sha256 algorithm.
func (api api_request) generateChecksum(action string, params string) string {
	switch api.shatype {
	case SHA1:
		return api.generateChecksumSHA1(action, params)
	case SHA256:
		return api.generateChecksumSHA256(action, params)
	default:
		return ""
	}
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
