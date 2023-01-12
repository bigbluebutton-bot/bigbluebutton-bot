# bigbluebutton-bot

```plantuml
@startuml
package bbb #DDDDDD {
    class Connection {
        + CurrentPoll:                      string
        ' + StreamAnnotations:                string'
        + UsersPersistentData:              string
        + GroupChatMsg:                     string
        + StreamCursor:                     string
        + PadsUpdates:                      string
        + PadsSessions:                     string
        + Pads:                             string
        + BreakoutsHistory:                 string
        + Breakouts:                        string
        + ExternalVideoMeetings:            string
        + VoiceCallStates:                  string
        + ConnectionStatus:                 string
        + VideoStreams:                     string
        + RecordMeetings:                   string
        + UserTyping:                       string
        + LocalSettings:                    string
        + MeetingTimeRemaining:             string
        + UsersInfos:                       string
        + GuestUser:                        string
        + UsersSettings:                    string
        + PresentationPods:                 string
        + GroupChat:                        string
        + Screenshare:                      string
        + WhiteboardMultiUser:              string
        + VoiceUsers:                       string
        + Captions:                         string
        + SlidePositions:                   string
        + Slides:                           string
        + Presentations:                    string
        + Polls:                            string
        + Meetings:                         string
        + Users:                            string
        + CurrentUser:                      string
        ' + MeteorAutoUpdateClientVersions:   string'
    }

    ' class Recording {
    '     + MeetingID:     string
    ' }

    class Meeting {
        + MeetingID:     string
        + MeetingName:   string
        + InternalID:    string
        + CreateTime:    int64
        + CreateDate:    string
        + VoiceBridge:   int
        + DialNumber:    string
        + AttendeePW:    string
        + ModeratorPW:   string
        + Running:       bool
        + Duration:      int
        + HasJoined:     bool
        + Recording:     bool
        + ForciblyEnded: bool
        + StartTime:     int64
        + EndTime:       int64
        + Participants:  int
        + Listeners:     int
        + VoiceCount:    int
        + VideoCount:    int
        + MaxUsers:      int
        + Moderators:    int
        + Attendees:     []attendee
        '+  Metadata:      metadata'
        + IsBreakout:    bool

        + Create():   error
        + End():      error
        + IsRunning(): bool
        + Join(Username string, isModerator bool):     User, error

    }
    
    class User {
        + UserID:         string
        + FullName:       string
        + Role:           string
        + IsPresenter:    bool
        + IsListening:    bool
        + HasJoinedVoice: bool
        + HasVideo:       bool
        + ClientType:     string
    }

    class Message {
        id:                 string
        timestamp:          int
        correlationId:      string
        chatEmphasizedText: bool
        message:            string
        sender:             string
        senderName:         string
        senderRole:         RoleType
        meetingId:          string
        chatId:             string
    }
}

package bot #DDDDDD {
    class Client {
        + Status: 	        StatusType
        + ClientURL:			string
        + ClientWSURL:		string
        + ApiURL:				string
        - apiSecret:			string
        + API: 				*api.ApiRequest
        - ddpClient: 			*ddp.Client
        - event: 				*event
    }
}

package api #DDDDDD {
    class ApiRequest {
        + Url:        string
        + Secret:     string
        + Sha:        ShaType
    }
}

@enduml
```

All subscriptions:
```go
	err = c.ddpClient.Sub("users")
	if err != nil {
		panic("")
	}
	if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("polls")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("presentations")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("slides")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("slide-positions")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("captions")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("voiceUsers")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("whiteboard-multi-user")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("screenshare")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("group-chat")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("presentation-pods")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("users-settings")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("guestUser")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("users-infos")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("meeting-time-remaining")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("local-settings")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("users-typing")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("record-meetings")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("video-streams")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("connection-status")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("voice-call-states")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("external-video-meetings")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("breakouts")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("breakouts-history")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("pads")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("pads-sessions")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("pads-updates")
    if err != nil {
        panic("")
    }

	type streamsettings struct {
		UseCollection bool		`json:"useCollection"`
		Args          []string	`json:"args"`
	}
	st := streamsettings{
		UseCollection: 	false,
		Args:		 	[]string{},
	}
	err = c.ddpClient.Sub("stream-cursor-" + internalMeetingID, "message", st)
    if err != nil {
        panic("stream-cursor-" + internalMeetingID)
    }
	err = c.ddpClient.Sub("stream-annotations-" + internalMeetingID, "removed", st)
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("stream-annotations-" + internalMeetingID, "added", st)
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("group-chat-msg", 0)
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("users-persistent-data")
    if err != nil {
        panic("")
    }
	err = c.ddpClient.Sub("current-poll", false, true)
    if err != nil {
        panic("")
    }
```