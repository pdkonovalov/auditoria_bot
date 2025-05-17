package handler

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) EventsByDate(c tele.Context) error {
	curDate, ok := c.Get("date").(string)
	if !ok {
		return fmt.Errorf("Failed get date from context")
	}
	curDateTime, err := time.ParseInLocation("02.01.2006", curDate, h.location)
	if err != nil {
		return err
	}
	events, err := h.eventRepository.GetByDate(curDateTime, nil)
	if len(events) == 0 {
		return c.Send(message.EventsNotFoundMessage)
	}
	var prevDate *string
	prevEvent, exists, err := h.eventRepository.GetFirstBefore(curDateTime, nil)
	if err != nil {
		return fmt.Errorf("Failed get first before event: %s", err)
	}
	if exists {
		prevDate = new(string)
		*prevDate = prevEvent.Time.Format("02.01.2006")
	}
	var nextDate *string
	nextEvent, exists, err := h.eventRepository.GetFirstAfter(curDateTime.Add(24*time.Hour-time.Minute), nil)
	if err != nil {
		return fmt.Errorf("Failed get first after event: %s", err)
	}
	if exists {
		nextDate = new(string)
		*nextDate = nextEvent.Time.Format("02.01.2006")
	}
	content, err := message.EventsByDateMessageContent(prevDate, curDate, nextDate, events)
	if err != nil {
		return err
	}
	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}
