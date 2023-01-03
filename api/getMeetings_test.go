package api

import (
	"testing"
)

// Thes tests only succeed if there is a BBB server running, which does not have any other meetings running!!!

// Test for GetMeetings
func TestGetMeetings(t *testing.T) {
	conf := readConfig("../config.json", t)	// Read config from file or environment. This function is in "api_test.go"

	bbbapi, err := NewRequest(conf.Url, conf.Secret, SHA1)
	if err != nil {
		t.Errorf("GetMeetings() %d FAILED: NewRequest: %s", 0, err)
		return
	}

	newmeeting, err := bbbapi.CreateMeeting("name", "meetingID", "attendeePW", "moderatorPW", "welcome text", false, false, false, 12345)
	if err != nil {
		t.Errorf("GetMeetings() %d FAILED: CreateMeeting: %s", 0, err)
		return
	}

	meetings, err := bbbapi.GetMeetings()
	if err != nil {
		t.Errorf("GetMeetings() %d FAILED: GetMeetings: %s", 0, err)
		return
	}

	value, ok := meetings[newmeeting.MeetingID]
	if(!ok) {
		t.Errorf("GetMeetings() %d FAILED: Meeting was not found in GetMeetings: %s", 0, value.MeetingName)
		t.Log(value)
		return
	}
	t.Logf("GetMeetings() %d PASSED %s", 0, value.MeetingName)
}