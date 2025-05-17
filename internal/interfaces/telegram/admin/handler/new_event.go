package handler

import (
	"fmt"
	"time"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/state"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) NewEventInit(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	event := entity.Event{
		Draft:     true,
		CreatedBy: user.UserID,
	}
	eventID, err := h.eventRepository.Create(&event)
	if err != nil {
		return fmt.Errorf("Failed create event: %s", err)
	}

	user.State = state.NewEventWaitInputFormat
	user.Context["eventID"] = eventID
	exists, err := h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputFormatMessage, message.WaitInputFormatReplyKeyboard)
}

func (h *AdminHandler) NewEventFormatInput(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	eventID, ok := user.Context["eventID"].(string)
	if !ok {
		return fmt.Errorf("Failed get event id from user context")
	}
	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get event, event not exists")
	}

	var nextInputFormat string

	var offline, online bool
	switch c.Message().Text {
	case message.WaitInputFormatReplyKeyboardButtonOffline:
		offline = true
		nextInputFormat = "offline"
	case message.WaitInputFormatReplyKeyboardButtonOnline:
		online = true
		nextInputFormat = "online"
	case message.WaitInputFormatReplyKeyboardButtonOfflineOnline:
		offline = true
		online = true
		nextInputFormat = "offline"
	default:
		return c.Send(message.WaitInputFormatInvalidInputMessage)
	}

	event.Offline = offline
	event.Online = online

	exists, err = h.eventRepository.Update(event)
	if err != nil {
		return fmt.Errorf("Failed update event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update event, event not exists")
	}

	user.State = state.NewEventWaitInputPaid
	user.Context["format"] = nextInputFormat

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	if offline && online {
		return c.Send(message.WaitInputOfflinePaidMessage, message.WaitInputPaidReplyKeyboard)
	}
	return c.Send(message.WaitInputPaidMessage, message.WaitInputPaidReplyKeyboard)
}

func (h *AdminHandler) NewEventPaidInput(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	eventID, ok := user.Context["eventID"].(string)
	if !ok {
		return fmt.Errorf("Failed get event id from user context")
	}
	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get event, event not exists")
	}

	format, ok := user.Context["format"].(string)
	if !ok {
		return fmt.Errorf("Failed get format from user context")
	}

	var paid bool
	switch c.Message().Text {
	case message.WaitInputPaidReplyKeyboardButtonTrue:
		paid = true
	case message.WaitInputPaidReplyKeyboardButtonFalse:
	default:
		return c.Send(message.WaitInputPaidInvalidInputMessage)
	}

	if format == "offline" {
		event.OfflinePaid = paid
		if event.Online {
			user.State = state.NewEventWaitInputPaid
			user.Context["format"] = "online"
		} else {
			user.State = state.NewEventWaitInputTitle
			delete(user.Context, "format")
		}
	} else if format == "online" {
		event.OnlinePaid = paid
		user.State = state.NewEventWaitInputTitle
		delete(user.Context, "format")
	} else {
		return fmt.Errorf("Unexpected format in user context: %s", format)
	}

	exists, err = h.eventRepository.Update(event)
	if err != nil {
		return fmt.Errorf("Failed update event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update event, event not exists")
	}

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	if user.State == state.NewEventWaitInputPaid {
		return c.Send(message.WaitInputOnlinePaidMessage, message.WaitInputPaidReplyKeyboard)
	}
	return c.Send(message.WaitInputTitleMessage)
}

func (h *AdminHandler) NewEventTitleInput(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	eventID, ok := user.Context["eventID"].(string)
	if !ok {
		return fmt.Errorf("Failed get event id from user context")
	}
	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get event, event not exists")
	}

	event.Title = c.Message().Text

	exists, err = h.eventRepository.Update(event)
	if err != nil {
		return fmt.Errorf("Failed update event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update event, event not exists")
	}

	user.State = state.NewEventWaitInputTime
	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputTimeMessage)
}

func (h *AdminHandler) NewEventTimeInput(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	eventID, ok := user.Context["eventID"].(string)
	if !ok {
		return fmt.Errorf("Failed get event id from user context")
	}
	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get event, event not exists")
	}

	time, err := time.ParseInLocation("15:04 02.01.2006", c.Message().Text, h.location)
	if err != nil {
		return c.Send(message.WaitInputTimeInvalidInputMessage)
	}

	event.Time = time

	exists, err = h.eventRepository.Update(event)
	if err != nil {
		return fmt.Errorf("Failed update event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update event, event not exists")
	}

	user.State = state.NewEventWaitInputPhotoText
	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputPhotoTextMessage)
}

func (h *AdminHandler) NewEventPhotoTextInput(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	eventID, ok := user.Context["eventID"].(string)
	if !ok {
		return fmt.Errorf("Failed get event id from user context")
	}
	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get event, event not exists")
	}

	if c.Message().Photo != nil {
		captionLen := len(utf16.Encode([]rune(c.Message().Caption)))
		if captionLen > 1024 {
			return c.Send(message.WaitInputPhotoTextInvalidInputMessage(captionLen))
		}
		event.Photo = c.Message().Photo
		event.Text = c.Message().Caption
		event.Entities = c.Message().CaptionEntities
	} else {
		event.Text = c.Message().Text
	}
	event.Draft = false
	event.CreatedAt = time.Now()
	event.CreatedBy = user.UserID

	exists, err = h.eventRepository.Update(event)
	if err != nil {
		return fmt.Errorf("Failed update event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update event, event not exists")
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

	eventURL := h.generateBotUrl(eventID)
	return c.Send(message.NewEventSuccessMessage(eventURL))
}
