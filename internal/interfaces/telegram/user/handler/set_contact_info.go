package handler

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/state"

	tele "gopkg.in/telebot.v4"
)

func (h *UserHandler) SetContactInfo(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}
	content := message.SetContactInfoMessageContent(&user)
	return c.Send(content[0], content[1:]...)
}

func (h *UserHandler) EditContactInfoInit(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}
	user.State = state.EditContactInfoWaitInput
	exists, err := h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}
	return c.Send(message.EditContactInfoWaitInputMessage, message.EditContactInfoWaitInputReplyKeyboard(user.Username))
}

func (h *UserHandler) EditContactInfoInput(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}
	if c.Message().Contact != nil {
		user.ContactInfo = fmt.Sprintf("+%s", c.Message().Contact.PhoneNumber)
	} else if c.Message().Text == message.EditContactInfoWaitInputReplyKeyboardTelegramText {
		if user.Username == "" {
			return c.Send(message.EditContactInfoWaitInputTelegramNotExists, message.EditContactInfoWaitInputReplyKeyboard(user.Username))
		}
		user.ContactInfo = fmt.Sprintf("https://t.me/%s", user.Username)
	} else {
		text := c.Message().Text
		ok = h.validator.ContactInfo(text)
		if !ok {
			return c.Send(message.EditContactInfoInvalidInputMessage, message.EditContactInfoWaitInputReplyKeyboard(user.Username))
		}
		user.ContactInfo = text
	}
	user.State = state.Init
	user.Context = make(map[string]any)
	exists, err := h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}
	err = c.Send(message.EditContactInfSuccessMessage)
	if err != nil {
		return err
	}
	content := message.SetContactInfoMessageContent(&user)
	return c.Send(content[0], content[1:]...)
}
