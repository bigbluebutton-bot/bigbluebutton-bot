package api

type responseEndMeeting struct {
	Script      string `xml:"script"`
	ReturnCode  string `xml:"returncode"`
	Errors   		[]responseerror `xml:"errors>error"`
	MessageKey  string `xml:"messageKey"`
	Message     string `xml:"message"`
}

// Makes a http get request to the BigBlueButton API and returns the closed meeting
func (api *api_request) EndMeeting(meetingID string) (meeting, error) {

	meetings, err:= api.GetMeetings()
	if(err != nil){
		return meeting{}, err
	}
	m := meetings[meetingID]

	params := []params{
		{
			name:  MEETING_ID,
			value: meetingID,
		},
		{
			name:  PASSWORD,
			value: m.ModeratorPW,
		},
	}

	var response responseEndMeeting
	err = api.makeRequest(&response, END, params...)
	if(err != nil){
		return meeting{}, err
	}

	return m, nil
}