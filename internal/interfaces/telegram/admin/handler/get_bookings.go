package handler

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) GetBookings(c tele.Context) error {
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

	content := message.GetBookingsMessageContent(
		event,
		createdBy,
		updatedBy,
		eventURL,
		len(bookingsOffline),
		len(bookingsOnline),
	)
	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}

func (h *AdminHandler) GetBookingsOffline(c tele.Context) error {
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

	bookingsOffline, err := h.bookingRepository.GetByEventID(eventID, true, false)
	if err != nil {
		return fmt.Errorf("Failed get offline bookings: %s", err)
	}

	if len(bookingsOffline) == 0 {
		if event.Offline && !event.Online {
			return c.Send(message.BookingsNotFoundMessage)
		} else {
			return c.Send(message.BookingsOfflineNotFoundMessage)
		}
	}

	for _, booking := range bookingsOffline {
		user, exists, err := h.userRepository.Get(booking.UserID)
		if err != nil {
			return fmt.Errorf("Failed get user: %s", err)
		}
		if !exists {
			return fmt.Errorf("Failed get user, user not exists")
		}
		content := message.BookingMessageContent(booking, user)
		err = c.Send(content[0], content[1:]...)
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (h *AdminHandler) GetBookingsOnline(c tele.Context) error {
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

	bookingsOnline, err := h.bookingRepository.GetByEventID(eventID, false, true)
	if err != nil {
		return fmt.Errorf("Failed get online bookings: %s", err)
	}

	if len(bookingsOnline) == 0 {
		if event.Online && !event.Offline {
			return c.Send(message.BookingsNotFoundMessage)
		} else {
			return c.Send(message.BookingsOnlineNotFoundMessage)
		}
	}

	for _, booking := range bookingsOnline {
		user, exists, err := h.userRepository.Get(booking.UserID)
		if err != nil {
			return fmt.Errorf("Failed get user: %s", err)
		}
		if !exists {
			return fmt.Errorf("Failed get user, user not exists")
		}
		content := message.BookingMessageContent(booking, user)
		err = c.Send(content[0], content[1:]...)
		if err != nil {
			return err
		}
	}
	return nil
}
