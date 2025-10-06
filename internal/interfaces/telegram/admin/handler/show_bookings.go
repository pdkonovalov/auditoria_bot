package handler

import (
	"fmt"
	"strconv"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"

	tele "gopkg.in/telebot.v4"
)

func (h *AdminHandler) ShowBookingsFormatSelection(c tele.Context) error {
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

	content := message.ShowBookingsFormatSelectionMessageContent(
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

func (h *AdminHandler) ShowBookings(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}

	format, ok := c.Get("format").(string)
	if !ok {
		return fmt.Errorf("Failed get format from context")
	}

	pageStr, ok := c.Get("page").(string)
	if !ok {
		return fmt.Errorf("Failed get page from context")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return fmt.Errorf("Invalid page value in context. Failed convert to int: %s", err)
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

	bookingsOnline, err := h.bookingRepository.GetByEventID(eventID, false, true)
	if err != nil {
		return fmt.Errorf("Failed get online bookings: %s", err)
	}

	var bookingsTarget []*entity.Booking

	if format == "offline" {
		bookingsTarget = bookingsOffline
	} else if format == "online" {
		bookingsTarget = bookingsOnline
	} else {
		return fmt.Errorf("Invalid format value in context: %s", format)
	}

	var pageCount int

	if len(bookingsTarget) == 0 {
		pageCount = 1
	} else if len(bookingsTarget)%h.bookingsPerPage == 0 {
		pageCount = len(bookingsTarget) / h.bookingsPerPage
	} else {
		pageCount = len(bookingsTarget)/h.bookingsPerPage + 1
	}

	if page < 0 {
		return fmt.Errorf("Invalid page value in context: %s", format)
	}

	if page >= pageCount {
		page = pageCount - 1
	}

	bookingsPage := make([]*entity.Booking, 0)
	for i := page * h.bookingsPerPage; i < len(bookingsTarget) && i < (page+1)*h.bookingsPerPage; i++ {
		bookingsPage = append(bookingsPage, bookingsTarget[i])
	}

	bookingsPageUsers := make([]*entity.User, 0)
	for _, booking := range bookingsPage {
		user, exists, err := h.userRepository.Get(booking.UserID)
		if err != nil {
			return fmt.Errorf("Failed get user for booking: %s", err)
		}
		if !exists {
			return fmt.Errorf("Failed get user for booking, user not exists")
		}
		bookingsPageUsers = append(bookingsPageUsers, user)
	}

	content := message.ShowBookingsMessageContent(
		event,
		page,
		bookingsPageUsers,
		page != 0,
		page != pageCount-1,
		len(bookingsOffline) != 0,
		len(bookingsOnline) != 0,
		format,
	)

	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}

func (h *AdminHandler) Booking(c tele.Context) error {
	eventID, ok := c.Get("eventID").(string)
	if !ok {
		return fmt.Errorf("Failed get event id from context")
	}

	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		return fmt.Errorf("Failed get user id from context")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid user id value in context. Failed convert to int64: %s", err)
	}

	bookingTarget, exists, err := h.bookingRepository.Get(userID, eventID)
	if err != nil {
		return fmt.Errorf("Failed get booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get booking, booking not exists")
	}

	bookings, err := h.bookingRepository.GetByEventID(eventID, bookingTarget.Offline, bookingTarget.Online)
	if err != nil {
		return fmt.Errorf("Failed get bookings: %s", err)
	}

	var page int

	for i, booking := range bookings {
		if booking.UserID == userID {
			page = i / h.bookingsPerPage
			break
		}
	}

	bookingTargetUser, exists, err := h.userRepository.Get(bookingTarget.UserID)
	if err != nil {
		return fmt.Errorf("Failed get user for booking: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed get user for booking, user not exists")
	}

	content := message.BookingMessageContent(bookingTarget, bookingTargetUser, page)

	err = c.EditOrSend(content[0], content[1:]...)
	if err != nil {
		return c.Send(content[0], content[1:]...)
	}
	return nil
}
