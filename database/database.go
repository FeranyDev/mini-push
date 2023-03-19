package database

import (
	"github.com/labstack/gommon/log"
	"time"

	"github.com/google/uuid"

	"github.com/feranydev/mini-push/config"
)

type sqlMessage struct {
	PushID    uuid.UUID `json:"pushID"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	SpeedTime int64     `json:"speedTime"`
}

type Message sqlMessage

func Connections() {
	redisStart()
}

func SaveMessage(pushID uuid.UUID, title, text string) (string, bool) {
	if config.Deploy.Redis.Addr != "" {
		messageID, err := redisSet(sqlMessage{
			PushID:    pushID,
			Title:     title,
			Text:      text,
			SpeedTime: time.Now().Unix(),
		})
		if err != nil {
			return "error", false
		}
		return messageID.String(), true
	}

	return "404", false

}

func GetMessageByMsgID(msgId uuid.UUID) (Message, bool) {
	if config.Deploy.Redis.Addr != "" {
		data, err := redisGet(msgId)
		if err != nil {
			log.Errorf("redis get message error: %v", err)
			return Message{}, false
		}
		return Message(data), true
	}

	return Message{}, false
}
