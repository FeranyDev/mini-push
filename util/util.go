package util

import (
	"os"

	"github.com/labstack/gommon/log"
)

func DefaultString(value, defaultValue string) string {
	if value != "" || value != *new(string) {
		return value
	} else {
		return defaultValue
	}
}

func DefaultInt(value, defaultValue int) int {
	if value != *new(int) {
		return value
	} else {
		return defaultValue
	}
}

func CheckFile(path string) (err error) {
	exist, err := pathExists(path)
	if err != nil {
		return
	}
	if exist {
		log.Infof("has dir![%v]", path)
	} else {
		log.Infof("no dir![%v]", path)
		err = os.WriteFile(path, []byte(""), 0644)
		if err != nil {
			return
		} else {
			log.Infof("touch success!")
		}
	}
	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
