package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/state"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) EditEvent(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}
	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return c.Send(message.EventNotFoundMessage)
	}
	createdBy, exists, err := h.userRepository.Get(event.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed get created by: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get created by, user not exists")
	}
	var updatedBy *entity.User
	if event.UpdatedBy != nil {
		updatedBy, exists, err = h.userRepository.Get(*event.UpdatedBy)
		if err != nil {
			return fmt.Errorf("Failed get updated by: %s", err)
		}
		if !exists {
			return fmt.Errorf("Failed get updated by, user not exists")
		}
	}
	eventURL := h.generateBotUrl(eventID)

	bookingsOffline, err := h.bookingRepository.GetByEventID(eventID, true, false)
	if err != nil {
		return fmt.Errorf("Failed get offline bookings: %s", err)
	}

	bookingsOnline, err := h.bookingRepository.GetByEventID(eventID, false, true)
	if err != nil {
		return fmt.Errorf("Failed get online bookings: %s", err)
	}

	content := message.EditEventMessageContent(event, createdBy, updatedBy, eventURL, len(bookingsOffline), len(bookingsOnline))
	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}

func (h *AdminHandler) EditEventFormatInit(c tele.Context) error {
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

	user.State = state.EditEventWaitInputFormat
	user.Context["eventID"] = eventID

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputFormatMessage, message.WaitInputFormatReplyKeyboard)
}

func (h *AdminHandler) EditEventFormatInput(c tele.Context) error {
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
		return c.Send(message.EventNotFoundMessage)
	}

	var offline, offlinePaid, online, onlinePaid bool
	switch c.Message().Text {
	case message.WaitInputFormatReplyKeyboardButtonOffline:
		offline = true
		if event.Offline {
			offlinePaid = event.OfflinePaid
		} else {
			offlinePaid = event.OnlinePaid
		}
	case message.WaitInputFormatReplyKeyboardButtonOnline:
		online = true
		if event.Online {
			onlinePaid = event.OnlinePaid
		} else {
			onlinePaid = event.OfflinePaid
		}
	case message.WaitInputFormatReplyKeyboardButtonOfflineOnline:
		offline = true
		if event.Offline {
			offlinePaid = event.OfflinePaid
		} else {
			offlinePaid = event.OnlinePaid
		}
		online = true
		if event.Online {
			onlinePaid = event.OnlinePaid
		} else {
			onlinePaid = event.OfflinePaid
		}
	default:
		return c.Send(message.WaitInputFormatInvalidInputMessage)
	}

	event.Offline = offline
	event.OfflinePaid = offlinePaid
	event.Online = online
	event.OnlinePaid = onlinePaid

	timeNow := time.Now()
	event.UpdatedAt = &timeNow
	event.UpdatedBy = &user.UserID

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

	err = c.Send(message.EditEventSuccessMessage)
	if err != nil {
		return err
	}

	createdBy, exists, err := h.userRepository.Get(event.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed get created by: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get created by, user not exists")
	}

	bookingsOffline, err := h.bookingRepository.GetByEventID(eventID, true, false)
	if err != nil {
		return fmt.Errorf("Failed get offline bookings: %s", err)
	}

	bookingsOnline, err := h.bookingRepository.GetByEventID(eventID, false, true)
	if err != nil {
		return fmt.Errorf("Failed get online bookings: %s", err)
	}

	eventURL := h.generateBotUrl(eventID)
	content := message.EditEventMessageContent(event, createdBy, &user, eventURL, len(bookingsOffline), len(bookingsOnline))

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}

func (h *AdminHandler) EditEventPaidInit(c tele.Context) error {
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

	format, ok := c.Get("format").(string)
	if !ok {
		return fmt.Errorf("Failed get format from context")
	}

	user.State = state.EditEventWaitInputPaid
	user.Context["eventID"] = eventID
	user.Context["format"] = format

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputPaidMessage, message.WaitInputPaidReplyKeyboard)
}

func (h *AdminHandler) EditEventPaidInput(c tele.Context) error {
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
	} else if format == "online" {
		event.OnlinePaid = paid
	} else {
		return fmt.Errorf("Unexpected format in user context: %s", format)
	}

	timeNow := time.Now()
	event.UpdatedAt = &timeNow
	event.UpdatedBy = &user.UserID

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

	err = c.Send(message.EditEventSuccessMessage)
	if err != nil {
		return err
	}

	createdBy, exists, err := h.userRepository.Get(event.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed get created by: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get created by, user not exists")
	}

	bookingsOffline, err := h.bookingRepository.GetByEventID(eventID, true, false)
	if err != nil {
		return fmt.Errorf("Failed get offline bookings: %s", err)
	}

	bookingsOnline, err := h.bookingRepository.GetByEventID(eventID, false, true)
	if err != nil {
		return fmt.Errorf("Failed get online bookings: %s", err)
	}

	eventURL := h.generateBotUrl(eventID)
	content := message.EditEventMessageContent(event, createdBy, &user, eventURL, len(bookingsOffline), len(bookingsOnline))

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}

func (h *AdminHandler) EditEventTitleInit(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}

	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	user.State = state.EditEventWaitInputTitle
	user.Context["eventID"] = eventID

	exists, err := h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputTitleMessage)
}

