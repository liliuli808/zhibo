package bot

import (
	"bytes"
	"encoding/json"
	"github.com/Shopify/sarama"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"zhibo/kafka"
	"zhibo/utils"
)

type Bot struct {
	C        *Config
	Consumer *kafka.ConsumerGroup
}

func Instance(c *Config) *Bot {
	bot := Bot{C: c}
	bot.Consumer = &kafka.ConsumerGroup{
		Config:  kafka.Config{Address: c.KafkaAddress},
		GroupID: c.GroupId,
		Topics:  []string{c.Topic},
		GroupHandler: kafka.ConsumerGroupHandler{
			Handler: bot.SendMessage,
		},
	}
	return &bot
}

func (b *Bot) SendMessage(m *sarama.ConsumerMessage, count int) {
	var resMess string
	var message kafka.Message
	err := json.Unmarshal(m.Value, &message)
	if err != nil {
		panic(err)
	}

	if message.Type == "answer" {
		resMess = "问：" + message.OriginalBody + "\n答：" + message.Body
	} else {
		resMess = message.Body
	}
	str, imageArr := getImagePath(resMess)

	if str != "" {
		b.sendTextMessage(str)
	}

	if len(imageArr) != 0 {
		b.sendImageMessage(imageArr)
	}
}

func (b *Bot) sendTextMessage(str string) {
	s := utils.TextToImage(strings.Split(str, "\n"))
	defer os.Remove(s)
	filePath, _ := filepath.Abs(s)
	marshal, err := json.Marshal(ImageBody{Id: b.C.QqGroupId,
		Message: ImageMessage{Type: "image", Data: ImageData{File: "file://" + filePath}}})
	if err != nil {
		return
	}
	rsp, err := http.Post(b.C.Api, "application/json", bytes.NewReader(marshal))
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	_, err = ioutil.ReadAll(rsp.Body)
}

func (b *Bot) sendImageMessage(arr []string) {
	for _, s := range arr {
		marshal, err := json.Marshal(ImageBody{Id: b.C.QqGroupId, Message: ImageMessage{Type: "image", Data: ImageData{File: s}}})
		if err != nil {
			return
		}
		rsp, err := http.Post(b.C.Api, "application/json", bytes.NewReader(marshal))
		if err != nil {
			panic(err)
		}
		_, err = ioutil.ReadAll(rsp.Body)
		rsp.Body.Close()
	}
}

func getImagePath(messageStr string) (string, []string) {
	var start []int
	var end []int
	var res []string
	for i, _ := range messageStr {
		if i+5 > len(messageStr) {
			continue
		}
		if i+6 > len(messageStr) {
			continue
		}

		if messageStr[i:i+5] == "[img]" {
			start = append(start, i+5)
		}
		if messageStr[i:i+6] == "[/img]" {
			end = append(end, i)
		}
	}
	resStr := messageStr
	for i, v := range start {
		res = append(res, messageStr[v:end[i]])
		if i == 0 {
			resStr = messageStr[:v-5]
		} else {
			resStr += messageStr[end[i-1]+6 : v-5]
		}
	}

	return resStr, res
}
