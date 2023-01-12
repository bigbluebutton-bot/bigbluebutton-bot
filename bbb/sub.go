package bbb

type SubType int

const (
	UsersSub SubType = iota
	PollsSub
	PresentationsSub
	SlidesSub
	SlidePositionsSub
	CaptionsSub
	VoiceUsersSub
	WhiteboardMultiUserSub
	ScreenshareSub
	GroupChatSub
	PresentationPodsSub
	UsersSettingsSub
	GuestUserSub
	UsersInfosSub
	MeetingTimeRemainingSub
	LocalSettingsSub
	UsersTypingSub
	RecordMeetingsSub
	VideoStreamsSub
	ConnectionStatusSub
	VoiceCallStatesSub
	ExternalVideoMeetingsSub
	BreakoutsSub
	BreakoutsHistorySub
	PadsSub
	PadsSessionsSub
	PadsUpdatesSub
	UsersPersistentDataSub
	StreamCursorSub
	StreamAnnotationsAddedSub
	StreamAnnotationsRemovedSub
	GroupChatMsgSub
	CurrentPollSub
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
	case UsersSub:
		return "users", []interface{}{}
	case PollsSub:
		return "polls", []interface{}{}
	case PresentationsSub:
		return "presentations", []interface{}{}
	case SlidesSub:
		return "slides", []interface{}{}
	case SlidePositionsSub:
		return "slide-positions", []interface{}{}
	case CaptionsSub:
		return "captions", []interface{}{}
	case VoiceUsersSub:
		return "voice-users", []interface{}{}
	case WhiteboardMultiUserSub:
		return "whiteboard-multi-user", []interface{}{}
	case ScreenshareSub:
		return "screenshare", []interface{}{}
	case GroupChatSub:
		return "group-chat", []interface{}{}
	case PresentationPodsSub:
		return "presentation-pods", []interface{}{}
	case UsersSettingsSub:
		return "users-settings", []interface{}{}
	case GuestUserSub:
		return "guest-user", []interface{}{}
	case UsersInfosSub:
		return "users-infos", []interface{}{}
	case MeetingTimeRemainingSub:
		return "meeting-time-remaining", []interface{}{}
	case LocalSettingsSub:
		return "local-settings", []interface{}{}
	case UsersTypingSub:
		return "users-typing", []interface{}{}
	case RecordMeetingsSub:
		return "record-meetings", []interface{}{}
	case VideoStreamsSub:
		return "video-streams", []interface{}{}
	case ConnectionStatusSub:
		return "connection-status", []interface{}{}
	case VoiceCallStatesSub:
		return "voice-call-states", []interface{}{}
	case ExternalVideoMeetingsSub:
		return "external-video-meetings", []interface{}{}
	case BreakoutsSub:
		return "breakouts", []interface{}{}
	case BreakoutsHistorySub:
		return "breakouts-history", []interface{}{}
	case PadsSub:
		return "pads", []interface{}{}
	case PadsSessionsSub:
		return "pads-sessions", []interface{}{}
	case PadsUpdatesSub:
		return "pads-updates", []interface{}{}
	case UsersPersistentDataSub:
		return "users-persistent-data", []interface{}{}
	case StreamCursorSub:
		return "stream-cursor-INTERNALID", []interface{}{st} //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	case StreamAnnotationsAddedSub:
		return "stream-annotations-INTERNALID", []interface{}{st} //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	case StreamAnnotationsRemovedSub:
		return "stream-annotations-INTERNALID", []interface{}{st} //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	case GroupChatMsgSub:
		return "group-chat-msg", []interface{}{0}
	case CurrentPollSub:
		return "current-poll", []interface{}{false, true}
	default:
		return "", []interface{}{}
	}
}
