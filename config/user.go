package config

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	TgBot *tgbotapi.BotAPI
)

type T struct {
	lck   sync.Mutex
	Users []User `json:"user"`
}

type User struct {
	PushId      string `json:"pushId"`
	TgId        int64  `json:"tgId"`
	PushServer  string `json:"pushServer"`
	ServerToken string `json:"serverToken,omitempty"`
}

func (t *T) ReadData() {
	t.lck.Lock()
	defer t.lck.Unlock()
	file, err := os.ReadFile("./users.json")
	if err != nil {
		log.Panicf("Error reading file: %v", err)
	}

	if err = json.Unmarshal(file, &t); err != nil {
		log.Panicf("Error while unmarshaling config: %s", err)
	}
	//log.Debug("Read users.json file: %s", t.Users)
}

func (t *T) AddUser(user User) {
	t.lck.Lock()
	defer t.lck.Unlock()
	t.Users = append(t.Users, user)
	tmp, err := json.Marshal(t)
	if err != nil {
		log.Debug("Error while marshaling config: %s", err)
	}
	if err = os.WriteFile("./users.json", tmp, 0644); err != nil {
		log.Errorf("Error while writing config: %s", err)
	}
}

func (t *T) FindAllPush(tgId int64) (users []User) {
	for _, user := range t.Users {
		if user.TgId == tgId {
			users = append(users, user)
		}
	}
	return
}

func (t *T) FindUserByTgId(tgId int64) (User, bool) {
	for _, user := range t.Users {
		if user.TgId == tgId && user.PushServer == "tg" {
			return user, true
		}
	}
	return User{}, false
}

func (t *T) FindUserByPushId(pushId string) (User, bool) {
	for _, user := range t.Users {
		if user.PushId == pushId {
			return user, true
		}
	}
	return User{}, false
}

func (user *User) Send(title, data string, copy bool) (err error) {
	switch user.PushServer {
	case "tg":
		{
			msg := tgbotapi.NewMessage(user.TgId, "")
			if copy {
				if title != "" {
					msg.Text = fmt.Sprintf("*%s*\n`%s`", title, data)
				} else {
					msg.Text = fmt.Sprintf("`%s`", data)
				}
			} else {
				if title != "" {
					msg.Text = fmt.Sprintf("*%s*\n%s", title, data)
				} else {
					msg.Text = fmt.Sprintf("%s", data)
				}
			}
			msg.ParseMode = "markdown"
			_, err = TgBot.Send(msg)
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
			} `json:"data"`
		}{
			To:          user.ServerToken,
			CollapseKey: "type_a",
			Data: struct {
				Body     string `json:"body"`
				Title    string `json:"title"`
				AutoCopy string `json:"autoCopy"`
				MsgType  string `json:"msgType"`
			}{
				Body:     data,
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
			return err
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			log.Debugf("fcm response: %s", body)
			return nil
		}
	}
	return fmt.Errorf("unknown push server")
}
