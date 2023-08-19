package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

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
	if config.Deploy.Redis.Addr == "" {
		log.Infof("redis not set")
		return
	} else {
		redisStart()
	}
	if config.Deploy.LocalMemory {
		lmStart()
	}
}

func SaveMessage(messageId, pushID uuid.UUID, title, text string) bool {
	data := sqlMessage{
		PushID:    pushID,
		Title:     title,
		Text:      text,
		SpeedTime: time.Now().Unix(),
	}

	success := true

	if config.Deploy.Redis.Addr != "" {
		err := redisSet(messageId, data)
		if err != nil {
			success = false
		}
	}

	if config.Deploy.LocalMemory {
		ok := lmSet(messageId, data)
		if !ok {
			success = false
		}
	}

	if success {
		return true
	}

	return false
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

	if config.Deploy.LocalMemory {
		data, ok := lmGet(msgId)
		if ok {
			return Message(data), true
		}
	}

	return Message{}, false
}
