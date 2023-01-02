package api

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

// Thes tests only succeed if there is a BBB server running, which does not have any other meetings running!!!

type config struct {
	Url    string `json:"url"`
	Secret string `json:"secret"`
}
func readConfig(file string) config {
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we initialize config
	var conf config
	// we unmarshal our byteArray which contains our jsonFile's content into conf
	json.Unmarshal([]byte(byteValue), &conf)

	return conf
}

// Test for GetMeetings
func TestGetMeetings(t *testing.T) {
	conf := readConfig("../config_test.json")

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