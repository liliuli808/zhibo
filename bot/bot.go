package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"zhibo/kafka"
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

	if resMess != "" {
		body := Body{Id: b.C.QqGroupId, Message: str}
		marshal, err := json.Marshal(body)
		if err != nil {
			return
		}
		rsp, err := http.Post(b.C.Api, "application/json", bytes.NewReader(marshal))
		if err != nil {
			log.Fatal(err)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(rsp.Body)
	}

	if len(imageArr) != 0 {
		b.sendImageMessage(imageArr)
	}
}

func (b *Bot) sendImageMessage(arr []string) {
	for _, s := range arr {
		marshal, err := json.Marshal(ImageBody{Id: b.C.QqGroupId, Message: ImageMessage{Type: "image", Data: ImageData{File: s}}})
		if err != nil {
			return
		}
		fmt.Println(string(marshal))
		rsp, err := http.Post(b.C.Api, "application/json", bytes.NewReader(marshal))
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(rsp.Body)
		fmt.Println(string(body))
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
