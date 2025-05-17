package handler

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/state"

	tele "gopkg.in/telebot.v4"
)

func (h *UserHandler) DeleteBookingInit(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}
	filter, ok := c.Get("filter").(string)
	if !ok {
		return fmt.Errorf("Failed get filter from context")
	}
	_, exists, err := h.eventRepository.Get(eventID)
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
	_, isBooked, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !isBooked {
		return fmt.Errorf("Failed get booking, booking not exists")
	}
	user.State = state.DeleteBookingWaitInputAreYouSure
	user.Context["eventID"] = eventID
	user.Context["filter"] = filter
	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}
	return c.Send(message.DeleteBookingWaitInputAreYouSureMessage)
}

func (h *UserHandler) DeleteBookingAreYouSureInput(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}
	eventID, ok := user.Context["eventID"].(string)
	if !ok {
		return fmt.Errorf("Failed get event id from user context")
	}
	filter, ok := user.Context["filter"].(string)
	if !ok {
		return fmt.Errorf("Failed get event id from user context")
	}
	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		user.State = state.Init
		user.Context = make(map[string]any)
		exists, err = h.userRepository.Update(&user)
		if err != nil {
			return fmt.Errorf("Failed update user: %s", err)
		}
		if !exists {
			return fmt.Errorf("Failed update user, user not exists")
		}
		return c.Send(message.EventNotFoundMessage)
	}
	if c.Message().Text == "да" {
		exists, err := h.bookingRepository.Delete(user.UserID, eventID)
		if err != nil {
			return fmt.Errorf("Failed delete booking: %s", err)
		}
		if !exists {
			return fmt.Errorf("Failed delete booking, booking not exists")
		}
	} else if c.Message().Text != "нет" {
		return c.Send(message.DeleteBookingWaitInputAreYouSureInvalidInputMessage)
	}
	user.State = state.Init
	user.Context = make(map[string]any)
	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}
	if c.Message().Text == "да" {
		err := c.Send(message.DeleteBookingSuccessMessage)
		if err != nil {
			return err
		}
		content := message.EventMessageContent(event, false, filter)
		return c.Send(content[0], content[1:]...)
	}
	return nil
}
