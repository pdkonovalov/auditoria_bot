package handler

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/message"

	tele "gopkg.in/telebot.v4"
)

func (h *UserHandler) Hello(c tele.Context) error {
	is_admin, ok := c.Get("is_admin").(bool)
	if !ok {
		return fmt.Errorf("Failed get is admin from context")
	}
	if is_admin {
		return c.Send(message.HelloAdminMessage, message.HelloEntities)
	}
	return c.Send(message.HelloMessage, message.HelloEntities)
}
