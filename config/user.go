package config

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"os"
	"sync"
)

type T struct {
	lck  sync.Mutex
	User []User `json:"user"`
}

type User struct {
	PushId string `json:"pushId"`
	TgId   int64  `json:"tgId"`
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
	log.Debug("Read users.json file: %s", t.User)
}

func (t *T) AddUser(user User) {
	t.lck.Lock()
	defer t.lck.Unlock()
	t.User = append(t.User, user)
	tmp, err := json.Marshal(t)
	if err != nil {
		log.Debug("Error while marshaling config: %s", err)
	}
	if err = os.WriteFile("./users.json", tmp, 0644); err != nil {
		log.Errorf("Error while writing config: %s", err)
	}
}

func (t *T) FindUserByTgId(tgId int64) (User, bool) {
	for _, user := range t.User {
		if user.TgId == tgId {
			return user, true
		}
	}
	return User{}, false
}

func (t *T) FindUserByPushId(pushId string) (User, bool) {
	for _, user := range t.User {
		if user.PushId == pushId {
			return user, true
		}
	}
	return User{}, false
}
