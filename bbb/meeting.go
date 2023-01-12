package bbb

type Meeting struct {
	MeetingId            string                   `json:"meetingId"`
	BreakoutProps        MeetingBreakoutProps     `json:"breakoutProps"`
	DurationProps        MeetingDurationProps     `json:"durationProps"`
	Groups               []interface{}            `json:"groups"`
	GuestLobbyMessage    string                   `json:"guestLobbyMessage"`
	Layout               string                   `json:"layout"`
	LockSettingsProps    MeetingLockSettingsProps `json:"lockSettingsProps"`
	MeetingEnded         bool                     `json:"meetingEnded"`
	MeetingProp          MeetingMeetingProp       `json:"meetingProp"`
	MetadataProp         MeetingMetadataProp      `json:"metadataProp"`
	PublishedPoll        bool                     `json:"publishedPoll"`
	RandomlySelectedUser []interface{}            `json:"randomlySelectedUser"`
	SystemProps          MeetingSystemProps       `json:"systemProps"`
	UsersProp            MeetingUsersProp         `json:"usersProp"`
	VoiceProp            MeetingVoiceProp         `json:"voiceProp"`
	WelcomeProp          MeetingWelcomeProp       `json:"welcomeProp"`
}
type MeetingBreakoutProps struct {
	BreakoutRooms      []interface{} `json:"breakoutRooms"`
	FreeJoin           bool          `json:"freeJoin"`
	ParentId           string        `json:"parentId"`
	PrivateChatEnabled bool          `json:"privateChatEnabled"`
	Record             bool          `json:"record"`
	Sequence           int           `json:"sequence"`
}
type MeetingDurationProps struct {
	CreatedDate                            string `json:"createdDate"`
	CreatedTime                            int    `json:"createdTime"`
	Duration                               int    `json:"duration"`
	EndWhenNoModerator                     bool   `json:"endWhenNoModerator"`
	EndWhenNoModeratorDelayInMinutes       int    `json:"endWhenNoModeratorDelayInMinutes"`
	MeetingExpireIfNoUserJoinedInMinutes   int    `json:"meetingExpireIfNoUserJoinedInMinutes"`
	MeetingExpireWhenLastUserLeftInMinutes int    `json:"meetingExpireWhenLastUserLeftInMinutes"`
	TimeRemaining                          int    `json:"timeRemaining"`
	UserActivitySignResponseDelayInMinutes int    `json:"userActivitySignResponseDelayInMinutes"`
	UserInactivityInspectTimerInMinutes    int    `json:"userInactivityInspectTimerInMinutes"`
	UserInactivityThresholdInMinutes       int    `json:"userInactivityThresholdInMinutes"`
}
type MeetingLockSettingsProps struct {
	DisableCam             bool   `json:"disableCam"`
	DisableMic             bool   `json:"disableMic"`
	DisableNotes           bool   `json:"disableNotes"`
	DisablePrivateChat     bool   `json:"disablePrivateChat"`
	DisablePublicChat      bool   `json:"disablePublicChat"`
	HideUserList           bool   `json:"hideUserList"`
	HideViewersCursor      bool   `json:"hideViewersCursor"`
	LockOnJoin             bool   `json:"lockOnJoin"`
	LockOnJoinConfigurable bool   `json:"lockOnJoinConfigurable"`
	LockedLayout           bool   `json:"lockedLayout"`
	SetBy                  string `json:"setBy"`
}
type MeetingMeetingProp struct {
	DisabledFeatures []string `json:"disabledFeatures"`
	ExtId            string   `json:"extId"`
	IntId            string   `json:"intId"`
	IsBreakout       bool     `json:"isBreakout"`
	MeetingCameraCap int      `json:"meetingCameraCap"`
	Name             string   `json:"name"`
}
type MeetingMetadataProp struct {
	Metadata map[string]interface{} `json:"metadata"`
}
type MeetingSystemProps struct {
	Html5InstanceId int `json:"html5InstanceId"`
}
type MeetingUsersProp struct {
	AllowModsToEjectCameras bool   `json:"allowModsToEjectCameras"`
	AllowModsToUnmuteUsers  bool   `json:"allowModsToUnmuteUsers"`
	AuthenticatedGuest      bool   `json:"authenticatedGuest"`
	GuestPolicy             string `json:"guestPolicy"`
	MaxUsers                int    `json:"maxUsers"`
	MeetingLayout           string `json:"meetingLayout"`
	UserCameraCap           int    `json:"userCameraCap"`
	WebcamsOnlyForModerator bool   `json:"webcamsOnlyForModerator"`
}
type MeetingVoiceProp struct {
	DialNumber  string `json:"dialNumber"`
	MuteOnStart bool   `json:"muteOnStart"`
	TelVoice    string `json:"telVoice"`
	VoiceConf   string `json:"voiceConf"`
}
type MeetingWelcomeProp struct {
	WelcomeMsg         string `json:"welcomeMsg"`
	WelcomeMsgTemplate string `json:"welcomeMsgTemplate"`
}
