package kafka

import (
	"encoding/json"
)

type Message struct {
	OriginalBody string `json:"OriginalBody"`
	MessageTime  string `json:"messageTime"`
	Type         string `json:"type"`
	Body         string `json:"Body"`
	Uuid         string `json:"uuid"`
}

func InitMessage(message string, originMessageBody string, messageType string, MessageId string, MessageTime string) *Message {
	return &Message{
		originMessageBody,
		MessageTime,
		messageType,
		message,
		MessageId,
	}
}

func (message *Message) ToJson() string {
	marshal, err := json.Marshal(message)
	if err != nil {
		return ""
	}
	return string(marshal)
}
