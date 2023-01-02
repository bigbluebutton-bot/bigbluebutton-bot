package api

import (
	"fmt"
	"reflect"
	"testing"
)

type testnewrequest struct {
	url        string
	secret     string
	shatype    sha
	expected   api_request
	shouldfail bool
}

func TestNewRequest(t *testing.T) {
	tests := []testnewrequest{
		{ //0
			url:     "https://example.com",
			secret:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			shatype: SHA256,
			expected: api_request{
				url:     "https://example.com/api/",
				secret:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				shatype: SHA256,
			},
			shouldfail: false,
		},
		{ //1
			url:     "https://example.com/",
			secret:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			shatype: SHA1,
			expected: api_request{
				url:     "https://example.com/api/",
				secret:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				shatype: SHA1,
			},
			shouldfail: false,
		},
		{ //2
			url:     "http://example.com/",
			secret:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			shatype: SHA256,
			expected: api_request{
				url:     "http://example.com/api/",
				secret:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				shatype: SHA256,
			},
			shouldfail: false,
		},
		{ //3
			url:        "example.com",
			secret:     "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			expected:   api_request{},
			shouldfail: true,
		},
		{ //4
			url:     "https://example.com",
			secret:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			shatype: 2,
			expected: api_request{
				url:     "https://example.com/api/",
				secret:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				shatype: SHA256,
			},
			shouldfail: false,
		},
		{ //5
			url:        "http://example.com",
			secret:     "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			shatype:    SHA256,
			expected:   api_request{},
			shouldfail: true,
		},
	}

	for num, test := range tests {
		result, err := NewRequest(test.url, test.secret, test.shatype)
		fmt.Println(test.shatype)

		if err != nil && !test.shouldfail {
			t.Errorf("NewRequest(%s,%s) %d FAILED: Error %s", test.url, test.secret, num, err)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("NewRequest(%s,%s) %d FAILED: Object is not correct", test.url, test.secret, num)
		} else {
			t.Logf("NewRequest(%s,%s) %d PASSED", test.url, test.secret, num)
		}
	}
}

type testgeneratechecksum struct {
	action          action
	params          []params
	expected_sha1   string
	expected_sha256 string
	shouldfail      bool
}

// Test for generateChecksum
// https://mconf.github.io/api-mate/
func TestGenerateChecksum(t *testing.T) {
	tests := []testgeneratechecksum{
		{ //0
			action: CREATE,
			params: []params{
				{
					name:  ALLOW_START_STOP_RECORDING,
					value: "true",
				},
				{
					name:  ATTENDEE_PW,
					value: "ap",
				},
				{
					name:  AUTO_START_RECORDING,
					value: "false",
				},
				{
					name:  MEETING_ID,
					value: "random-4026116",
				},
				{
					name:  MODERATOR_PW,
					value: "mp",
				},
				{
					name:  NAME,
					value: "random-4026116",
				},
				{
					name:  RECORD,
					value: "false",
				},
				{
					name:  VOICE_BRIDGE,
					value: "70848",
				},
				{
					name:  WELCOME,
					value: "Hello you there",
				},
			},
			expected_sha1:   "2c2f2b2f6050bda0ff2c6dacd9d51e09951810ae",
			expected_sha256: "ae982d76751077e4e1eae8a667d5f74fe4f9c9a9df7d30ff2e56b3a025f1828d",
			shouldfail:      false,
		},
	}

	bbbapi_sha1, _ := NewRequest("https://example.com/bigbluebutton/api/", "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", SHA1)
	bbbapi_sha256, _ := NewRequest("https://example.com/bigbluebutton/api/", "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", SHA256)

	for num, test := range tests {
		params := buildParams(test.params...)

		//Sha1
		resultsha1 := bbbapi_sha1.generateChecksum(test.action, params)
		if resultsha1 != test.expected_sha1 {
			if test.shouldfail {
				t.Logf("generateChecksumSHA1(%s,...) %d PASSED", test.action, num)
			} else {
				t.Errorf("generateChecksumSHA1(%s,...) %d FAILED: Cheksum is wrong: %s", test.action, num, bbbapi_sha1.url+string(test.action)+"?"+params+"&checksum="+resultsha1)
			}
		} else {
			t.Logf("generateChecksumSHA1(%s,...) %d PASSED", test.action, num)
		}

		//Sha256
		resultsha256 := bbbapi_sha256.generateChecksum(test.action, params)
		if resultsha256 != test.expected_sha256 {
			if test.shouldfail {
				t.Logf("generateChecksumSHA256(%s,...) %d PASSED", test.action, num)
			} else {
				t.Errorf("generateChecksumSHA256(%s,...) %d FAILED: Cheksum is wrong: %s", test.action, num, bbbapi_sha256.url+string(test.action)+"?"+params+"&checksum="+resultsha256)
			}
		} else {
			t.Logf("generateChecksumSHA256(%s,...) %d PASSED", test.action, num)
		}
	}
}
