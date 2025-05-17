package handler

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/message"

	tele "gopkg.in/telebot.v4"
)

func (h *UserHandler) MyEvents(c tele.Context) error {
	nowDateString := time.Now().Format("02.01.2006")
	nowDateTime, err := time.ParseInLocation("02.01.2006", nowDateString, h.location)
	if err != nil {
		return err
	}
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}
	event, exists, err := h.eventRepository.GetFirstAfter(nowDateTime.Add(-time.Minute), &user.UserID)
	if err != nil {
		return fmt.Errorf("Failed get first after event: %s", err)
	}
	if exists {
		c.Set("date", event.Time.Format("02.01.2006"))
		c.Set("filter", "my")
		return h.EventsByDate(c)
	}
	event, exists, err = h.eventRepository.GetFirstBefore(nowDateTime.Add(time.Minute), &user.UserID)
	if err != nil {
		return fmt.Errorf("Failed get first before event: %s", err)
	}
	if exists {
		c.Set("date", event.Time.Format("02.01.2006"))
		c.Set("filter", "my")
		return h.EventsByDate(c)
	}
	return c.Send(message.MyEventsNotFoundMessage)
}
