package api

import (
	"errors"
	"fmt"
	"strconv"
)

type responseCreateMeeting struct {
	Script          string `xml:"script"`
	ReturnCode      string `xml:"returncode"`
	Errors   		[]responseerror `xml:"errors>error"`
	MeetingID       string `xml:"meetingID"`
	InternalMeeting string `xml:"internalMeetingID"`
	ParentMeeting   string `xml:"parentMeetingID"`
	AttendeePW      string `xml:"attendeePW"`
	ModeratorPW     string `xml:"moderatorPW"`
	CreateTime      int64  `xml:"createTime"`
	VoiceBridge     int64  `xml:"voiceBridge"`
	DialNumber      string `xml:"dialNumber"`
	CreateDate      string `xml:"createDate"`
	HasUserJoined   bool   `xml:"hasUserJoined"`
	Duration        int64  `xml:"duration"`
	HasBeenEnded    bool   `xml:"hasBeenForciblyEnded"`
	MessageKey      string `xml:"messageKey"`
	Message         string `xml:"message"`
}


// Makes a http get request to the BigBlueButton API and returns a list of meetings
func (api *api_request) CreateMeeting(name string, meetingID string, attendeePW string, moderatorPW string, welcome string, allowStartStopRecording bool, autoStartRecording bool, record bool, voiceBridge int64) (meeting, error) {

	params := []params{
		{
			name:  ALLOW_START_STOP_RECORDING,
			value: strconv.FormatBool(allowStartStopRecording),
		},
		{
			name:  ATTENDEE_PW,
			value: attendeePW,
		},
		{
			name:  AUTO_START_RECORDING,
			value: strconv.FormatBool(autoStartRecording),
		},
		{
			name:  MEETING_ID,
			value: meetingID,
		},
		{
			name:  MODERATOR_PW,
			value: moderatorPW,
		},
		{
			name:  NAME,
			value: name,
		},
		{
			name:  RECORD,
			value: strconv.FormatBool(record),
		},
		{
			name:  VOICE_BRIDGE,
			value: strconv.FormatInt(voiceBridge, 10),
		},
		{
			name:  WELCOME,
			value: welcome,
		},
	}

	//Make the request
	var response responseCreateMeeting
	err := api.makeRequest(&response, CREATE, params...)
	if err != nil {
		return meeting{}, err
	}

	//Check if the request was successful
	if response.ReturnCode != "SUCCESS" {
		if(response.MessageKey != "" && response.Message != "") {
			return meeting{}, errors.New(response.MessageKey + ": " + response.Message)
		}
		if(response.Errors != nil) {
			if(response.Errors[0].Key != "" && response.Errors[0].Message != "") {
				return meeting{}, errors.New(response.Errors[0].Key + ": " + response.Errors[0].Message)
			}
		}
		return meeting{}, errors.New("API response was not successful")
	}

	//Get the meeting info
	meetings, err := api.GetMeetings()
	if(err != nil) {
		return meeting{}, err
	}

	// Check if meeting already exists (duplicateWarning)
	if response.MessageKey == "duplicateWarning" {
		fmt.Println("Warning: Meeting already exists")
	}

	return meetings[response.MeetingID], nil
}