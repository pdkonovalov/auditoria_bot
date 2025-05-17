package handler

import (
	"fmt"

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

	var (
		getBookingsOffline bool
		getBookingsOnline  bool
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
		getBookingsOffline = true
	}
	if event.Online || bookingsOnlineExists {
		getBookingsOnline = true
	}

	content := message.GetBookingsMessageContent(
		event,
		createdBy,
		updatedBy,
		eventURL,
		getBookingsOffline,
		getBookingsOnline,
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

	bookings, err := h.bookingRepository.GetByEventID(eventID)
	if err != nil {
		return fmt.Errorf("Failed get bookings: %s", err)
	}

	offlineBookings := make([]*entity.Booking, 0)
	for _, booking := range bookings {
		if booking.Offline {
			offlineBookings = append(offlineBookings, booking)
		}
	}

	if len(offlineBookings) == 0 {
		if event.Offline && !event.Online {
			return c.Send(message.BookingsNotFoundMessage)
		} else {
			return c.Send(message.BookingsOfflineNotFoundMessage)
		}
	}

	for _, booking := range offlineBookings {
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

	bookings, err := h.bookingRepository.GetByEventID(eventID)
	if err != nil {
		return fmt.Errorf("Failed get bookings: %s", err)
	}

	onlineBookings := make([]*entity.Booking, 0)
	for _, booking := range bookings {
		if booking.Online {
			onlineBookings = append(onlineBookings, booking)
		}
	}

	if len(onlineBookings) == 0 {
		if event.Online && !event.Offline {
			return c.Send(message.BookingsNotFoundMessage)
		} else {
			return c.Send(message.BookingsOnlineNotFoundMessage)
		}
	}

	for _, booking := range onlineBookings {
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
