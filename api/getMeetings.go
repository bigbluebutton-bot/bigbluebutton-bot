package api

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
)


type response struct {
    Script      string   `xml:"script"`
    ReturnCode  string   `xml:"returncode"`
    Meetings    []meeting`xml:"meetings>meeting"`
    MessageKey  string   `xml:"messageKey"`
    Message     string   `xml:"message"`
}

type meeting struct {
	MeetingName   string    `xml:"meetingName"`
	MeetingID     string    `xml:"meetingID"`
	InternalID    string    `xml:"internalMeetingID"`
	CreateTime    int64     `xml:"createTime"`
	CreateDate    string    `xml:"createDate"`
	VoiceBridge   int       `xml:"voiceBridge"`
	DialNumber    string    `xml:"dialNumber"`
	AttendeePW    string    `xml:"attendeePW"`
	ModeratorPW   string    `xml:"moderatorPW"`
	Running       bool      `xml:"running"`
	Duration      int       `xml:"duration"`
	HasJoined     bool      `xml:"hasUserJoined"`
	Recording     bool      `xml:"recording"`
	ForciblyEnded bool      `xml:"hasBeenForciblyEnded"`
	StartTime     int64     `xml:"startTime"`
	EndTime       int64     `xml:"endTime"`
	Participants  int       `xml:"participantCount"`
	Listeners     int       `xml:"listenerCount"`
	VoiceCount    int       `xml:"voiceParticipantCount"`
	VideoCount    int       `xml:"videoCount"`
	MaxUsers      int       `xml:"maxUsers"`
	Moderators    int      	`xml:"moderatorCount"`
	Attendees     []attendee`xml:"attendees>attendee"`
	Metadata 	  metadata 	`xml:"metadata"`
	IsBreakout 	  bool 		`xml:"isBreakout"`
}

type attendee struct {
	UserID         string `xml:"userID"`
	FullName       string `xml:"fullName"`
	Role           string `xml:"role"`
	IsPresenter    bool   `xml:"isPresenter"`
	IsListening    bool   `xml:"isListeningOnly"`
	HasJoinedVoice bool   `xml:"hasJoinedVoice"`
	HasVideo       bool   `xml:"hasVideo"`
	ClientType     string `xml:"clientType"`
}

type metadata struct {
	OriginVersion    string `xml:"bbb-origin-version"`
	OriginServerName string `xml:"bbb-origin-server-name"`
	Origin           string `xml:"bbb-origin"`
	Listed           bool   `xml:"gl-listed"`
}

// Makes a http get request to the BigBlueButton API and returns a list of meetings
func (bbb *api_request) GetMeetings() ([]meeting, error) {
	//Make a http get request to the BigBlueButton API
	resp, err := http.Get(bbb.url + "getMeetings?checksum=" + bbb.generateChecksumSHA256("getMeetings", ""))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//Unmarshal xml
	var r response
	err = xml.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	//Check if the request was successful
	if r.ReturnCode != "SUCCESS" {
		return nil, errors.New(r.Message)
	}

	return r.Meetings, nil
}