package api

import (
	"errors"
)

type responsegetmeetings struct {
	Script     string          `xml:"script"`
	ReturnCode string          `xml:"returncode"`
	Errors     []responseerror `xml:"errors>error"`
	Meetings   []meeting       `xml:"meetings>meeting"`
	MessageKey string          `xml:"messageKey"`
	Message    string          `xml:"message"`
}

type meeting struct {
	MeetingName   string     `xml:"meetingName"`
	MeetingID     string     `xml:"meetingID"`
	InternalID    string     `xml:"internalMeetingID"`
	CreateTime    int64      `xml:"createTime"`
	CreateDate    string     `xml:"createDate"`
	VoiceBridge   int        `xml:"voiceBridge"`
	DialNumber    string     `xml:"dialNumber"`
	AttendeePW    string     `xml:"attendeePW"`
	ModeratorPW   string     `xml:"moderatorPW"`
	Running       bool       `xml:"running"`
	Duration      int        `xml:"duration"`
	HasJoined     bool       `xml:"hasUserJoined"`
	Recording     bool       `xml:"recording"`
	ForciblyEnded bool       `xml:"hasBeenForciblyEnded"`
	StartTime     int64      `xml:"startTime"`
	EndTime       int64      `xml:"endTime"`
	Participants  int        `xml:"participantCount"`
	Listeners     int        `xml:"listenerCount"`
	VoiceCount    int        `xml:"voiceParticipantCount"`
	VideoCount    int        `xml:"videoCount"`
	MaxUsers      int        `xml:"maxUsers"`
	Moderators    int        `xml:"moderatorCount"`
	Attendees     []attendee `xml:"attendees>attendee"`
	Metadata      metadata   `xml:"metadata"`
	IsBreakout    bool       `xml:"isBreakout"`
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
func (api *ApiRequest) GetMeetings() (map[string]meeting, error) {

	//Make the request
	var response responsegetmeetings
	err := api.makeRequest(&response, GET_MEETINGS)
	if err != nil {
		return map[string]meeting{}, err
	}

	//Check if the request was successful
	if response.ReturnCode != "SUCCESS" {
		if response.MessageKey != "" && response.Message != "" {
			return map[string]meeting{}, errors.New(response.MessageKey + ": " + response.Message)
		}
		if response.Errors != nil {
			if response.Errors[0].Key != "" && response.Errors[0].Message != "" {
				return map[string]meeting{}, errors.New(response.Errors[0].Key + ": " + response.Errors[0].Message)
			}
		}
		return map[string]meeting{}, errors.New("API response was not successful")
	}

	// Create map of meetings with InternalID as key
	meetings := map[string]meeting{}
	for _, meeting := range response.Meetings {
		meetings[meeting.MeetingID] = meeting
	}

	return meetings, nil
}
