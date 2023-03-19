package usergroup

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"sync"
)

var DefaultUserGroup T

type User struct {
	PushId      string `json:"pushId"`
	TgId        int64  `json:"tgId"`
	PushServer  string `json:"pushServer"`
	ServerToken string `json:"serverToken,omitempty"`
}

type T struct {
	lck  sync.Mutex
	Data []User `json:"user"`
}

func (t *T) ReadUsers() {
	t.lck.Lock()
	defer t.lck.Unlock()
	file, err := os.ReadFile("./users.json")
	if err != nil {
		log.Panicf("Error reading file: %v", err)
	}

	if err = json.Unmarshal(file, &t); err != nil {
		log.Panicf("Error while unmarshaling config: %s", err)
	}
	log.Infof("Read users.json success")
}

func (t *T) AddUser(user User) {
	t.lck.Lock()
	defer t.lck.Unlock()
	t.Data = append(t.Data, user)
	tmp, err := json.Marshal(t)
	if err != nil {
		log.Debug("Error while marshaling config: %s", err)
	}
	if err = os.WriteFile("./users.json", tmp, 0644); err != nil {
		log.Errorf("Error while writing config: %s", err)
	}
}

func (t *T) DeletePushToken(pushId string, tgId int64) error {
	t.lck.Lock()
	defer t.lck.Unlock()
	for i, u := range t.Data {
		if u.PushId == pushId && u.TgId == tgId {
			t.Data = append(t.Data[:i], t.Data[i+1:]...)
			tmp, err := json.Marshal(t)
			if err != nil {
				log.Debug("Error while marshaling config: %s", err)
			}
			if err = os.WriteFile("./users.json", tmp, 0644); err != nil {
				log.Errorf("Error while writing config: %s", err)
			}
			return nil
		}
	}
	return fmt.Errorf("not found user")
}

func (t *T) FindAllPush(tgId int64) (users []User) {
	for _, user := range t.Data {
		if user.TgId == tgId {
			users = append(users, user)
		}
	}
	return
}

func (t *T) FindUserByTgId(tgId int64) (User, bool) {
	for _, user := range t.Data {
		if user.TgId == tgId && user.PushServer == "tg" {
			return user, true
		}
	}
	return User{}, false
}

func (t *T) FindUserByPushId(pushId string) (User, bool) {
	for _, user := range t.Data {
		if user.PushId == pushId {
			return user, true
		}
	}
	return User{}, false
}
