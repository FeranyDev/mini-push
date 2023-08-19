package main

import (
	"encoding/json"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"github.com/feranydev/mini-push/api"
	"github.com/feranydev/mini-push/config"
	"github.com/feranydev/mini-push/database"
	"github.com/feranydev/mini-push/push"
	"github.com/feranydev/mini-push/usergroup"
)

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

func init() {

	if config.Deploy.Debug {
		log.SetLevel(log.DEBUG)
		log.SetHeader("${time_rfc3339} ${level} ${short_file}:${line} ${message}")
	} else {
		log.SetHeader("${time_rfc3339} ${level} ${message}")
	}
}

func main() {

	deploy := config.Deploy

	usergroup.DefaultUserGroup.ReadUsers()

	database.Connections()

	bot, err := tgbotapi.NewBotAPI(deploy.TgBot)
	if err != nil {
		log.Panic(err)
	}
	if deploy.Debug {
		bot.Debug = true
		log.Infof("%s", deploy)
		log.Infof("%s", usergroup.DefaultUserGroup.Data)
	}
	push.TgBot = bot

	log.Infof("Authorized on account %s", bot.Self.UserName)
	log.Infof("Service started successfully currently one with %d users", len(usergroup.DefaultUserGroup.Data))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//go api.Start(bot, deploy)
	go api.Start(bot)

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			go func(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "create":
						log.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)
						user, presence := usergroup.DefaultUserGroup.FindUserByTgId(update.Message.Chat.ID)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
						if presence {
							msg.Text = "你已经注册过了，推送ID为: " + user.PushId
						} else {
							pushID := uuid.NewString()
							usergroup.DefaultUserGroup.AddUser(usergroup.User{PushId: pushID, TgId: update.Message.Chat.ID, PushServer: "tg"})
							msg.Text = "注册成功，推送ID为: " + pushID
						}
						_, _ = bot.Send(msg)
					case "bindpushlite":
						if update.Message.CommandArguments() != "" {
							token := uuid.NewString()
							usergroup.DefaultUserGroup.AddUser(usergroup.User{
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
						users := usergroup.DefaultUserGroup.FindAllPush(update.Message.Chat.ID)
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
					case "delete":
						if update.Message.CommandArguments() != "" {
							if err := usergroup.DefaultUserGroup.DeletePushToken(update.Message.CommandArguments(), update.Message.Chat.ID); err != nil {
								log.Errorf("删除推送ID失败: %s", err)
								_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "删除失败，可能是推送ID不存在或者你没有权限删除"))
							} else {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "删除成功")
								_, _ = bot.Send(msg)
							}
						} else {
							_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "请使用 /delete [推送ID] 来删除推送设备"))
						}
					default:
						_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "未知指令"))
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
