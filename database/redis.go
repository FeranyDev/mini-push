package database

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/feranydev/mini-push/config"
	"github.com/feranydev/mini-push/util"
)

var ctx = context.Background()

var rdb = &redis.Client{}

func redisStart() {

	deploy := config.Deploy

	rdb = redis.NewClient(&redis.Options{
		Addr:     util.DefaultString(deploy.Redis.Addr, "localhost:6379"),
		Username: util.DefaultString(deploy.Redis.User, ""),
		Password: util.DefaultString(deploy.Redis.Pass, ""),
		DB:       util.DefaultInt(deploy.Redis.Db, 0),
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	res, err := rdb.Get(ctx, "key").Result()
	if err != nil || res != "value" {
		panic(err)
	}

	err = rdb.Del(ctx, "key").Err()
	if err != nil {
		panic(err)
	}

	log.Infof("redis check success")
}

func redisSet(data sqlMessage) (messageID uuid.UUID, err error) {
	messageID = uuid.New()
	marshal, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = rdb.Set(ctx, messageID.String(), base64.StdEncoding.EncodeToString(marshal), 7*24*time.Hour).Err()
	return
}

func redisGet(msgId uuid.UUID) (data sqlMessage, err error) {
	result64, err := rdb.Get(ctx, msgId.String()).Result()
	if err != nil {
		return data, err
	}
	result, err := base64.StdEncoding.DecodeString(result64)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(result, &data)
	if err != nil {
		return sqlMessage{}, err
	}
	return
}
