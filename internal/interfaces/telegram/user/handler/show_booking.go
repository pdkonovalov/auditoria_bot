package handler

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/message"

	tele "gopkg.in/telebot.v4"
)

func (h *UserHandler) ShowBooking(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}
	filter, ok := c.Get("filter").(string)
	if !ok {
		return fmt.Errorf("Failed get filter from context")
	}
	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return err
	}
	if !exists {
		return c.Send(message.EventNotFoundMessage)
	}
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}
	booking, isBooked, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !isBooked {
		return fmt.Errorf("Failed get booking, booking not exists")
	}
	content := message.ShowBookingMessageContent(event, booking, filter)
	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}
