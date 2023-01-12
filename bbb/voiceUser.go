package bbb

// For voice-users
type VoiceUser struct {
	IntId       string `json:"intId"`
	MeetingId   string `json:"meetingId"`
	CallerName  string `json:"callerName"`
	CallerNum   string `json:"callerNum"`
	CallingWith string `json:"callingWith"`
	Joined      bool   `json:"joined"`
	ListenOnly  bool   `json:"listenOnly"`
	Muted       bool   `json:"muted"`
	Spoke       bool   `json:"spoke"`
	Talking     bool   `json:"talking"`
	VoiceConf   string `json:"voiceConf"`
	VoiceUserId string `json:"voiceUserId"`
}