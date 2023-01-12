package bbb

// For GoupChat and Private Chats
type Chat struct {
	ChatId       string             `json:"chatId"`
	MeetingId    string             `json:"meetingId"`
	Access       string             `json:"access"`
	CreatedBy    string             `json:"createdBy"`
	Participants []ChatParticipants `json:"participants"`
	Users        []string           `json:"users"`
}
type ChatParticipants struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
