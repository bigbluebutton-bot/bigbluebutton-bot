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
	}

	for num, test := range tests {
		result, err := NewRequest(test.url, test.secret, test.shatype)
		fmt.Println(test.shatype)

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

type testgeneratechecksum struct {
	action          string
	params          []params
	expected_sha1   string
	expected_sha256 string
	shouldfail      bool
}

// Test for generateChecksum
func TestGenerateChecksum(t *testing.T) {
	tests := []testgeneratechecksum{
		{ //0 https://example.com/api/create?allowStartStopRecording=true&attendeePW=ap&autoStartRecording=false&meetingID=random-4026116&moderatorPW=mp&name=random-4026116&record=false&voiceBridge=70848&welcome=%3Cbr%3EWelcome+to+%3Cb%3E%25%25CONFNAME%25%25%3C%2Fb%3E%21&checksum=420805f07ee9e1b75537bcc22cea3586ec01be82565843ca83bfc3625aeac0ad
			action: "create",
			params: []params{
				{
					name:  "allowStartStopRecording",
					value: "true",
				},
				{
					name:  "attendeePW",
					value: "ap",
				},
				{
					name:  "autoStartRecording",
					value: "false",
				},
				{
					name:  "meetingID",
					value: "random-4026116",
				},
				{
					name:  "moderatorPW",
					value: "mp",
				},
				{
					name:  "name",
					value: "random-4026116",
				},
				{
					name:  "record",
					value: "false",
				},
				{
					name:  "voiceBridge",
					value: "70848",
				},
				{
					name:  "welcome",
					value: "Hello",
				},
			},
			expected_sha1:   "e9a7a7a51c8dfc1b53ab313b096307215325fc15",
			expected_sha256: "e29adcdff42ebe39f9e3f3a5d528d798cce15442f6bb9115f9fe288a482f16b2",
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
				t.Logf("generateChecksumSHA1(%s,...) %b PASSED", test.action, num)
			} else {
				t.Errorf("generateChecksumSHA1(%s,...) %b FAILED: Cheksum is wrong: %s", test.action, num, bbbapi_sha1.url+test.action+"?"+params+"&checksum="+resultsha1)
			}
		} else {
			t.Logf("generateChecksumSHA1(%s,...) %b PASSED", test.action, num)
		}

		//Sha256
		resultsha256 := bbbapi_sha256.generateChecksum(test.action, params)
		if resultsha256 != test.expected_sha256 {
			if test.shouldfail {
				t.Logf("generateChecksumSHA256(%s,...) %b PASSED", test.action, num)
			} else {
				t.Errorf("generateChecksumSHA256(%s,...) %b FAILED: Cheksum is wrong: %s", test.action, num, bbbapi_sha256.url+test.action+"?"+params+"&checksum="+resultsha256)
			}
		} else {
			t.Logf("generateChecksumSHA256(%s,...) %b PASSED", test.action, num)
		}
	}
}
