package handler

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/state"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) SendNotification(c tele.Context) error {
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

	var (
		sendNotificationOffline bool
		sendNotificationOnline  bool
	)
	bookings, err := h.bookingRepository.GetByEventID(eventID)
	if err != nil {
		return fmt.Errorf("Failed get bookings: %s", err)
	}
	var (
		bookingsOfflineExists bool
		bookingsOnlineExists  bool
	)
	for _, booking := range bookings {
		if booking.Offline {
			bookingsOfflineExists = true
		}
		if booking.Online {
			bookingsOnlineExists = true
		}
		if bookingsOfflineExists && bookingsOnlineExists {
			break
		}
	}
	if event.Offline || bookingsOfflineExists {
		sendNotificationOffline = true
	}
	if event.Online || bookingsOnlineExists {
		sendNotificationOnline = true
	}

	content := message.SendNotificationMessageContent(
		event,
		createdBy,
		updatedBy,
		eventURL,
		sendNotificationOffline,
		sendNotificationOnline)
	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}

func (h *AdminHandler) SendNotificationInit(c tele.Context) error {
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

	format, ok := c.Get("format").(string)
	if !ok {
		return fmt.Errorf("Failed get format from context")
	}
	if format != "offline" && format != "online" {
		return fmt.Errorf("Unexpected format value: %s", format)
	}

	allBookings, err := h.bookingRepository.GetByEventID(eventID)
	if err != nil {
		return fmt.Errorf("Failed get bookings: %s", err)
	}

	bookings := make([]*entity.Booking, 0)
	if format == "offline" {
		for _, booking := range allBookings {
			if booking.Offline {
				bookings = append(bookings, booking)
			}
		}
		if len(bookings) == 0 {
			if event.Offline && !event.Online {
				return c.Send(message.BookingsNotFoundMessage)
			} else {
				return c.Send(message.BookingsOfflineNotFoundMessage)
			}
		}
	} else {
		for _, booking := range allBookings {
			if booking.Online {
				bookings = append(bookings, booking)
			}
		}
		if len(bookings) == 0 {
			if event.Online && !event.Offline {
				return c.Send(message.BookingsNotFoundMessage)
			} else {
				return c.Send(message.BookingsOnlineNotFoundMessage)
			}
		}
	}

	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	user.State = state.SendNotificationWaitInputPhotoText
	user.Context["eventID"] = eventID
	user.Context["format"] = format

	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.SendNotificationWaitInputPhotoTextMessage)
}

func (h *AdminHandler) SendNotificationPhotoTextInput(c tele.Context) error {
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

	format, ok := user.Context["format"].(string)
	if !ok {
		return fmt.Errorf("Failed get format from user context")
	}
	if format != "offline" && format != "online" {
		return fmt.Errorf("Unexpected format value: %s", format)
	}

	content := make([]any, 0)
	if c.Message().Photo != nil {
		photo := c.Message().Photo
		photo.Caption = c.Message().Caption
		entities := c.Message().Entities
		content = append(content, photo, entities)
	} else {
		text := c.Message().Text
		entities := c.Message().Entities
		content = append(content, text, entities)
	}

	allBookings, err := h.bookingRepository.GetByEventID(eventID)
	if err != nil {
		return fmt.Errorf("Failed get bookings: %s", err)
	}

	bookings := make([]*entity.Booking, 0)
	if format == "offline" {
		for _, booking := range allBookings {
			if booking.Offline {
				bookings = append(bookings, booking)
			}
		}
		if len(bookings) == 0 {
			if event.Offline && !event.Online {
				return c.Send(message.BookingsNotFoundMessage)
			} else {
				return c.Send(message.BookingsOfflineNotFoundMessage)
			}
		}
	} else {
		for _, booking := range allBookings {
			if booking.Online {
				bookings = append(bookings, booking)
			}
		}
		if len(bookings) == 0 {
			if event.Online && !event.Offline {
				return c.Send(message.BookingsNotFoundMessage)
			} else {
				return c.Send(message.BookingsOnlineNotFoundMessage)
			}
		}
	}

	for _, booking := range bookings {
		_, err := c.Bot().Send(&tele.User{ID: booking.UserID}, content[0], content[1:]...)
		if err != nil {
			return err
		}
	}

	user.State = state.Init
	delete(user.Context, "eventID")
	exists, err = h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	return c.Send(message.SendNotificationSuccessMessage)
}
