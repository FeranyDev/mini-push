package api

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/feranydev/push-server/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

//*标题 title*
//_斜体 italic text_
//[显示文本](网址) [text](URL)
//`内联固定宽度的代码 inline fixed-width code`
//```预先格式化的 固定宽度的代码块 pre-formatted fixed-width code block```

type back struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
}

func Start(bot *tgbotapi.BotAPI, t *config.T) {
	e := echo.New()

	e.HideBanner = true

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
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
		if user, ok := t.FindUserByPushId(pushId); ok {
			msg := tgbotapi.NewMessage(user.TgId, text)
			bot.Send(msg)
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.GET("/push/:pushId/:title/:text", func(c echo.Context) error {
		pushId := c.Param("pushId")
		title := c.Param("title")
		text := c.Param("text")
		if user, ok := t.FindUserByPushId(pushId); ok {
			data := fmt.Sprintf("*%s*\n%s", title, text)
			msg := tgbotapi.NewMessage(user.TgId, data)
			msg.ParseMode = "markdown"
			bot.Send(msg)
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.GET("/push/:pushId/copy/:text", func(c echo.Context) error {
		pushId := c.Param("pushId")
		text := c.Param("text")
		if user, ok := t.FindUserByPushId(pushId); ok {
			msg := tgbotapi.NewMessage(user.TgId, fmt.Sprintf("`%s`", text))
			msg.ParseMode = "markdown"
			bot.Send(msg)
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.GET("/push/:pushId/copy/:title/:text", func(c echo.Context) error {
		pushId := c.Param("pushId")
		title := c.Param("title")
		text := c.Param("text")
		if user, ok := t.FindUserByPushId(pushId); ok {
			msg := tgbotapi.NewMessage(user.TgId, fmt.Sprintf("*%s*\n`%s`", title, text))
			msg.ParseMode = "markdown"
			bot.Send(msg)
			return success(c)
		} else {
			return notFound(c)
		}
	})

	e.Logger.Fatal(e.Start(":8080"))
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
