package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"zhibo/bot"
	"zhibo/kafka"
	"zhibo/levelDb"
	"zhibo/mysql"
	"zhibo/utils"
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
	Mysql  *mysql.Mysql
	Config *Config
	Db     *levelDb.LevelDb
	Bot    *bot.Bot
}

type plantBotResponse struct {
	Succeeded bool `json:"succeeded"`
	RespData  struct {
		Topics []struct {
			TopicId int64 `json:"topic_id"`
			Group   struct {
				GroupId int64  `json:"group_id"`
				Name    string `json:"name"`
				Type    string `json:"type"`
			} `json:"group"`
			Type string `json:"type"`
			Talk struct {
				Owner struct {
					UserId    int64  `json:"user_id"`
					Name      string `json:"name"`
					AvatarUrl string `json:"avatar_url"`
				} `json:"owner"`
				Text   string `json:"text"`
				Images []struct {
					ImageId   int64  `json:"image_id"`
					Type      string `json:"type"`
					Thumbnail struct {
						Url    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"thumbnail"`
					Large struct {
						Url    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"large"`
					Original struct {
						Url    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
						Size   int    `json:"size"`
					} `json:"original,omitempty"`
				} `json:"images,omitempty"`
			} `json:"talk"`
			LikesCount    int    `json:"likes_count"`
			RewardsCount  int    `json:"rewards_count"`
			CommentsCount int    `json:"comments_count"`
			ReadingCount  int    `json:"reading_count"`
			ReadersCount  int    `json:"readers_count"`
			Digested      bool   `json:"digested"`
			Sticky        bool   `json:"sticky"`
			CreateTime    string `json:"create_time"`
			UserSpecific  struct {
				Liked      bool `json:"liked"`
				Subscribed bool `json:"subscribed"`
			} `json:"user_specific"`
			LatestLikes []struct {
				CreateTime string `json:"create_time"`
				Owner      struct {
					UserId    int64  `json:"user_id"`
					Name      string `json:"name"`
					AvatarUrl string `json:"avatar_url"`
				} `json:"owner"`
			} `json:"latest_likes,omitempty"`
			ShowComments []struct {
				CommentId  int64  `json:"comment_id"`
				CreateTime string `json:"create_time"`
				Owner      struct {
					UserId    int64  `json:"user_id"`
					Name      string `json:"name"`
					AvatarUrl string `json:"avatar_url"`
				} `json:"owner"`
				Text            string `json:"text"`
				LikesCount      int    `json:"likes_count"`
				RewardsCount    int    `json:"rewards_count"`
				Sticky          bool   `json:"sticky"`
				RepliesCount    int    `json:"replies_count,omitempty"`
				ParentCommentId int64  `json:"parent_comment_id,omitempty"`
				Repliee         struct {
					UserId    int64  `json:"user_id"`
					Name      string `json:"name"`
					AvatarUrl string `json:"avatar_url"`
				} `json:"repliee,omitempty"`
			} `json:"show_comments,omitempty"`
		} `json:"topics"`
	} `json:"resp_data"`
}

func NewAgent(config *Config) *Agent {
	agent := &Agent{}
	agent.Config = config
	agent.Mysql = &mysql.Mysql{Config: agent.Config.MysqlConfig}
	agent.Mysql.Init()
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
	if err != nil {
		log.Println(err)
	}
	agent.parasJson(body)
}

func (agent *Agent) parasJson(s []byte) {
	db := levelDb.NewLevelDbInstance("./data/leveldb")
	defer db.Close()
	var resp Response
	err := json.Unmarshal(s, &resp)
	if err != nil {
		log.Println("error:", err)
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
			originMessageBody = message.NickName + ":" + originMessage.Body
		}
		one, err := db.HasOne(message.MessageId)

		if one == true {
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		err, err2 := db.Put(message.MessageId, "true")
		if err != nil || err2 != nil {
			log.Println(err, err2)
		}
		_, err = mysql.StructInsert(agent.Mysql.MysqlDb, trimHtml(message.Body), filterEmoji(trimHtml(originMessageBody)), messageType, message.MessageId, message.MessageTime)
		if err != nil {
			log.Println(err)
		}
		message := kafka.InitMessage(message.NickName, trimHtml(message.Body), filterEmoji(trimHtml(originMessageBody)), messageType, message.MessageId, message.MessageTime)
		var resMess string
		if err != nil {
			panic(err)
		}
		if message.Type == "answer" {
			resMess = "问：" + message.OriginalBody + "\n答：" + message.NickName + message.Body
		} else {
			resMess = message.NickName + ":" + message.Body
		}
		str, imageArr := getImagePath(resMess)

		if str != "" {
			agent.sendTextMessage(str)
		}

		if len(imageArr) != 0 {
			agent.sendImageMessage(imageArr)
		}
	}
}

func (agent *Agent) StartSendPlantBot() {
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
	req.Header.Set("authority", "api.zsxq.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("origin", "https://wx.zsxq.com")
	req.Header.Set("referer", "https://wx.zsxq.com/")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	agent.parasJsonPlantBot(body)
}

func (agent *Agent) parasJsonPlantBot(s []byte) {
	db := levelDb.NewLevelDbInstance("./data/leveldb")
	defer db.Close()
	var resp plantBotResponse
	err := json.Unmarshal(s, &resp)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, message := range resp.RespData.Topics {

		hasTalk, err := db.HasOne("talkId" + strconv.FormatInt(message.TopicId, 10))

		if err != nil {
			fmt.Println(err)
		}
		var commentArr []string

		if message.Talk.Text == "" {
			continue
		}

		ti, _ := time.Parse("2006-01-02T15:04:05.000+0800", message.CreateTime)
		commentArr = append(commentArr, ti.Format("2006-01-02 15:04"))
		commentArr = append(commentArr, message.Talk.Text)
		var hasComment bool

		for _, comment := range message.ShowComments {
			if comment.Text == "" {
				continue
			}
			var resArr []string
			resArr = append(resArr, ti.Format("2006-01-02 15:04"))
			resArr = append(resArr, message.Talk.Text)
			resArr = append(resArr, "评论: "+comment.Text)

			hasComment, err = db.HasOne("commentId" + strconv.FormatInt(comment.CommentId, 10))
			if err != nil {
				fmt.Println(err)
			}
			if hasComment {
				continue
			}

			db.Put("commentId"+strconv.FormatInt(comment.CommentId, 10), "1")
			s := utils.TextToImage(resArr)
			abs, err := filepath.Abs(s)
			if err != nil {
				return
			}
			marshal, err := json.Marshal(ImageBody{Id: "703653853",
				Message: ImageMessage{Type: "image", Data: ImageData{File: "file://" + abs}}})
			if err != nil {
				return
			}

			rsp, err := http.Post("http://127.0.0.1:5700/send_group_msg", "application/json", bytes.NewReader(marshal))
			if err != nil {
				panic(err)
			}
			rsp.Body.Close()
			ioutil.ReadAll(rsp.Body)
		}

		if hasTalk {
			continue
		}

		db.Put("talkId"+strconv.FormatInt(message.TopicId, 10), "1")
		s := utils.TextToImage(commentArr)
		abs, err := filepath.Abs(s)
		if err != nil {
			return
		}
		marshal, err := json.Marshal(ImageBody{Id: "703653853",
			Message: ImageMessage{Type: "image", Data: ImageData{File: "file://" + abs}}})
		if err != nil {
			return
		}

		rsp, err := http.Post("http://127.0.0.1:5700/send_group_msg", "application/json", bytes.NewReader(marshal))
		if err != nil {
			panic(err)
		}
		rsp.Body.Close()
		ioutil.ReadAll(rsp.Body)

		if len(message.Talk.Images) > 0 {
			for _, image := range message.Talk.Images {
				marshal, err := json.Marshal(ImageBody{Id: "703653853",
					Message: ImageMessage{Type: "image", Data: ImageData{File: image.Large.Url}}})
				rsp, err := http.Post("http://127.0.0.1:5700/send_group_msg", "application/json", bytes.NewReader(marshal))
				if err != nil {
					panic(err)
				}
				rsp.Body.Close()
			}
		}

	}
}

type ImageMessage struct {
	Type string    `json:"type"`
	Data ImageData `json:"data"`
}

type ImageBody struct {
	Id      string       `json:"group_id"`
	Message ImageMessage `json:"message"`
}

type ImageData struct {
	File string `json:"file"`
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

func (agent *Agent) sendTextMessage(str string) {
	s := utils.TextToImage(strings.Split(str, "\n"))
	defer os.Remove(s)
	filePath, _ := filepath.Abs(s)
	marshal, err := json.Marshal(ImageBody{Id: agent.Config.QqGroupId,
		Message: ImageMessage{Type: "image", Data: ImageData{File: "file://" + filePath}}})
	if err != nil {
		return
	}
	rsp, err := http.Post(agent.Config.Api, "application/json", bytes.NewReader(marshal))
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	_, err = ioutil.ReadAll(rsp.Body)
}

func (agent *Agent) sendImageMessage(arr []string) {
	for _, s := range arr {
		marshal, err := json.Marshal(ImageBody{Id: agent.Config.QqGroupId, Message: ImageMessage{Type: "image", Data: ImageData{File: s}}})
		if err != nil {
			return
		}
		rsp, err := http.Post(agent.Config.Api, "application/json", bytes.NewReader(marshal))
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
