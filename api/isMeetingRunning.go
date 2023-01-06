package api

type responseIsMeetingRunning struct {
	Script      string `xml:"script"`
	ReturnCode  string `xml:"returncode"`
	Errors   		[]responseerror `xml:"errors>error"`
	Running	 	bool   `xml:"running"`
	MessageKey  string `xml:"messageKey"`
	Message     string `xml:"message"`
}

// Makes a http get request to the BigBlueButton API and returs the running state of the meeting.
// If an error occurs the returned value is false
func (api *api_request) IsMeetingRunning(meetingID string) bool {

	params := []params{
		{
			name:  MEETING_ID,
			value: meetingID,
		},
	}

	var response responseIsMeetingRunning
	err := api.makeRequest(&response, IS_MEETING_RUNNING, params...)
	if(err != nil){
		return false
	}

	return response.Running
}