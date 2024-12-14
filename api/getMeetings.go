package api

import (
	"errors"
)

type Responsegetmeetings struct {
	Script     string          `xml:"script" json:"script"`
	ReturnCode string          `xml:"returncode" json:"returnCode"`
	Errors     []responseerror `xml:"errors>error" json:"errors"`
	Meetings   []Meeting       `xml:"meetings>meeting" json:"meetings"`
	MessageKey string          `xml:"messageKey" json:"messageKey"`
	Message    string          `xml:"message" json:"message"`
}

type Meeting struct {
	MeetingName   string     `xml:"meetingName" json:"meetingName"`
	MeetingID     string     `xml:"meetingID" json:"meetingID"`
	InternalID    string     `xml:"internalMeetingID" json:"internalMeetingID"`
	CreateTime    int64      `xml:"createTime" json:"createTime"`
	CreateDate    string     `xml:"createDate" json:"createDate"`
	VoiceBridge   int        `xml:"voiceBridge" json:"voiceBridge"`
	DialNumber    string     `xml:"dialNumber" json:"dialNumber"`
	AttendeePW    string     `xml:"attendeePW" json:"attendeePW"`
	ModeratorPW   string     `xml:"moderatorPW" json:"moderatorPW"`
	Running       bool       `xml:"running" json:"running"`
	Duration      int        `xml:"duration" json:"duration"`
	HasJoined     bool       `xml:"hasUserJoined" json:"hasUserJoined"`
	Recording     bool       `xml:"recording" json:"recording"`
	ForciblyEnded bool       `xml:"hasBeenForciblyEnded" json:"hasBeenForciblyEnded"`
	StartTime     int64      `xml:"startTime" json:"startTime"`
	EndTime       int64      `xml:"endTime" json:"endTime"`
	Participants  int        `xml:"participantCount" json:"participantCount"`
	Listeners     int        `xml:"listenerCount" json:"listenerCount"`
	VoiceCount    int        `xml:"voiceParticipantCount" json:"voiceParticipantCount"`
	VideoCount    int        `xml:"videoCount" json:"videoCount"`
	MaxUsers      int        `xml:"maxUsers" json:"maxUsers"`
	Moderators    int        `xml:"moderatorCount" json:"moderatorCount"`
	Attendees     []Attendee `xml:"attendees>attendee" json:"attendees"`
	Metadata      Metadata   `xml:"metadata" json:"metadata"`
	IsBreakout    bool       `xml:"isBreakout" json:"isBreakout"`
}

type Attendee struct {
	UserID         string `xml:"userID" json:"userID"`
	FullName       string `xml:"fullName" json:"fullName"`
	Role           string `xml:"role" json:"role"`
	IsPresenter    bool   `xml:"isPresenter" json:"isPresenter"`
	IsListening    bool   `xml:"isListeningOnly" json:"isListeningOnly"`
	HasJoinedVoice bool   `xml:"hasJoinedVoice" json:"hasJoinedVoice"`
	HasVideo       bool   `xml:"hasVideo" json:"hasVideo"`
	ClientType     string `xml:"clientType" json:"clientType"`
}

type Metadata struct {
	OriginVersion    string `xml:"bbb-origin-version" json:"bbbOriginVersion"`
	OriginServerName string `xml:"bbb-origin-server-name" json:"bbbOriginServerName"`
	Origin           string `xml:"bbb-origin" json:"bbbOrigin"`
	Listed           bool   `xml:"gl-listed" json:"glListed"`
}

// Makes a http get request to the BigBlueButton API and returns a list of meetings
func (api *ApiRequest) GetMeetings() (map[string]Meeting, error) {

	//Make the request
	var response Responsegetmeetings
	err := api.makeRequest(&response, GET_MEETINGS)
	if err != nil {
		return map[string]Meeting{}, err
	}

	//Check if the request was successful
	if response.ReturnCode != "SUCCESS" {
		if response.MessageKey != "" && response.Message != "" {
			return map[string]Meeting{}, errors.New(response.MessageKey + ": " + response.Message)
		}
		if response.Errors != nil {
			if response.Errors[0].Key != "" && response.Errors[0].Message != "" {
				return map[string]Meeting{}, errors.New(response.Errors[0].Key + ": " + response.Errors[0].Message)
			}
		}
		return map[string]Meeting{}, errors.New("API response was not successful")
	}

	// Create map of meetings with InternalID as key
	meetings := map[string]Meeting{}
	for _, meeting := range response.Meetings {
		meetings[meeting.MeetingID] = meeting
	}

	return meetings, nil
}
