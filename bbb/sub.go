package bbb

type SubType int

const (
	Users SubType = iota
	Polls
	Presentations
	Slides
	SlidePositions
	Captions
	VoiceUsers
	WhiteboardMultiUser
	Screenshare
	GroupChat
	PresentationPods
	UsersSettings
	GuestUser
	UsersInfos
	MeetingTimeRemaining
	LocalSettings
	UsersTyping
	RecordMeetings
	VideoStreams
	ConnectionStatus
	VoiceCallStates
	ExternalVideoMeetings
	Breakouts
	BreakoutsHistory
	Pads
	PadsSessions
	PadsUpdates
	UsersPersistentData
	StreamCursor
	StreamAnnotationsAdded
	StreamAnnotationsRemoved
	GroupChatMsg
	CurrentPoll
)

type streamsettings struct {
	UseCollection bool     `json:"useCollection"`
	Args          []string `json:"args"`
}

// Returns the sub and all parameters
// IMPORTENT: If the name conatins the word `INTERNALID`, it MUST be replaced with the internal meeting id!!!
func GetSub(SubName SubType) (string, []interface{}) {
	st := streamsettings{
		UseCollection: false,
		Args:          []string{},
	}

	switch SubName {
	case Users:
		return "users", []interface{}{}
	case Polls:
		return "polls", []interface{}{}
	case Presentations:
		return "presentations", []interface{}{}
	case Slides:
		return "slides", []interface{}{}
	case SlidePositions:
		return "slide-positions", []interface{}{}
	case Captions:
		return "captions", []interface{}{}
	case VoiceUsers:
		return "voice-users", []interface{}{}
	case WhiteboardMultiUser:
		return "whiteboard-multi-user", []interface{}{}
	case Screenshare:
		return "screenshare", []interface{}{}
	case GroupChat:
		return "group-chat", []interface{}{}
	case PresentationPods:
		return "presentation-pods", []interface{}{}
	case UsersSettings:
		return "users-settings", []interface{}{}
	case GuestUser:
		return "guest-user", []interface{}{}
	case UsersInfos:
		return "users-infos", []interface{}{}
	case MeetingTimeRemaining:
		return "meeting-time-remaining", []interface{}{}
	case LocalSettings:
		return "local-settings", []interface{}{}
	case UsersTyping:
		return "users-typing", []interface{}{}
	case RecordMeetings:
		return "record-meetings", []interface{}{}
	case VideoStreams:
		return "video-streams", []interface{}{}
	case ConnectionStatus:
		return "connection-status", []interface{}{}
	case VoiceCallStates:
		return "voice-call-states", []interface{}{}
	case ExternalVideoMeetings:
		return "external-video-meetings", []interface{}{}
	case Breakouts:
		return "breakouts", []interface{}{}
	case BreakoutsHistory:
		return "breakouts-history", []interface{}{}
	case Pads:
		return "pads", []interface{}{}
	case PadsSessions:
		return "pads-sessions", []interface{}{}
	case PadsUpdates:
		return "pads-updates", []interface{}{}
	case UsersPersistentData:
		return "users-persistent-data", []interface{}{}
	case StreamCursor:
		return "stream-cursor-INTERNALID", []interface{}{st} //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	case StreamAnnotationsAdded:
		return "stream-annotations-INTERNALID", []interface{}{st} //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	case StreamAnnotationsRemoved:
		return "stream-annotations-INTERNALID", []interface{}{st} //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	case GroupChatMsg:
		return "group-chat-msg", []interface{}{0}
	case CurrentPoll:
		return "current-poll", []interface{}{false, true}
	default:
		return "", []interface{}{}
	}
}