func (h *AdminHandler) EditEventTitleInput(c tele.Context) error {
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

	timeNow := time.Now()
	event.UpdatedAt = &timeNow
	event.UpdatedBy = &user.UserID

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

	err = c.Send(message.EditEventSuccessMessage)
	if err != nil {
		return err
	}

	createdBy, exists, err := h.userRepository.Get(event.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed get created by: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get created by, user not exists")
	}

	bookingsOffline, err := h.bookingRepository.GetByEventID(eventID, true, false)
	if err != nil {
		return fmt.Errorf("Failed get offline bookings: %s", err)
	}

	bookingsOnline, err := h.bookingRepository.GetByEventID(eventID, false, true)
	if err != nil {
		return fmt.Errorf("Failed get online bookings: %s", err)
	}

	eventURL := h.generateBotUrl(eventID)
	content := message.EditEventMessageContent(event, createdBy, &user, eventURL, len(bookingsOffline), len(bookingsOnline))

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}

func (h *AdminHandler) EditEventTimeInit(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}

	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	user.State = state.EditEventWaitInputTime
	user.Context["eventID"] = eventID

	exists, err := h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputTimeMessage)
}

func (h *AdminHandler) EditEventTimeInput(c tele.Context) error {
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

	newTime, err := time.ParseInLocation("15:04 02.01.2006", c.Message().Text, h.location)
	if err != nil {
		return c.Send(message.WaitInputTimeInvalidInputMessage)
	}

	event.Time = newTime

	timeNow := time.Now()
	event.UpdatedAt = &timeNow
	event.UpdatedBy = &user.UserID

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

	err = c.Send(message.EditEventSuccessMessage)
	if err != nil {
		return err
	}

	createdBy, exists, err := h.userRepository.Get(event.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed get created by: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get created by, user not exists")
	}

	bookingsOffline, err := h.bookingRepository.GetByEventID(eventID, true, false)
	if err != nil {
		return fmt.Errorf("Failed get offline bookings: %s", err)
	}

	bookingsOnline, err := h.bookingRepository.GetByEventID(eventID, false, true)
	if err != nil {
		return fmt.Errorf("Failed get online bookings: %s", err)
	}

	eventURL := h.generateBotUrl(eventID)
	content := message.EditEventMessageContent(event, createdBy, &user, eventURL, len(bookingsOffline), len(bookingsOnline))

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}

func (h *AdminHandler) EditEventPaymentDetailsInit(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}

	event, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return fmt.Errorf("Failed get event: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get event, event not exists")
	}

	if !event.OfflinePaid && !event.OnlinePaid {
		return nil
	}

	user.State = state.EditEventWaitInputPaymentDetails
	user.Context["eventID"] = eventID

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputPaymentDetailsMessage, message.WaitInputPaymentDetailsReplyKeyboard(h.defaultPaymentDetails))
}

func (h *AdminHandler) EditEventPaymentDetailsInput(c tele.Context) error {
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

	parts := strings.Split(c.Message().Text, " ")
	if len(parts) != 3 ||
		len(parts[0]) == 0 ||
		len(parts[1]) == 0 ||
		len(parts[2]) == 0 {
		return c.Send(message.WaitInputPaymentDetailsInvalidInputMessage)
	}

	event.PaymentDetailsFirstName = parts[0]
	event.PaymentDetailsLastName = parts[1]
	event.PaymentDetailsAccount = parts[2]

	timeNow := time.Now()
	event.UpdatedAt = &timeNow
	event.UpdatedBy = &user.UserID

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

	err = c.Send(message.EditEventSuccessMessage)
	if err != nil {
		return err
	}

	createdBy, exists, err := h.userRepository.Get(event.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed get created by: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get created by, user not exists")
	}

	bookingsOffline, err := h.bookingRepository.GetByEventID(eventID, true, false)
	if err != nil {
		return fmt.Errorf("Failed get offline bookings: %s", err)
	}

	bookingsOnline, err := h.bookingRepository.GetByEventID(eventID, false, true)
	if err != nil {
		return fmt.Errorf("Failed get online bookings: %s", err)
	}

	eventURL := h.generateBotUrl(eventID)
	content := message.EditEventMessageContent(event, createdBy, &user, eventURL, len(bookingsOffline), len(bookingsOnline))

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}

func (h *AdminHandler) EditEventPhotoTextInit(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}

	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	user.State = state.EditEventWaitInputPhotoText
	user.Context["eventID"] = eventID

	exists, err := h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.WaitInputPhotoTextMessage)
}

func (h *AdminHandler) EditEventPhotoTextInput(c tele.Context) error {
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
		event.Photo = c.Message().Photo
		event.Text = c.Message().Caption
		event.Entities = c.Message().CaptionEntities
	} else {
		event.Photo = nil
		event.Text = c.Message().Text
		event.Entities = c.Message().Entities
	}

	timeNow := time.Now()
	event.UpdatedAt = &timeNow
	event.UpdatedBy = &user.UserID

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

	err = c.Send(message.EditEventSuccessMessage)
	if err != nil {
		return err
	}
	content := message.EventPhotoTextMessageContent(event)

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}
