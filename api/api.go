package api

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type SHA string

const (
	SHA1   SHA = "SHA1"
	SHA256 SHA = "SHA256"
)

type ApiRequest struct {
	url     string
	secret  string
	shatype SHA
}

// Create an object for making http get api requests to the BBB server.
// The requests are described here: https://bigbluebutton.org/api-mate/ and
// https://docs.bigbluebutton.org/dev/api.html
func NewRequest(url string, secret string, shatype SHA) (*ApiRequest, error) {

	switch shatype {
	case SHA1:
		break
	case SHA256:
		break
	default:
		shatype = SHA256
	}

	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		return &ApiRequest{}, errors.New("url has the wrong format. It should look like this: https://example.com/api/")
	}

	if !strings.HasSuffix(url, "/") {
		//Add / to the end of the url
		url = url + string("/")
	}
	if !strings.HasSuffix(url, "api/") {
		//Add api/ to the end of the url
		url = url + string("api/")
	}

	return &ApiRequest{
		url:     url,
		secret:  secret,
		shatype: shatype,
	}, nil
}

// Only those actions are allowed
type action string

const (
	CREATE             action = "create"
	END                action = "end"
	GET_MEETINGS       action = "getMeetings"
	IS_MEETING_RUNNING action = "isMeetingRunning"
	JOIN               action = "join"

	// GET_RECORDINGS 				action = "getRecordings"
	// PUBLISH_RECORDINGS 			action = "publishRecordings"
	// DELETE_RECORDINGS 			action = "deleteRecordings"
	// UPDATE_RECORDINGS 			action = "updateRecordings"
	// GET_RECORDING_TEXT_TRACKS 	action = "getRecordingTextTracks"

	// GET_MEETING_INFO 			action = "getMeetingInfo"
	// GET_DEFAULT_CONFIG_XML 		action = "getDefaultConfigXML"
	// SET_CONFIG_XML 				action = "setConfigXML"
	// ENTER 						action = "enter"
	// CONFIG_XML 					action = "configXML"
	// SIGN_OUT 					action = "signOut"
)

// Only those parames are allowed
type paramname string

const (
	MEETING_ID                 paramname = "meetingID"
	RECORD_ID                  paramname = "recordID"
	NAME                       paramname = "name"
	ATTENDEE_PW                paramname = "attendeePW"
	MODERATOR_PW               paramname = "moderatorPW"
	PASSWORD                   paramname = "password" //same as moderatorPW (I dont know why its sometimse called password and not moderatorPW)
	FULL_NAME                  paramname = "fullName"
	WELCOME                    paramname = "welcome"
	VOICE_BRIDGE               paramname = "voiceBridge"
	RECORD                     paramname = "record"
	AUTO_START_RECORDING       paramname = "autoStartRecording"
	ALLOW_START_STOP_RECORDING paramname = "allowStartStopRecording"
	DIAL_NUMBER                paramname = "dialNumber"
	WEB_VOICE                  paramname = "webVoice"
	LOGOUT_URL                 paramname = "logoutURL"
	MAX_PARTICIPANTS           paramname = "maxParticipants"
	DURATION                   paramname = "duration"
	USER_ID                    paramname = "userID"
	CREATE_TIME                paramname = "createTime"
	WEB_VOICE_CONF             paramname = "webVoiceConf"
	PUBLISH                    paramname = "publish"
	REDIRECT                   paramname = "redirect"
	CLIENT_URL                 paramname = "clientURL"
	CONFIG_TOKEN               paramname = "configToken"
	AVATAR_URL                 paramname = "avatarURL"
	MODERATOR_ONLY_MESSAGE     paramname = "moderatorOnlyMessage"
)

type params struct {
	name  paramname
	value string
}

func (api *ApiRequest) buildParams(params ...params) string {
	var param string
	for count, p := range params {

		//Replace special chars
		name := url.QueryEscape(string(p.name))
		value := url.QueryEscape(p.value)

		if count == 0 {
			param = name + string("=") + value
			continue
		}
		param = param + string("&") + name + string("=") + value
	}

	//Replace some chars with origanal char
	param = strings.ReplaceAll(param, url.QueryEscape(" "), "+")

	return param
}

// Generate the checksum for a api request.
// The checksum is generated with the sha1 or sha256 algorithm.
func (api *ApiRequest) generateChecksum(action action, params string) string {
	if api.shatype == SHA1 {
		return api.generateChecksumSHA1(action, params)
	} else {
		return api.generateChecksumSHA256(action, params)
	}
}

// Generate the SHA256 checksum for a api request.
func (api ApiRequest) generateChecksumSHA256(action action, params string) string {
	//Generate sha256 and sha1 checksum
	checksum := sha256.New()
	checksum.Write([]byte(string(action) + params + api.secret))
	return hex.EncodeToString(checksum.Sum(nil))
}

// Generate the SHA1 checksum for a api request.
func (api *ApiRequest) generateChecksumSHA1(action action, params string) string {
	//Generate sha256 and sha1 checksum
	checksum := sha1.New()
	checksum.Write([]byte(string(action) + params + api.secret))
	return hex.EncodeToString(checksum.Sum(nil))
}

// The response from the BigBlueButton API
// EXAMPLES:
// type response struct {
//     Script      string   `xml:"script"`
//     ReturnCode  string   `xml:"returncode"`
// 	   Errors   	[]responseerror `xml:"errors>error"`
//     MessageKey  string   `xml:"messageKey"`
//     Message     string   `xml:"message"`
// }

type responseerror struct {
	Key     string `xml:"key"`
	Message string `xml:"message"`
}

func (api *ApiRequest) buildURL(action action, params ...params) string {
	param := api.buildParams(params...)
	checksum := api.generateChecksum(action, param)

	var url string
	if len([]rune(param)) > 0 {
		url = api.url + string(action) + string("?") + param + string("&checksum=") + checksum
	} else {
		url = api.url + string(action) + string("?checksum=") + checksum
	}

	return url
}

func (api *ApiRequest) makeRequest(response any, action action, params ...params) error {

	url := api.buildURL(action, params...)

	//Make a http get request to the BigBlueButton API
	client := new(http.Client)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req) //send request
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Server returned: " + resp.Status)
	}

	cookies := resp.Cookies() //get cookies

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//Unmarshal xml
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	// Set Cookie to response.Cookie
	ps := reflect.ValueOf(response)
	// struct
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		// exported field
		f := s.FieldByName("Cookie")
		if f.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if f.CanSet() {
				// change value of Cookie
				if f.Kind() == reflect.Slice {
					//Set Cookie
					f.Set(reflect.ValueOf(cookies))
				}
			}
		}
	}

	return nil
}
