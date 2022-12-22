package api

import (
	"errors"
	"strings"
)

type api_request struct {
	url string
	secret string
}

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