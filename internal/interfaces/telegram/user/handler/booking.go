package handler

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/state"

	tele "gopkg.in/telebot.v4"
)

func (h *UserHandler) BookingInit(c tele.Context) error {
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

	filter, ok := c.Get("filter").(string)
	if !ok {
		return fmt.Errorf("Failed get filter from context")
	}

	_, isBooked, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if isBooked {
		return c.Send(message.BookingAlredyBookedMessage)
	}

	booking := entity.Booking{
		EventID: eventID,
		UserID:  user.UserID,
		Draft:   true,
	}

	if len(user.ContactInfo) == 0 {
		user.State = state.BookingWaitInputContactInfo
	} else if event.Offline && event.Online {
		user.State = state.BookingWaitInputFormat
	} else if event.Offline {
		booking.Offline = true
		if event.OfflinePaid {
			user.State = state.BookingWaitInputPayment
		} else {
			user.State = state.BookingWaitInputAdditionalInfo
		}
	} else if event.Online {
		booking.Online = true
		if event.OnlinePaid {
			user.State = state.BookingWaitInputPayment
		} else {
			user.State = state.BookingWaitInputAdditionalInfo
		}
	} else {
		return fmt.Errorf("Invalid event, offline and online is false")
	}

	isBooked, err = h.bookingRepository.Create(&booking)
	if err != nil {
		return fmt.Errorf("Failed create booking: %s", err)
	}
	if isBooked {
		return fmt.Errorf("Failed create booking, booking exists")
	}

	user.Context["eventID"] = eventID
	user.Context["filter"] = filter

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	switch user.State {
	case state.BookingWaitInputContactInfo:
		return c.Send(message.BookingWaitInputContactInfoMessage, message.BookingWaitInputContactInfoReplyKeyboard(user.Username))
	case state.BookingWaitInputFormat:
		return c.Send(message.BookingWaitInputFormatMessage, message.BookingWaitInputFormatReplyKeyboard)
	case state.BookingWaitInputPayment:
		content := message.BookingWaitInputPaymentMessageContent(event)
		return c.Send(content[0], content[1:]...)
	case state.BookingWaitInputAdditionalInfo:
		return c.Send(message.BookingWaitInputAdditionalInfoMessage, message.BookingWaitInputAdditionalInfoReplyKeyboard)
	}

	return fmt.Errorf("Unexpected user state: %s", user.State)
}

func (h *UserHandler) BookingContactInfoInput(c tele.Context) error {
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
		err := c.Send(message.BookingContactInfoSuccessMessage)
		if err != nil {
			return err
		}
		content := message.SetContactInfoMessageContent(&user)

		time.Sleep(time.Second)

		err = c.Send(content[0])
		if err != nil {
			return err
		}

		time.Sleep(time.Second)

		return c.Send(message.EventNotFoundMessage)
	}

	booking, exists, err := h.bookingRepository.Get(user.UserID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get booking, booking not exists")
	}

	if c.Message().Contact != nil {
		user.ContactInfo = fmt.Sprintf("+%s", c.Message().Contact.PhoneNumber)
	} else if c.Message().Text == message.BookingWaitInputContactInfoReplyKeyboardTelegramText {
		if user.Username == "" {
			return c.Send(message.BookingWaitInputContactInfoTelegramNotExists, message.BookingWaitInputContactInfoReplyKeyboard(user.Username))
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

	if event.Offline && event.Online {
		user.State = state.BookingWaitInputFormat
	} else if event.Offline {
		booking.Offline = true
		if event.OfflinePaid {
			user.State = state.BookingWaitInputPayment
		} else {
			user.State = state.BookingWaitInputAdditionalInfo
		}
	} else if event.Online {
		booking.Online = true
		if event.OnlinePaid {
			user.State = state.BookingWaitInputPayment
		} else {
			user.State = state.BookingWaitInputAdditionalInfo
		}
	} else {
		return fmt.Errorf("Invalid event, offline and online is false")
	}

	exists, err = h.bookingRepository.Update(booking)
	if err != nil {
		return fmt.Errorf("Failed update booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update booking, booking not exists")
	}

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	err = c.Send(message.BookingContactInfoSuccessMessage)
	if err != nil {
		return err
	}
	content := message.SetContactInfoMessageContent(&user)

	time.Sleep(time.Second)

	err = c.Send(content[0], content[1])
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	switch user.State {
	case state.BookingWaitInputFormat:
		return c.Send(message.BookingWaitInputFormatMessage, message.BookingWaitInputFormatReplyKeyboard)
	case state.BookingWaitInputPayment:
		content := message.BookingWaitInputPaymentMessageContent(event)
		return c.Send(content[0], content[1:]...)
	case state.BookingWaitInputAdditionalInfo:
		return c.Send(message.BookingWaitInputAdditionalInfoMessage, message.BookingWaitInputAdditionalInfoReplyKeyboard)
	}

	return fmt.Errorf("Unexpected user state: %s", user.State)
}

func (h *UserHandler) BookingFormatInput(c tele.Context) error {
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

	switch c.Message().Text {
	case message.BookingWaitInputFormatReplyKeyboardOfflineText:
		booking.Offline = true
		if event.OfflinePaid {
			user.State = state.BookingWaitInputPayment
		} else {
			user.State = state.BookingWaitInputAdditionalInfo
		}
	case message.BookingWaitInputFormatReplyKeyboardOnlineText:
		booking.Online = true
		if event.OnlinePaid {
			user.State = state.BookingWaitInputPayment
		} else {
			user.State = state.BookingWaitInputAdditionalInfo
		}
	default:
		return c.Send(message.BookingWaitInputFormatInvalidInputMessage)
	}

	exists, err = h.bookingRepository.Update(booking)
	if err != nil {
		return fmt.Errorf("Failed update booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update booking, booking not exists")
	}

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	switch user.State {
	case state.BookingWaitInputPayment:
		content := message.BookingWaitInputPaymentMessageContent(event)
		return c.Send(content[0], content[1:]...)
	case state.BookingWaitInputAdditionalInfo:
		return c.Send(message.BookingWaitInputAdditionalInfoMessage, message.BookingWaitInputAdditionalInfoReplyKeyboard)
	}

	return fmt.Errorf("Unexpected user state: %s", user.State)
}

func (h *UserHandler) BookingPaymentInput(c tele.Context) error {
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
		booking.Payment = c.Message().Photo
	} else if c.Message().Text == message.BookingWaitInputPaymentReplyKeyboardText {
		booking.Payment = nil
	} else {
		return c.Send(message.BookingWaitInputPaymentInvalidInputMessage)
	}

	exists, err = h.bookingRepository.Update(booking)
	if err != nil {
		return fmt.Errorf("Failed гupdate booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update booking, booking not exists")
	}

	user.State = state.BookingWaitInputAdditionalInfo

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.BookingWaitInputAdditionalInfoMessage, message.BookingWaitInputAdditionalInfoReplyKeyboard)
}

func (h *UserHandler) BookingAdditionalInfoInput(c tele.Context) error {
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

	if c.Message().Text == message.BookingWaitInputAdditionalInfoReplyKeyboardText {
		booking.Text = ""
	} else {
		booking.Text = c.Message().Text
	}

	booking.Draft = false
	booking.CreatedAt = time.Now()

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

	err = c.Send(message.BookingSuccessMessage)
	if err != nil {
		return err
	}

	content := message.EventMessageContent(event, true, filter)
	return c.Send(content[0], content[1:]...)
}
