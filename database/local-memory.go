package database

import (
	"github.com/google/uuid"
)

var lm map[string]sqlMessage

func lmStart() {
	lm = make(map[string]sqlMessage)
}

func lmSet(messageId uuid.UUID, data sqlMessage) bool {
	lm[messageId.String()] = data
	return true
}

func lmGet(msgId uuid.UUID) (sqlMessage, bool) {
	message, ok := lm[msgId.String()]
	return message, ok
}
