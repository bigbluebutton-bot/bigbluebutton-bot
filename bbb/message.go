package bbb

type Message struct {
	ID                 string `json:"id"`
	Timestamp          int64  `json:"timestamp"`
	CorrelationID      string `json:"correlationId"`
	ChatEmphasizedText bool   `json:"chatEmphasizedText"`
	Message            string `json:"message"`
	Sender             string `json:"sender"`
	SenderName         string `json:"senderName"`
	SenderRole         string `json:"senderRole"`
	MeetingId          string `json:"meetingId"`
	ChatId             string `json:"chatId"`
}
