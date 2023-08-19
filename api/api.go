package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/feranydev/mini-push/database"
)

func getMessageInfo(c echo.Context) error {
	messageId := c.Param("messageId")
	uid, err := uuid.Parse(messageId)
	if err != nil {
		return notFound(c, "This message ID is not legitimate")
	}
	if message, ok := database.GetMessageByMsgID(uid); ok {
		return backJson(c, message)
	} else {
		return notFound(c, "This message was not found")
	}

}
