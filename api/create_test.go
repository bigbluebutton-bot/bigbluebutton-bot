package api

import "testing"

// Thes tests only succeed if there is a BBB server running, which does not have any other meetings running!!!

type testcreatemeeting struct {
	allow_start_stop_recording bool
	attendee_pw                string
	auto_start_recording       bool
	meeting_id                 string
	moderator_pw               string
	name                       string
	record                     bool
	voice_bridge               int64
	welcome                    string

	shouldfail bool
}

// Test for CreateMeeting
func TestCreateMeeting(t *testing.T) {
	tests := []testcreatemeeting{
		{ //0
			allow_start_stop_recording: true,
			attendee_pw: "ap",
			auto_start_recording: false,
			meeting_id: "random-4026116",
			moderator_pw: "mp",
			name: "random-4026116",
			record: false,
			voice_bridge: 70848,
			welcome: "Hello you there",
			shouldfail: false,
		},
		{ //1 same voice bridge as 0
			allow_start_stop_recording: true,
			attendee_pw: "ap",
			auto_start_recording: false,
			meeting_id: "random-4026115",
			moderator_pw: "mp",
			name: "random-4026115",
			record: false,
			voice_bridge: 70848,
			welcome: "Hello you there",
			shouldfail: true,
		},
		{ //2 same meeting id as 0
			allow_start_stop_recording: true,
			attendee_pw: "ap",
			auto_start_recording: false,
			meeting_id: "random-4026116",
			moderator_pw: "mp",
			name: "random-dfdfdf",
			record: false,
			voice_bridge: 12345,
			welcome: "Hello you there",
			shouldfail: false,
		},
		{ //3
			allow_start_stop_recording: true,
			attendee_pw: "apvvvvvvvvvvvvvvvvjhkhkjhkhjkhjkhjkhjkllllllllllllllllllkooiiaahf", // max sizte 64 chars
			auto_start_recording: false,
			meeting_id: "random-54321",
			moderator_pw: "mp",
			name: "random-54321",
			record: false,
			voice_bridge: 54321,
			welcome: "Hello you there",
			shouldfail: true,
		},
		{ //4
			allow_start_stop_recording: true,
			attendee_pw: "sds", // max sizte 64 chars
			auto_start_recording: false,
			meeting_id: "random-5432112",
			moderator_pw: "apvvvvvvvvvvvvvvvvjhkhkjhkhjkhjkhjkhjkllllllllllllllllllkooiiaahf",
			name: "random-5432112",
			record: false,
			voice_bridge: 5432112,
			welcome: "Hello you there",
			shouldfail: true,
		},
	}

	conf := readConfig("../_example/config.json", t)	// Read config from file or environment. This function is in "api_test.go"

	bbbapi, err := NewRequest(conf.BBB.API.URL, conf.BBB.API.Secret, SHA1)
	if err != nil {
		t.Errorf("CreateMeeting() FAILED: NewRequest: %s", err)
		return
	}


	for num, test := range tests {
		newmeeting, err := bbbapi.CreateMeeting(test.name, test.meeting_id, test.attendee_pw, test.moderator_pw, test.welcome, test.allow_start_stop_recording, test.auto_start_recording, test.record, test.voice_bridge)
		if err != nil {
			if test.shouldfail {
				t.Logf("CreateMeeting() %d PASSED: CreateMeeting: %s", num, err)
				continue
			}
			t.Errorf("CreateMeeting() %d FAILED: CreateMeeting: %s", num, err)
			continue
		}
	
		meetings, err := bbbapi.GetMeetings()
		if err != nil {
			t.Errorf("CreateMeeting() %d FAILED: GetMeetings: %s", num, err)
			continue
		}
	
		value, ok := meetings[newmeeting.MeetingID]
		if(!ok) {
			t.Errorf("CreateMeeting() %d FAILED: Meeting was not found in GetMeetings: %s", num, value.MeetingName)
			t.Log(value)
			continue
		}

		if test.shouldfail {
			t.Errorf("CreateMeeting() %d FAILED: There was no error. But one was expected! %s", num, value.MeetingName)
			continue
		}

		t.Logf("CreateMeeting() %d PASSED %s", 0, value.MeetingName)
	}

}