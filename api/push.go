package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/feranydev/mini-push/push"
	"github.com/feranydev/mini-push/usergroup"
)

func text(c echo.Context) error {
	pushId := c.Param("pushId")
	text := c.Param("text")
	if user, ok := usergroup.DefaultUserGroup.FindUserByPushId(pushId); ok {
		if msgId, err := push.Send(user, "", text, false); err != nil {
			log.Errorf("send message to %d error: %s", user.PushId, err)
			return c.JSON(http.StatusInternalServerError, back{
				Code:      500,
				Msg:       "send message error",
				Timestamp: time.Now().Unix(),
			})
		} else {
			return success(c, msgId)
		}
	} else {
		return userNotFound(c)
	}
}

func textAndTitle(c echo.Context) error {
	pushId := c.Param("pushId")
	title := c.Param("title")
	text := c.Param("text")
	if user, ok := usergroup.DefaultUserGroup.FindUserByPushId(pushId); ok {
		if msgId, err := push.Send(user, title, text, false); err != nil {
			log.Errorf("send message to %d error: %s", user.PushId, err)
			return c.JSON(http.StatusInternalServerError, back{
				Code:      500,
				Msg:       "send message error",
				Timestamp: time.Now().Unix(),
			})
		} else {
			return success(c, msgId)
		}
	} else {
		return userNotFound(c)
	}
}

func textCopy(c echo.Context) error {
	pushId := c.Param("pushId")
	text := c.Param("text")
	if user, ok := usergroup.DefaultUserGroup.FindUserByPushId(pushId); ok {
		if msgId, err := push.Send(user, "", text, true); err != nil {
			log.Errorf("send message to %d error: %s", user.PushId, err)
			return c.JSON(http.StatusInternalServerError, back{
				Code:      500,
				Msg:       "send message error",
				Timestamp: time.Now().Unix(),
			})
		} else {
			return success(c, msgId)
		}
	} else {
		return userNotFound(c)
	}
}

func textAndTitleCopy(c echo.Context) error {
	pushId := c.Param("pushId")
	title := c.Param("title")
	text := c.Param("text")
	if user, ok := usergroup.DefaultUserGroup.FindUserByPushId(pushId); ok {
		if msgId, err := push.Send(user, title, text, true); err != nil {
			log.Errorf("send message to %d error: %s", user.PushId, err)
			return c.JSON(http.StatusInternalServerError, back{
				Code:      500,
				Msg:       "send message error",
				Timestamp: time.Now().Unix(),
			})
		} else {
			return success(c, msgId)
		}
	} else {
		return userNotFound(c)
	}
}

func post(c echo.Context) error {
	if user, ok := usergroup.DefaultUserGroup.FindUserByPushId(c.Param("pushId")); ok {
		data := struct {
			Title   string `json:"title"`
			Text    string `json:"text"`
			Message string `json:"message"`
			Copy    bool   `json:"copy"`
		}{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, back{
				Code:      400,
				Msg:       "bad request",
				Timestamp: time.Now().Unix(),
			})
		}
		log.Debugf("push data: %+v", data)
		tmp := data.Text
		if tmp == "" {
			tmp = data.Message
		}
		if msgId, err := push.Send(user, data.Title, tmp, data.Copy); err != nil {
			log.Errorf("send message to %d error: %s", user.PushId, err)
			return c.JSON(http.StatusInternalServerError, back{
				Code:      500,
				Msg:       "send message error",
				Timestamp: time.Now().Unix(),
			})
		} else {
			return success(c, msgId)
		}
	} else {
		return userNotFound(c)
	}
}

func pushDeer(c echo.Context) error {
	data := struct {
		PushKey string `json:"pushkey"`
		Text    string `json:"text"`
		Desp    string `json:"desp"`
		Type    string `json:"type"`
	}{}
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, back{
			Code:      400,
			Msg:       "bad request",
			Timestamp: time.Now().Unix(),
		})
	}
	log.Debugf("push data: %+v", data)
	if user, ok := usergroup.DefaultUserGroup.FindUserByPushId(data.PushKey); ok {
		if msgId, err := push.Send(user, data.Text, data.Desp, true); err != nil {
			log.Errorf("send message to %d error: %s", user.TgId, err)
			return sendError(c)
		} else {
			return success(c, msgId)
		}
	} else {
		return userNotFound(c)
	}
}
