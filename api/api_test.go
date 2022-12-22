package api

import (
	"reflect"
	"testing"
)


type testnewrequest struct {
	url string
	secret string
	expected api_request
	shouldfail bool
}
func TestNewRequest(t * testing.T) {
	tests := []testnewrequest {
		{//0
			url: "https://example.com",
			secret: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			expected: api_request {
				"https://example.com/api/",
				"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			},
			shouldfail: false,
		},
		{//1
			url: "https://example.com/",
			secret: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			expected: api_request {
				"https://example.com/api/",
				"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			},
			shouldfail: false,
		},
		{//2
			url: "http://example.com/",
			secret: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			expected: api_request {
				"http://example.com/api/",
				"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			},
			shouldfail: false,
		},
		{//3
			url: "example.com",
			secret: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			expected: api_request {},
			shouldfail: true,
		},
	}

	for num, test := range tests {
		result, err := NewRequest(test.url, test.secret)

		if err != nil && !test.shouldfail {
			t.Errorf("NewRequest(%s,%s) %b FAILED: Error %s", test.url, test.secret, num, err)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("NewRequest(%s,%s) %b FAILED: Object is not correct", test.url, test.secret, num)
		} else {
			t.Logf("NewRequest(%s,%s) %b PASSED", test.url, test.secret, num)
		}
	}
}