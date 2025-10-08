package handler

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/config"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/repository"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/callback"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/command"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/state"
	user_message "github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/message"

	tele "gopkg.in/telebot.v4"
)

type AdminHandler struct {
	botUsername string

	userRepository    repository.UserRepository
	eventRepository   repository.EventRepository
	bookingRepository repository.BookingRepository

	location *time.Location

	defaultPaymentDetails config.PaymentDetailsList

	bookingsPerPage int
}

func NewAdminHandler(
	cfg *config.Config,
	userRepository repository.UserRepository,
	eventRepository repository.EventRepository,
	bookingRepository repository.BookingRepository,
) (*AdminHandler, error) {
	location, err := time.LoadLocation(cfg.TelegramBotTimezone)
	if err != nil {
		return nil, err
	}
	return &AdminHandler{
		botUsername:           cfg.TelegramBotUsername,
		userRepository:        userRepository,
		eventRepository:       eventRepository,
		bookingRepository:     bookingRepository,
		location:              location,
		defaultPaymentDetails: cfg.TelegramBotDefaultPaymentDetails,
		bookingsPerPage:       cfg.TelegramBotBookingsPerPage,
	}, nil
}

func (h *AdminHandler) OnCallback(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}
	user.State = state.Init
	if eventID, ok := user.Context["eventID"].(string); ok {
		event, exists, err := h.eventRepository.Get(eventID)
		if err != nil {
			return fmt.Errorf("Failed get event: %s", err)
		}
		if exists && event.Draft {
			exists, err := h.eventRepository.Delete(eventID)
			if err != nil {
				return fmt.Errorf("Failed delete draft event: %s", err)
			}
			if !exists {
				return fmt.Errorf("Failed delete draft event, event not exists")
			}
		}
	}
	user.Context = make(map[string]any)
	exists, err := h.userRepository.Update(&user)
	if err != nil {
		return fmt.Errorf("Failed update user: %s", err)
	}
	if !exists {
		return fmt.Errorf("Failed update user, user not exists")
	}

	unique, ok := c.Get("callback_unique").(string)
	if !ok {
		return fmt.Errorf("Failed get callback unique from context")
	}
	if len(unique) == 0 {
		return fmt.Errorf("Failed get callback unique, callback unique is empty")
	}
	data, ok := c.Get("callback_data").(string)
	if !ok {
		return fmt.Errorf("Failed get callback data from context")
	}
	if len(data) == 0 {
		return fmt.Errorf("Failed get callback data, callback data is empty")
	}

	data_decoded, err := callback.Decode(data)
	if err != nil {
		return fmt.Errorf("Failed decode callback data: %s", err)
	}

	switch unique {
	// event button
	case callback.Event:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.Event(c)

	// event photo text button
	case callback.EventPhotoText:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.EventPhotoText(c)

	// date button
	case callback.EventsByDate:
		date, ok := data_decoded["date"]
		if !ok {
			return fmt.Errorf("Failed get date from callback data")
		}
		c.Set("date", date)
		return h.EventsByDate(c)

	// show bookings format selection button
	case callback.ShowBookingsFormatSelection:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.ShowBookingsFormatSelection(c)

	// show bookings button
	case callback.ShowBookings:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		format, ok := data_decoded["format"]
		if !ok {
			return fmt.Errorf("Failed get format from callback data")
		}
		c.Set("format", format)
		page, ok := data_decoded["page"]
		if !ok {
			return fmt.Errorf("Failed get page from callback data")
		}
		c.Set("page", page)
		return h.ShowBookings(c)

	// booking button
	case callback.Booking:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		userID, ok := data_decoded["userID"]
		if !ok {
			return fmt.Errorf("Failed get user id from callback data")
		}
		c.Set("userID", userID)
		return h.Booking(c)

	// booking check in button
	case callback.BookingCheckIn:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		userID, ok := data_decoded["userID"]
		if !ok {
			return fmt.Errorf("Failed get user id from callback data")
		}
		c.Set("userID", userID)
		return h.BookingCheckIn(c)

	// edit event button
	case callback.EditEvent:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.EditEvent(c)

	// edit format button
	case callback.EditFormat:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.EditEventFormatInit(c)

	// edit event offline paid button
	case callback.EditOfflinePaid:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		c.Set("format", "offline")
		return h.EditEventPaidInit(c)

	// edit event online paid button
	case callback.EditOnlinePaid:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		c.Set("format", "online")
		return h.EditEventPaidInit(c)

	// edit event title button
	case callback.EditTitle:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.EditEventTitleInit(c)

	// edit event time button
	case callback.EditTime:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.EditEventTimeInit(c)

	// edit event payment details button
	case callback.EditPaymentDetails:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.EditEventPaymentDetailsInit(c)

	// edit event photo and text button
	case callback.EditPhotoText:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.EditEventPhotoTextInit(c)

	// send notification button
	case callback.SendNotificationFormatSelection:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.SendNotificationFormatSelection(c)

	// send notification button
	case callback.SendNotification:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		format, ok := data_decoded["format"]
		if !ok {
			return fmt.Errorf("Failed get format from callback data")
		}
		c.Set("format", format)
		return h.SendNotificationInit(c)

	// delete event button
	case callback.DeleteEvent:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		return h.DeleteEventInit(c)
	}

	return fmt.Errorf("Unexpected callback unique: %s", unique)
}

