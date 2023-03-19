package api

import (
	"fmt"
	"github.com/feranydev/mini-push/usergroup"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/feranydev/mini-push/config"
)

type back struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
}

func Start(bot *tgbotapi.BotAPI) {
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

	e.GET("/push/:pushId/:text", text)

	e.GET("/push/:pushId/:title/:text", textAndTitle)

	e.GET("/push/:pushId/copy/:text", textCopy)

	e.GET("/push/:pushId/copy/:title/:text", textAndTitleCopy)

	e.POST("/push/:pushId", post)

	// 兼容PushPeer格式
	e.POST("/message/push", pushDeer)

	//*标题 title*
	//_斜体 italic text_
	//[显示文本](网址) [text](URL)
	//`内联固定宽度的代码 inline fixed-width code`
	//```预先格式化的 固定宽度的代码块 pre-formatted fixed-width code block```
	e.POST("/tg-format/:pushId", func(c echo.Context) error {
		pushId := c.Param("pushId")
		if user, ok := usergroup.DefaultUserGroup.FindUserByPushId(pushId); ok {
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
			return success(c, "tg 格式发送没有id")
		} else {
			return userNotFound(c)
		}
	})

	e.GET("/api/get-message-info/:messageId", getMessageInfo)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Deploy.Port)))
}

func sendError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, back{
		Code:      500,
		Msg:       "send message error",
		Timestamp: time.Now().Unix(),
	})
}

func success(c echo.Context, msgId string) error {
	return c.JSON(http.StatusOK, back{
		Code:      200,
		Msg:       msgId,
		Timestamp: time.Now().Unix(),
	})
}

func backJson(c echo.Context, msg interface{}) error {
	return c.JSON(http.StatusOK, struct {
		Code      int         `json:"code"`
		Msg       interface{} `json:"msg"`
		Timestamp int64       `json:"timestamp"`
	}{
		Code:      200,
		Msg:       msg,
		Timestamp: time.Now().Unix(),
	})
}

func userNotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, back{
		Code:      404,
		Msg:       "user not found",
		Timestamp: time.Now().Unix(),
	})
}

func notFound(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotFound, back{
		Code:      404,
		Msg:       msg,
		Timestamp: time.Now().Unix(),
	})
}
