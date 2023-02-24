package bbb

type CallType int

const (
	StreamCursorCall CallType = iota //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	VoidConnectionCall
	StopUserTypingCall
	SendGroupChatMsgCall
	StartUserTypingCall
	ChatMessageBeforeJoinCounterCall
	UserChangedLocalSettingsCall
	ValidateAuthTokenCall
	GetPadIdCall
	CreateSessionCall
	UpdateCaptionsOwnerCall
	FetchMessagePerPageCall
	UserShareWebcamCall
	ZoomSlideCall
	SwitchSlideCall
	SetEmojiStatusCall
	ToggleVoiceCall
	UserUnshareWebcamCall
	SendBulkAnnotationsCall
	ClearWhiteboardCall
	UndoAnnotationCall
	AddGlobalAccessCall
	MuteAllExceptPresenterCall
	MuteAllUsersCall
	ToggleWebcamsOnlyForModeratorCall
	ToggleLockSettingsCall
	ChangeGuestPolicyCall
	EndAllBreakoutsCall
	CreateBreakoutRoomCall
	CreateGroupCall
	RemovePresentationCall
	UserLeftMeetingCall
	SetExitReasonCall
)

func GetCall(callName CallType) (string) {
	switch callName {
	case StreamCursorCall:
		return "stream-cursor-INTERNALID" //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	case VoidConnectionCall:
		return "voidConnection"
	case StopUserTypingCall:
		return "stopUserTyping"
	case SendGroupChatMsgCall:
		return "sendGroupChatMsg"
	case StartUserTypingCall:
		return "startUserTyping"
	case ChatMessageBeforeJoinCounterCall:
		return "chatMessageBeforeJoinCounter"
	case UserChangedLocalSettingsCall:
		return "userChangedLocalSettings"
	case ValidateAuthTokenCall:
		return "validateAuthToken"
	case GetPadIdCall:
		return "getPadId"
	case CreateSessionCall:
		return "createSession"
	case UpdateCaptionsOwnerCall:
		return "updateCaptionsOwner"
	case FetchMessagePerPageCall:
		return "fetchMessagePerPage"
	case UserShareWebcamCall:
		return "userShareWebcam"
	case ZoomSlideCall:
		return "zoomSlide"
	case SwitchSlideCall:
		return "switchSlide"
	case SetEmojiStatusCall:
		return "setEmojiStatus"
	case ToggleVoiceCall:
		return "toggleVoice"
	case UserUnshareWebcamCall:
		return "userUnshareWebcam"
	case SendBulkAnnotationsCall:
		return "sendBulkAnnotations"
	case ClearWhiteboardCall:
		return "clearWhiteboard"
	case UndoAnnotationCall:
		return "undoAnnotation"
	case AddGlobalAccessCall:
		return "addGlobalAccess"
	case MuteAllExceptPresenterCall:
		return "muteAllExceptPresenter"
	case MuteAllUsersCall:
		return "muteAllUsers"
	case ToggleWebcamsOnlyForModeratorCall:
		return "toggleWebcamsOnlyForModerator"
	case ToggleLockSettingsCall:
		return "toggleLockSettings"
	case ChangeGuestPolicyCall:
		return "changeGuestPolicy"
	case EndAllBreakoutsCall:
		return "endAllBreakouts"
	case CreateBreakoutRoomCall:
		return "createBreakoutRoom"
	case CreateGroupCall:
		return "createGroup"
	case RemovePresentationCall:
		return "removePresentation"
	case UserLeftMeetingCall:
		return "userLeftMeeting"
	case SetExitReasonCall:
		return "setExitReason"
	default:
		return ""
	}
}