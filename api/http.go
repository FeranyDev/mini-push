package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/feranydev/mini-push/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type back struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
}

func Start(bot *tgbotapi.BotAPI, users *config.T, deploy *config.Config) {
	e := echo.New()

	e.HideBanner = true

	e.IPExtractor = echo.ExtractIPFromRealIPHeader()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP: true,
		LogMethod:   true,
		LogURI:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Infof("%v %v %v", v.RemoteIP, v.Method, v.URI)
			return nil
		},
	}))

	//e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//	Format: "${time_rfc3339} ${remote_ip} ${method} ${uri}\n",
	//}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Oh! It's working!")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, back{
			Code:      200,
			Msg:       "pong",
			Timestamp: time.Now().Unix(),
		})
	})
	e.GET("/push/:pushId/:text", func(c echo.Context) error {
		pushId := c.Param("pushId")
		text := c.Param("text")
		if user, ok := users.FindUserByPushId(pushId); ok {
			if err := user.Send("", text, false); err != nil {
				log.Errorf("send message to %d error: %s", user.PushId, err)
				return c.JSON(http.StatusInternalServerError, back{
					Code:      500,
					Msg:       "send message error",
					Timestamp: time.Now().Unix(),
				})
			}
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.GET("/push/:pushId/:title/:text", func(c echo.Context) error {
		pushId := c.Param("pushId")
		title := c.Param("title")
		text := c.Param("text")
		if user, ok := users.FindUserByPushId(pushId); ok {
			if err := user.Send(title, text, false); err != nil {
				log.Errorf("send message to %d error: %s", user.PushId, err)
				return c.JSON(http.StatusInternalServerError, back{
					Code:      500,
					Msg:       "send message error",
					Timestamp: time.Now().Unix(),
				})
			}
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.GET("/push/:pushId/copy/:text", func(c echo.Context) error {
		pushId := c.Param("pushId")
		text := c.Param("text")
		if user, ok := users.FindUserByPushId(pushId); ok {
			if err := user.Send("", text, true); err != nil {
				log.Errorf("send message to %d error: %s", user.PushId, err)
				return c.JSON(http.StatusInternalServerError, back{
					Code:      500,
					Msg:       "send message error",
					Timestamp: time.Now().Unix(),
				})
			}
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.GET("/push/:pushId/copy/:title/:text", func(c echo.Context) error {
		pushId := c.Param("pushId")
		title := c.Param("title")
		text := c.Param("text")
		if user, ok := users.FindUserByPushId(pushId); ok {
			if err := user.Send(title, text, true); err != nil {
				log.Errorf("send message to %d error: %s", user.PushId, err)
				return c.JSON(http.StatusInternalServerError, back{
					Code:      500,
					Msg:       "send message error",
					Timestamp: time.Now().Unix(),
				})
			}
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.POST("/push/:pushId", func(c echo.Context) error {
		if user, ok := users.FindUserByPushId(c.Param("pushId")); ok {
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
			if err := user.Send(data.Title, tmp, data.Copy); err != nil {
				log.Errorf("send message to %d error: %s", user.PushId, err)
				return c.JSON(http.StatusInternalServerError, back{
					Code:      500,
					Msg:       "send message error",
					Timestamp: time.Now().Unix(),
				})
			}
			return success(c)
		} else {
			return notFound(c)
		}
	})

	//*标题 title*
	//_斜体 italic text_
	//[显示文本](网址) [text](URL)
	//`内联固定宽度的代码 inline fixed-width code`
	//```预先格式化的 固定宽度的代码块 pre-formatted fixed-width code block```
	e.POST("/tg-format/:pushId", func(c echo.Context) error {
		pushId := c.Param("pushId")
		if user, ok := users.FindUserByPushId(pushId); ok {
			data := struct {
				Text      string `json:"text"`
				ChatId    int64  `json:"chat_id"`
				ParseMode string `json:"parse_mode"`
			}{}
			if err := c.Bind(&data); err != nil {
				return c.JSON(http.StatusBadRequest, back{
					Code:      400,
					Msg:       "bad request",
					Timestamp: time.Now().Unix(),
				})
			}
			msg := tgbotapi.NewMessage(user.TgId, data.Text)
			msg.ParseMode = data.ParseMode
			if _, err := bot.Send(msg); err != nil {
				log.Errorf("send tg message to %d error: %s", user.TgId, err)
				return c.JSON(http.StatusInternalServerError, back{
					Code:      500,
					Msg:       "send message error",
					Timestamp: time.Now().Unix(),
				})
			}
			return success(c)
		} else {
			return notFound(c)
		}
	})

	// 兼容PushPeer格式
	e.POST("/message/push", func(c echo.Context) error {
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
		if user, ok := users.FindUserByPushId(data.PushKey); ok {
			if err := user.Send(data.Text, data.Desp, true); err != nil {
				log.Errorf("send message to %d error: %s", user.TgId, err)
				return sendError(c)
			}
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", deploy.Port)))
}

func sendError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, back{
		Code:      500,
		Msg:       "send message error",
		Timestamp: time.Now().Unix(),
	})
}

func success(c echo.Context) error {
	return c.JSON(http.StatusOK, back{
		Code:      200,
		Msg:       "success",
		Timestamp: time.Now().Unix(),
	})
}
func notFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, back{
		Code:      404,
		Msg:       "user not found",
		Timestamp: time.Now().Unix(),
	})
}
