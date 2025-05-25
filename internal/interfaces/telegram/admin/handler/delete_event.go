package handler

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/state"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) DeleteEventInit(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}

	_, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return c.Send(message.EventNotFoundMessage)
	}

	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	user.State = state.DeleteEventWaitInputAreYouSure
	user.Context["eventID"] = eventID

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.DeleteEventWaitInputAreYouSureMessage)
}

func (h *AdminHandler) DeleteEventAreYouSureInput(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	eventID, ok := user.Context["eventID"].(string)
	if !ok {
		return fmt.Errorf("Failed get event id from user context")
	}
	_, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return c.Send(message.EventNotFoundMessage)
	}

	var needDelete bool
	switch c.Message().Text {
	case "да":
		needDelete = true
	case "нет":
		needDelete = false
	default:
		return c.Send(message.DeleteEventWaitInputAreYouSureInvalidInputMessage)
	}

	if needDelete {
		exists, err := h.eventRepository.Delete(eventID)
		if err != nil {
			return fmt.Errorf("Failed delete event: %s", err)
		}
		if !exists {
			return c.Send(message.EventNotFoundMessage)
		}
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

	if needDelete {
		return c.Send(message.DeleteEventSuccessMessage)
	}
	return nil
}
