package handler

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) Events(c tele.Context) error {
	nowDateString := time.Now().Format("02.01.2006")
	nowDateTime, err := time.ParseInLocation("02.01.2006", nowDateString, h.location)
	if err != nil {
		return err
	}
	event, exists, err := h.eventRepository.GetFirstAfter(nowDateTime.Add(-time.Minute), nil)
	if err != nil {
		return fmt.Errorf("Failed get first after event: %s", err)
	}
	if exists {
		c.Set("date", event.Time.Format("02.01.2006"))
		return h.EventsByDate(c)
	}
	event, exists, err = h.eventRepository.GetFirstBefore(nowDateTime.Add(time.Minute), nil)
	if err != nil {
		return fmt.Errorf("Failed get first before event: %s", err)
	}
	if exists {
		c.Set("date", event.Time.Format("02.01.2006"))
		return h.EventsByDate(c)
	}
	return c.Send(message.EventsNotFoundMessage)
}
