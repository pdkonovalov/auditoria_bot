package handler

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) Event(c tele.Context) error {
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
		getBookingsOffline      bool
		getBookingsOnline       bool
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
		getBookingsOffline = true
		sendNotificationOffline = true
	}
	if event.Online || bookingsOnlineExists {
		getBookingsOnline = true
		sendNotificationOnline = true
	}

	content := message.EventMessageContent(
		event,
		createdBy,
		updatedBy,
		eventURL,
		getBookingsOffline,
		getBookingsOnline,
		sendNotificationOffline,
		sendNotificationOnline)
	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}
