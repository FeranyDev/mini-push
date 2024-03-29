package push

import (
	"encoding/json"
	"fmt"
	"github.com/feranydev/mini-push/database"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"github.com/feranydev/mini-push/usergroup"
)

var (
	TgBot *tgbotapi.BotAPI
)

func Send(user usergroup.User, title, text string, copy bool) (string, error) {

	parse, err := uuid.Parse(user.PushId)
	if err != nil {
		return "", err
	}

	messageId := uuid.New()

	go database.SaveMessage(messageId, parse, title, text)
	go send(user, title, text, copy)

	return messageId.String(), nil
}

func send(user usergroup.User, title, text string, copy bool) {

	switch user.PushServer {
	case "tg":
		{
			msg := tgbotapi.NewMessage(user.TgId, "")
			if copy {
				if title != "" {
					msg.Text = fmt.Sprintf("*%s*\n`%s`", title, text)
				} else {
					msg.Text = fmt.Sprintf("`%s`", text)
				}
			} else {
				if title != "" {
					msg.Text = fmt.Sprintf("*%s*\n%s", title, text)
				} else {
					msg.Text = fmt.Sprintf("%s", text)
				}
			}
			msg.ParseMode = "markdown"
			_, err := TgBot.Send(msg)
			if err != nil {
				log.Errorf("telegram bot send message error: %s", err)
			}
			return
		}

	// https://github.com/xlvecle/PushLite
	case "PushLite":
		tmp := "0"
		if copy {
			tmp = "1"
		}
		fcmData := struct {
			To          string `json:"to"`
			CollapseKey string `json:"collapse_key"`
			Data        struct {
				Body     string `json:"body"`
				Title    string `json:"title"`
				AutoCopy string `json:"autoCopy"`
				MsgType  string `json:"msgType"`
			} `json:"text"`
		}{
			To:          user.ServerToken,
			CollapseKey: "type_a",
			Data: struct {
				Body     string `json:"body"`
				Title    string `json:"title"`
				AutoCopy string `json:"autoCopy"`
				MsgType  string `json:"msgType"`
			}{
				Body:     text,
				Title:    title,
				AutoCopy: tmp,
				MsgType:  "normal",
			},
		}

		marshal, _ := json.Marshal(fcmData)
		req, _ := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", strings.NewReader(string(marshal)))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("%s=%s", "key", "AIzaSyAd-JC3NxVeGRHyo5ZZB2BUmhSA7Z_IqHY"))
		if resp, err := (&http.Client{}).Do(req); err != nil {
			log.Errorf(err.Error())
		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			log.Debugf("fcm response: %s", body)
		}
	default:
		log.Errorf("Unknown push server: %s", user.PushServer)
	}
}
