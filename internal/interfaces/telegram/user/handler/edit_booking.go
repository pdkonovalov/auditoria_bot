package handler

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/state"
	tele "gopkg.in/telebot.v4"
)

func (h *UserHandler) EditBooking(c tele.Context) error {
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
		return err
	}
	if !exists {
		return c.Send(message.EventNotFoundMessage)
	}

	booking, isBooked, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !isBooked {
		return fmt.Errorf("Failed get booking, booking not exists")
	}

	filter, ok := c.Get("filter").(string)
	if !ok {
		return fmt.Errorf("Failed get filter from context")
	}

	content := message.EditBookingMessageContent(event, booking, filter)
	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}

func (h *UserHandler) EditBookingFormatInit(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}
	_, exists, err := h.eventRepository.Get(eventID)
	if err != nil {
		return err
	}
	if !exists {
		return c.Send(message.EventNotFoundMessage)
	}

	_, isBooked, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !isBooked {
		return fmt.Errorf("Failed get booking, booking not exists")
	}

	filter, ok := c.Get("filter").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}

	user.State = state.EditBookingWaitInputFormat
	user.Context["eventID"] = eventID
	user.Context["filter"] = filter

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.EditBookingWaitInputFormatMessage, message.EditBookingWaitInputFormatReplyKeyboard)
}

func (h *UserHandler) EditBookingFormatInput(c tele.Context) error {
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

	booking, exists, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get booking, booking not exists")
	}

	filter, ok := user.Context["filter"].(string)
	if !ok {
		return fmt.Errorf("Failed get filter from user context")
	}

	switch c.Message().Text {
	case message.BookingWaitInputFormatReplyKeyboardOfflineText:
		booking.Offline = true
		booking.Online = false
	case message.BookingWaitInputFormatReplyKeyboardOnlineText:
		booking.Offline = false
		booking.Online = true
	default:
		return c.Send(message.BookingWaitInputFormatInvalidInputMessage)
	}

	booking.UpdatedAt = new(time.Time)
	*booking.UpdatedAt = time.Now()

	exists, err = h.bookingRepository.Update(booking)
	if err != nil {
		return fmt.Errorf("Failed update booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update booking, booking not exists")
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

	err = c.Send(message.EditBookingSuccessMessage)
	if err != nil {
		return err
	}
	content := message.EditBookingMessageContent(event, booking, filter)

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}

func (h *UserHandler) EditBookingPaymentInit(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}
	filter, ok := c.Get("filter").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
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
	_, isBooked, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !isBooked {
		return fmt.Errorf("Failed get booking, booking not exists")
	}
	user.State = state.EditBookingWaitInputPayment
	user.Context["eventID"] = eventID
	user.Context["filter"] = filter
	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}
	content := message.EditBookingWaitInputPaymentMessageContent(event)
	return c.Send(content[0], content[1:]...)
}

func (h *UserHandler) EditBookingPaymentInput(c tele.Context) error {
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
		return fmt.Errorf("Failed get filter from user context")
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
	booking, exists, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get booking, booking not exists")
	}
	if c.Message().Photo != nil {
		booking.Payment = true
		booking.PaymentPhoto = c.Message().Photo
		booking.PaymentDocument = nil
	} else if c.Message().Document != nil {
		booking.Payment = true
		booking.PaymentPhoto = nil
		booking.PaymentDocument = c.Message().Document
	} else if c.Message().Text == message.BookingWaitInputPaymentReplyKeyboardText {
		booking.Payment = false
	} else {
		return c.Send(message.BookingWaitInputPaymentInvalidInputMessage)
	}
	booking.UpdatedAt = new(time.Time)
	*booking.UpdatedAt = time.Now()
	exists, err = h.bookingRepository.Update(booking)
	if err != nil {
		return fmt.Errorf("Failed update booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update booking, booking not exists")
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
	err = c.Send(message.EditBookingSuccessMessage)
	if err != nil {
		return err
	}
	content := message.EditBookingMessageContent(event, booking, filter)

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}

func (h *UserHandler) EditBookingAdditionalInfoInit(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}
	filter, ok := c.Get("filter").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
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
	user.State = state.EditBookingWaitInputAdditionalInfo
	user.Context["eventID"] = eventID
	user.Context["filter"] = filter
	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}
	return c.Send(message.EditBookingWaitInputAdditionalInfoMessage, message.EditBookingWaitInputAdditionalInfoReplyKeyboard)
}

func (h *UserHandler) EditBookingAdditionalInfoInput(c tele.Context) error {
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
	booking, exists, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get booking, booking not exists")
	}
	if c.Message().Text == message.BookingWaitInputAdditionalInfoReplyKeyboardText {
		booking.Text = ""
	} else {
		booking.Text = c.Message().Text
	}
	booking.UpdatedAt = new(time.Time)
	*booking.UpdatedAt = time.Now()
	exists, err = h.bookingRepository.Update(booking)
	if err != nil {
		return fmt.Errorf("Failed update booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update booking, booking not exists")
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
	err = c.Send(message.EditBookingSuccessMessage)
	if err != nil {
		return err
	}

	content := message.EditBookingMessageContent(event, booking, filter)

	time.Sleep(time.Second)

	return c.Send(content[0], content[1:]...)
}
