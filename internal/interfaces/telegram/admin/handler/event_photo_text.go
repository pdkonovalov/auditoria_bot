package handler

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) EventPhotoText(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}
	event, exist, err := h.eventRepository.Get(eventID)
	if err != nil {
		return err
	}
	if !exist {
		return c.Send(message.EventNotFoundMessage)
	}
	content := message.EventPhotoTextMessageContent(event)
	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}
