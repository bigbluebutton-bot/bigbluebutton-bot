package api

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
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



// Only those actions are allowed
type action string
const (
	CREATE 						action = "create"
	JOIN   						action = "join"
	IS_MEETING_RUNNING 			action = "isMeetingRunning"
	GET_MEETING_INFO 			action = "getMeetingInfo"
	MEETING 					action = "meeting"
	GET_MEETINGS 				action = "getMeetings"
	GET_DEFAULT_CONFIG_XML 		action = "getDefaultConfigXML"
	SET_CONFIG_XML 				action = "setConfigXML"
	ENTER 						action = "enter"
	CONFIG_XML 					action = "configXML"
	SIGN_OUT 					action = "signOut"
	GET_RECORDINGS 				action = "getRecordings"
	PUBLISH_RECORDINGS 			action = "publishRecordings"
	DELETE_RECORDINGS 			action = "deleteRecordings"
	UPDATE_RECORDINGS 			action = "updateRecordings"
	GET_RECORDING_TEXT_TRACKS 	action = "getRecordingTextTracks"
)

// Only those parames are allowed
type paramname string
const (
	MEETING_ID 					paramname = "meetingID"
	RECORD_ID 					paramname = "recordID"
	NAME 						paramname = "name"
	ATTENDEE_PW 				paramname = "attendeePW"
	MODERATOR_PW 				paramname = "moderatorPW"
	PASSWORD 					paramname = "password"	//same as moderatorPW (I dont know why its sometimse called password and not moderatorPW)
	FULL_NAME 					paramname = "fullName"
	WELCOME 					paramname = "welcome"
	VOICE_BRIDGE 				paramname = "voiceBridge"
	RECORD 						paramname = "record"
	AUTO_START_RECORDING 		paramname = "autoStartRecording"
	ALLOW_START_STOP_RECORDING 	paramname = "allowStartStopRecording"
	DIAL_NUMBER 				paramname = "dialNumber"
	WEB_VOICE 					paramname = "webVoice"
	LOGOUT_URL 					paramname = "logoutURL"
	MAX_PARTICIPANTS 			paramname = "maxParticipants"
	DURATION 					paramname = "duration"
	USER_ID 					paramname = "userID"
	CREATE_TIME 				paramname = "createTime"
	WEB_VOICE_CONF 				paramname = "webVoiceConf"
	PUBLISH 					paramname = "publish"
	REDIRECT 					paramname = "redirect"
	CLIENT_URL 					paramname = "clientURL"
	CONFIG_TOKEN 				paramname = "configToken"
	AVATAR_URL 					paramname = "avatarURL"
	MODERATOR_ONLY_MESSAGE 		paramname = "moderatorOnlyMessage"
)


type params struct {
	name  paramname
	value string
}

func buildParams(params ...params) string {
	var param string
	for count, p := range params {
		name := strings.ReplaceAll(string(p.name), " ", "+")
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
func (api api_request) generateChecksum(action action, params string) string {
	if(api.shatype == SHA1) {
		return api.generateChecksumSHA1(action, params)
	} else {
		return api.generateChecksumSHA256(action, params)
	}
}

// Generate the SHA256 checksum for a api request.
func (api api_request) generateChecksumSHA256(action action, params string) string {
	//Generate sha256 and sha1 checksum
	checksum := sha256.New()
	checksum.Write([]byte(string(action) + params + api.secret))
	return hex.EncodeToString(checksum.Sum(nil))
}

// Generate the SHA1 checksum for a api request.
func (api api_request) generateChecksumSHA1(action action, params string) string {
	//Generate sha256 and sha1 checksum
	checksum := sha1.New()
	checksum.Write([]byte(string(action) + params + api.secret))
	return hex.EncodeToString(checksum.Sum(nil))
}

func (api api_request) makeRequest(response any, action action, params ...params) (error) {
	param := buildParams(params...)
	checksum := api.generateChecksum(action, param)

	var url string
	if(len([]rune(param)) > 0) {
		url = api.url + string(action) + string("?") + param + string("&checksum=") + checksum
	} else {
		url = api.url + string(action) + string("?checksum=") + checksum
	}

	//Make a http get request to the BigBlueButton API
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//Unmarshal xml
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	return nil
}