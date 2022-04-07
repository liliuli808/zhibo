package mysql

import (
	"database/sql"
	"fmt"
	"time"
)

type Message struct {
	Id           int64  `db:"id"`
	Uuid         string `db:"uuid"`
	Body         string `db:"body"`
	OriginalBody string `db:"original_body"`
	Type         string `db:"type"`
}

func StructInsert(mysqlDb *sql.DB, body string, originalBody string, typeStr string, uuid string, messageTime string) {
	parse, _ := time.Parse("20060102150405", messageTime)
	res, err := mysqlDb.Exec(
		"insert INTO message(body,original_body,type,uuid,created_at) values(?,?, ? ,?,?)",
		body,
		originalBody,
		typeStr,
		uuid,
		parse.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.LastInsertId())
	}
}

func GetNotSendMessage(mysqlDb *sql.DB) *Message {
	rows, _ := mysqlDb.Query("SELECT id,uuid,body,original_body,type FROM `message` where is_send = 0 limit 1")
	var message Message
	if rows == nil {
		return nil
	}
	for rows.Next() {
		rows.Scan(&message.Id, &message.Uuid, &message.Body, &message.OriginalBody, &message.Type)
	}
	fmt.Println(message.Uuid)
	return &message
}

func UpdateSendState(mysqlDb *sql.DB, messageId string) {
	exec, err := mysqlDb.Exec("UPDATE zhibo.message t SET t.is_send = 1 WHERE t.id = ?", messageId)
	if err != nil {
		panic(err)
	}
	fmt.Println(exec.LastInsertId())
}
