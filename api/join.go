package api

import (
	"net/http"
)

type responseJoin struct {
	Script       string          `xml:"script"`
	ReturnCode   string          `xml:"returncode"`
	Errors       []responseerror `xml:"errors>error"`
	MessageKey   string          `xml:"messageKey"`
	Message      string          `xml:"message"`
	MeetingID    string          `xml:"meeting_id"`
	UserID       string          `xml:"user_id"`
	AuthToken    string          `xml:"auth_token"`
	SessionToken string          `xml:"session_token"`
	GuestStatus  string          `xml:"guestStatus"`
	URL          string          `xml:"url"`
	Cookie       []*http.Cookie
}

// Makes a http get request to the BigBlueButton API to join a meeting and returs:
// - url
// - cookie
// - userid
// - auth_token
// - session_token
func (api *ApiRequest) Join(meetingID string, userName string, moderator bool) (string, []*http.Cookie, string, string, string, error) {

	meetings, err := api.GetMeetings()
	if err != nil {
		return "", nil, "", "", "", err
	}

	m := meetings[meetingID]

	var password string
	if moderator {
		password = m.ModeratorPW
	} else {
		password = m.AttendeePW
	}

	params := []params{
		{
			name:  FULL_NAME,
			value: userName,
		},
		{
			name:  MEETING_ID,
			value: meetingID,
		},
		{
			name:  PASSWORD,
			value: password,
		},
		{
			name:  REDIRECT,
			value: "false",
		},
	}

	var response responseJoin
	err = api.makeRequest(&response, JOIN, params...)
	if err != nil {
		return "", nil, "", "", "", err
	}

	return response.URL, response.Cookie, response.UserID, response.AuthToken, response.SessionToken, nil
}

// Makes a http get request to the BigBlueButton API to join a meeting and returs:
// - url
func (api *ApiRequest) JoinGetURL(meetingID string, userName string, moderator bool) (string, error) {

	meetings, err := api.GetMeetings()
	if err != nil {
		return "", err
	}

	m := meetings[meetingID]

	var password string
	if moderator {
		password = m.ModeratorPW
	} else {
		password = m.AttendeePW
	}

	params := []params{
		{
			name:  FULL_NAME,
			value: userName,
		},
		{
			name:  MEETING_ID,
			value: meetingID,
		},
		{
			name:  PASSWORD,
			value: password,
		},
		{
			name:  REDIRECT,
			value: "true",
		},
	}

	url := api.buildURL(JOIN, params...)

	return url, nil
}
