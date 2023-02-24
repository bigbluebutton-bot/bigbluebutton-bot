package bbb

type CallType int

const (
	StreamCursor CallType = iota //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	VoidConnection
	StopUserTyping
	SendGroupChatMsg
	StartUserTyping
	ChatMessageBeforeJoinCounter
	UserChangedLocalSettings
	ValidateAuthToken
	GetPadId
	CreateSession
	UpdateCaptionsOwner
	FetchMessagePerPage
	UserShareWebcam
	ZoomSlide
	SwitchSlide
	SetEmojiStatus
	ToggleVoice
	UserUnshareWebcam
	SendBulkAnnotations
	ClearWhiteboard
	UndoAnnotation
	AddGlobalAccess
	MuteAllExceptPresenter
	MuteAllUsers
	ToggleWebcamsOnlyForModerator
	ToggleLockSettings
	ChangeGuestPolicy
	EndAllBreakouts
	CreateBreakoutRoom
	CreateGroup
	RemovePresentation
	UserLeftMeeting
	SetExitReason
)

func GetCall(callName CallType) (string) {
	switch callName {
	case StreamCursor:
		return "stream-cursor-INTERNALID" //BECAREFUL: The word INTERNALID must be replaced with the internal meeting id!!!
	case VoidConnection:
		return "voidConnection"
	case StopUserTyping:
		return "stopUserTyping"
	case SendGroupChatMsg:
		return "sendGroupChatMsg"
	case StartUserTyping:
		return "startUserTyping"
	case ChatMessageBeforeJoinCounter:
		return "chatMessageBeforeJoinCounter"
	case UserChangedLocalSettings:
		return "userChangedLocalSettings"
	case ValidateAuthToken:
		return "validateAuthToken"
	case GetPadId:
		return "getPadId"
	case CreateSession:
		return "createSession"
	case UpdateCaptionsOwner:
		return "updateCaptionsOwner"
	case FetchMessagePerPage:
		return "fetchMessagePerPage"
	case UserShareWebcam:
		return "userShareWebcam"
	case ZoomSlide:
		return "zoomSlide"
	case SwitchSlide:
		return "switchSlide"
	case SetEmojiStatus:
		return "setEmojiStatus"
	case ToggleVoice:
		return "toggleVoice"
	case UserUnshareWebcam:
		return "userUnshareWebcam"
	case SendBulkAnnotations:
		return "sendBulkAnnotations"
	case ClearWhiteboard:
		return "clearWhiteboard"
	case UndoAnnotation:
		return "undoAnnotation"
	case AddGlobalAccess:
		return "addGlobalAccess"
	case MuteAllExceptPresenter:
		return "muteAllExceptPresenter"
	case MuteAllUsers:
		return "muteAllUsers"
	case ToggleWebcamsOnlyForModerator:
		return "toggleWebcamsOnlyForModerator"
	case ToggleLockSettings:
		return "toggleLockSettings"
	case ChangeGuestPolicy:
		return "changeGuestPolicy"
	case EndAllBreakouts:
		return "endAllBreakouts"
	case CreateBreakoutRoom:
		return "createBreakoutRoom"
	case CreateGroup:
		return "createGroup"
	case RemovePresentation:
		return "removePresentation"
	case UserLeftMeeting:
		return "userLeftMeeting"
	case SetExitReason:
		return "setExitReason"
	default:
		return ""
	}
}