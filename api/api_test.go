package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)


type config struct {
	Url    string `json:"url"`
	Secret string `json:"secret"`
}
// For reading config from a file or from environment variables
func readConfig(file string, t *testing.T) config {

	// we initialize config
	var conf config

	// Try to read from env
	conf.Url = os.Getenv("BBB_URL")
	conf.Secret = os.Getenv("BBB_SECRET")

	if(conf.Url != "" && conf.Secret != "") {
		t.Log("Using env variables for config")
		return conf
	}

	// Try to read from config file
	t.Log("Using config file for config")
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our jsonFile's content into conf
	json.Unmarshal([]byte(byteValue), &conf)

	return conf
}






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


type testmakeRequest struct {
	url			 string
	secret		 string
	action       action
	params       []params
	expected   	 any
	shouldfail   bool
}
// Test for makeRequest
func TestMakeRequest(t *testing.T) {
	tests := []testmakeRequest{
		{ //0
			url: "https://examfgfgfgfffple.com/bigbluebutton/api/",
			secret: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			action: GET_MEETINGS,
			params: []params{},
			expected: "",
			shouldfail: true,
		},
		{ //1
			url: "https://test-install.blindsidenetworks.com/bigbluebutton/api/",
			secret: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			action: GET_MEETINGS,
			params: []params{},
			expected: "",
			shouldfail: true,
		},
		{ //2
			url: "https://test-install.blindsidenetworks.com/bigbluebutton/api/",
			secret: "8cd8ef52e8e101574e400365b55e11a6",
			action: GET_MEETINGS,
			params: []params{},
			expected: "",
			shouldfail: false,
		},
		{ //3
			url: "https://test-install.blindsidenetworks.com/wrong/api/",
			secret: "8cd8ef52e8e101574e400365b55e11a6",
			action: GET_MEETINGS,
			params: []params{},
			expected: "",
			shouldfail: true,
		},
	}

	

	for num, test := range tests {
		bbbapi, err := NewRequest(test.url, test.secret, SHA1)
		if(err != nil) {
			t.Errorf("makeRequest(...,%s,...) %d FAILED: NewRequest: %s", test.action, num, err)
			continue
		}

		var response responsegetmeetings
		err = bbbapi.makeRequest(&response, test.action, test.params...)
		if(err != nil) {
			if(!test.shouldfail) {
				t.Errorf("makeRequest(...,%s,...) %d FAILED: err: %s", test.action, num, err)
				continue
			}
		}
		t.Logf("makeRequest(...,%s,...) %d PASSED", test.action, num)
	}
}