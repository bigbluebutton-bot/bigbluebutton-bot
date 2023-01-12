package bbb

import (
	ddp "ddp"
	convert "github.com/benpate/convert"
)

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

// Converts a map[string]interface{} (from ddp.Update) into a Message object
func ConvertInToMessage(content ddp.Update) Message {
	return Message{
		ID:                 convert.String(content["id"]),
		Timestamp:          convert.Int64(content["timestamp"]),
		CorrelationID:      convert.String(content["correlationId"]),
		ChatEmphasizedText: convert.Bool(content["chatEmphasizedText"]),
		Message:            convert.String(content["message"]),
		Sender:             convert.String(content["sender"]),
		SenderName:         convert.String(content["senderName"]),
		SenderRole:         convert.String(content["senderRole"]),
		MeetingId:          convert.String(content["meetingId"]),
		ChatId:             convert.String(content["chatId"]),
	}
}
