package main

import (
	"encoding/json"
	"flag"
	"github.com/feranydev/mini-push/api"
	"github.com/feranydev/mini-push/config"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	flag.Parse()
	args := flag.Args()

	deploy := config.Deploy

	t.ReadData()

	bot, err := tgbotapi.NewBotAPI(deploy.TgBot)
	if err != nil {
		log.Panic(err)
	}

	if len(args) != 0 {
		if args[0] == "debug" {
			bot.Debug = true
			log.SetLevel(log.DEBUG)
		}
	}

	config.TgBot = bot
	log.Infof("Authorized on account %s", bot.Self.UserName)
	log.Infof("Service started successfully currently one with %d users", len(t.Users))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go api.Start(bot, &t, deploy)

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
							t.AddUser(config.User{PushId: pushID, TgId: update.Message.Chat.ID, PushServer: "tg"})
							msg.Text = "注册成功，推送ID为: " + pushID
						}
						_, _ = bot.Send(msg)
					case "bindpushlite":
						if update.Message.CommandArguments() != "" {
							token := uuid.NewString()
							t.AddUser(config.User{
								PushId:      token,
								TgId:        update.Message.Chat.ID,
								PushServer:  "PushLite",
								ServerToken: update.Message.CommandArguments(),
							})
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "绑定成功，推送ID为: "+token)
							_, _ = bot.Send(msg)
						} else {
							_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "请使用 /bindPushLite [客户端Token] 来绑定推送设备"))
						}
					case "select":
						users := t.FindAllPush(update.Message.Chat.ID)
						marshal, _ := json.Marshal(users)
						data := string(marshal)
						data = strings.Replace(data, "pushId", "推送ID", -1)
						data = strings.Replace(data, "pushServer", "推送服务", -1)
						data = strings.Replace(data, "serverToken", "客户端Token", -1)
						data = strings.Replace(data, "tgId", "TelegramID", -1)
						data = strings.Replace(data, ",", ",\n", -1)
						data = strings.Replace(data, "},\n{", "},\n\n{", -1)
						data = strings.Replace(data, "[", "", -1)
						data = strings.Replace(data, "]", "", -1)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, data)
						_, _ = bot.Send(msg)
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
					_, _ = bot.Send(msg)
				}
			}(bot, &update)
		}
	}
}