func (h *AdminHandler) OnMessage(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	if len(c.Message().Text) != 0 && c.Message().Text[0] == '/' {
		user.State = state.Init
		if eventID, ok := user.Context["eventID"].(string); ok {
			event, exists, err := h.eventRepository.Get(eventID)
			if err != nil {
				return fmt.Errorf("Failed get event: %s", err)
			}
			if exists && event.Draft {
				exists, err := h.eventRepository.Delete(eventID)
				if err != nil {
					return fmt.Errorf("Failed delete draft event: %s", err)
				}
				if !exists {
					return fmt.Errorf("Failed delete draft event, event not exists")
				}
			}
		}
		user.Context = make(map[string]any)
		exists, err := h.userRepository.Update(&user)
		if err != nil {
			return fmt.Errorf("Failed update user: %s", err)
		}
		if !exists {
			return fmt.Errorf("Failed update user, user not exists")
		}
		switch {
		case len(c.Message().Text) >= len("/start") && c.Message().Text[:len("/start")] == "/start":
			if len(c.Message().Payload) == 0 {
				return h.Hello(c)
			}
			c.Set("eventID", c.Message().Payload)
			return h.Event(c)
		case c.Message().Text == command.Events.Text:
			return h.Events(c)
		case c.Message().Text == command.NewEvent.Text:
			return h.NewEventInit(c)
		case c.Message().Text == command.User.Text:
			user.Admin = false
			exists, err := h.userRepository.Update(&user)
			if err != nil {
				return fmt.Errorf("Failed update user: %s", err)
			}
			if !exists {
				return fmt.Errorf("Failed update user, user not exists")
			}
			return c.Send(user_message.HelloAdminMessage, user_message.HelloEntities)
		case c.Message().Text == command.Cancel.Text:
			return h.Hello(c)
		}
		return h.Hello(c)
	}

	switch user.State {
	case state.Init:
		return h.Hello(c)

	// new event
	case state.NewEventWaitInputFormat:
		return h.NewEventFormatInput(c)
	case state.NewEventWaitInputPaid:
		return h.NewEventPaidInput(c)
	case state.NewEventWaitInputTitle:
		return h.NewEventTitleInput(c)
	case state.NewEventWaitInputTime:
		return h.NewEventTimeInput(c)
	case state.NewEventWaitInputPaymentDetails:
		return h.NewEventPaymentDetailsInput(c)
	case state.NewEventWaitInputPhotoText:
		return h.NewEventPhotoTextInput(c)

	// edit event
	case state.EditEventWaitInputFormat:
		return h.EditEventFormatInput(c)
	case state.EditEventWaitInputPaid:
		return h.EditEventPaidInput(c)
	case state.EditEventWaitInputTitle:
		return h.EditEventTitleInput(c)
	case state.EditEventWaitInputTime:
		return h.EditEventTimeInput(c)
	case state.EditEventWaitInputPaymentDetails:
		return h.EditEventPaymentDetailsInput(c)
	case state.EditEventWaitInputPhotoText:
		return h.EditEventPhotoTextInput(c)

	// send notification
	case state.SendNotificationWaitInputPhotoText:
		return h.SendNotificationPhotoTextInput(c)

	// delete event
	case state.DeleteEventWaitInputAreYouSure:
		return h.DeleteEventAreYouSureInput(c)
	}

	return fmt.Errorf("Unexpected user state: %s", user.State)
}
