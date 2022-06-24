package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
	"zhibo/kafka"
	"zhibo/levelDb"
	"zhibo/mysql"
)

type OriginMessage struct {
	Body string `json:"body"`
}

type Response struct {
	ResultKey string `json:"resultKey"`
	Data      struct {
		Messages []Message `json:"messages"`
	} `json:"data"`
}

type Message struct {
	MessageTime       string `json:"messageTime"`
	PrimaryTeacher    bool   `json:"primaryTeacher"`
	FromRoomId        int    `json:"fromRoomId"`
	IsTeachersStudent bool   `json:"isTeachersStudent"`
	NickName          string `json:"nickName"`
	MessageId         string `json:"messageId"`
	NewStudent        bool   `json:"newStudent"`
	VerifyTime        int64  `json:"verifyTime"`
	Type              string `json:"type"`
	Body              string `json:"body"`
	UserId            int    `json:"userId"`
	Uuid              string `json:"uuid,omitempty"`
	UserImage         string `json:"userImage"`
	MultFlag          string `json:"multFlag"`
	IsMedal           bool   `json:"isMedal"`
	Topic             string `json:"topic"`
	From              string `json:"from"`
	Attributes        string `json:"attributes"`
	IsCrown           bool   `json:"isCrown"`
	ContentType       string `json:"contentType"`
	IsComment         int    `json:"isComment"`
	OriginalMessageId int64  `json:"originalMessageId,omitempty"`
	OriginalMessage   string `json:"originalMessage,omitempty"`
}

type Agent struct {
	Mysql   *mysql.Mysql
	Config  *Config
	Product *kafka.Product
	Wg      *sync.WaitGroup
	Db      *levelDb.LevelDb
}

func NewAgent(config *Config) *Agent {
	agent := &Agent{}
	agent.Mysql = &mysql.Mysql{Config: config.MysqlConfig}
	agent.Mysql.Init()
	agent.Config = config
	agent.Product = &kafka.Product{Config: kafka.Config{Address: config.KafkaConfig.Address}}
	agent.Product.Instance()
	agent.Wg = &sync.WaitGroup{}
	return agent
}

func (agent *Agent) Start() {
	data, err := ioutil.ReadFile(agent.Config.ApiConfig.CookiePath)
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		agent.Config.ApiConfig.ApiAddress,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Cookie", strings.TrimSpace(string(data)))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	agent.parasJson(body)
}

func (agent *Agent) parasJson(s []byte) {
	db := levelDb.NewLevelDbInstance("./data/leveldb")
	defer db.Close()
	var resp Response
	err := json.Unmarshal(s, &resp)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, message := range resp.Data.Messages {
		messageType := "anno"
		originMessageBody := ""
		if message.OriginalMessage != "" {
			var originMessage OriginMessage
			err := json.Unmarshal([]byte(message.OriginalMessage), &originMessage)
			if err != nil {
				log.Fatal(err)
			}
			messageType = "answer"
			originMessageBody = originMessage.Body
		}
		one, err := db.HasOne(message.MessageId)

		if one == true {
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		err, err2 := db.Put(message.MessageId, "true")
		if err != nil {
			fmt.Println(err, err2)
		}
		_, err = mysql.StructInsert(agent.Mysql.MysqlDb, trimHtml(message.Body), filterEmoji(trimHtml(originMessageBody)), messageType, message.MessageId, message.MessageTime)
		if err != nil {
			continue
		}
		err = agent.Product.Push(
			agent.Config.KafkaConfig.Topic,
			kafka.InitMessage(trimHtml(message.Body), filterEmoji(trimHtml(originMessageBody)), messageType, message.MessageId, message.MessageTime).ToJson(),
		)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func trimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	//去除表情
	re, _ = regexp.Compile("/ud([8-9a-f][0-9a-z]{2})/i")
	src = re.ReplaceAllString(src, "")
	return strings.TrimSpace(src)
}

func filterEmoji(content string) string {
	newContent := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			newContent += string(value)
		}
	}
	return newContent
}
