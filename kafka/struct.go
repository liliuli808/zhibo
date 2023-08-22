package kafka

import (
	"encoding/json"
)

type Message struct {
	OriginalBody string `json:"OriginalBody"`
	MessageTime  string `json:"messageTime"`
	NickName     string `json:"nickName"`
	Type         string `json:"type"`
	Body         string `json:"Body"`
	Uuid         string `json:"uuid"`
}

func InitMessage(nickName string, message string, originMessageBody string, messageType string, MessageId string, MessageTime string) *Message {
	return &Message{
		originMessageBody,
		MessageTime,
		nickName,
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
