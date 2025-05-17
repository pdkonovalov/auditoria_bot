package handler

import (
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) Hello(c tele.Context) error {
	return c.Send(message.HelloMessage, message.HelloEntities)
}
