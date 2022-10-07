package main

import (
	"encoding/json"
	"github.com/feranydev/push-server/api"
	"github.com/feranydev/push-server/config"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	tgToken = "5137003298:AAFKiOeK6z-EaujbZ2u4cj_iBOYcb77ZUlE"
)

var t config.T

type V1Hitokoto struct {
	Id         int    `json:"id"`
	Uuid       string `json:"uuid"`
	Hitokoto   string `json:"hitokoto"`
	Type       string `json:"type"`
	From       string `json:"from"`
	FromWho    string `json:"from_who"`
	Creator    string `json:"creator"`
	CreatorUid int    `json:"creator_uid"`
	Reviewer   int    `json:"reviewer"`
	CommitFrom string `json:"commit_from"`
	CreatedAt  string `json:"created_at"`
	Length     int    `json:"length"`
}

func main() {

	t.ReadData()

	log.SetLevel(log.DEBUG)

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Infof("Authorized on account %s", bot.Self.UserName)
	log.Infof("Service started successfully currently one with %d users", len(t.User))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go api.Start(bot, &t)

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			go func(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "create":
						log.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)
						user, presence := t.FindUserByTgId(update.Message.Chat.ID)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
						if presence {
							msg.Text = "你已经注册过了，推送ID为: " + user.PushId
						} else {
							pushID := uuid.NewString()
							t.AddUser(config.User{PushId: pushID, TgId: update.Message.Chat.ID})
							msg.Text = "注册成功，推送ID为: " + pushID
						}
						bot.Send(msg)
					}
				} else {
					log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
					resp, err := http.Get("https://international.v1.hitokoto.cn")
					if err != nil {
						log.Errorf("Error while sending hitokoto: %s", err)
					}
					buff := make([]byte, 10240)
					read, err := resp.Body.Read(buff)
					if err != nil {
						log.Errorf("Error while reading response body: %s", err)
					}
					v1Hitokoto := V1Hitokoto{}
					if err = json.Unmarshal(buff[:read], &v1Hitokoto); err != nil {
						log.Errorf("Error while unmarshaling hitokoto: %s", err)
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, v1Hitokoto.Hitokoto)
					bot.Send(msg)
				}
			}(bot, &update)
		}
	}
}
