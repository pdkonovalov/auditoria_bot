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
	botUsername       string
	userRepository    repository.UserRepository
	eventRepository   repository.EventRepository
	bookingRepository repository.BookingRepository
	location          *time.Location
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
		botUsername:       cfg.TelegramBotUsername,
		userRepository:    userRepository,
		eventRepository:   eventRepository,
		bookingRepository: bookingRepository,
		location:          location,
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

	switch unique {
	// event button
	case callback.Event:
		c.Set("eventID", data)
		return h.Event(c)

	// event photo text button
	case callback.EventPhotoText:
		c.Set("eventID", data)
		return h.EventPhotoText(c)

	// date button
	case callback.EventsByDate:
		c.Set("date", data)
		return h.EventsByDate(c)

	// get bookings button
	case callback.GetBookings:
		c.Set("eventID", data)
		return h.GetBookings(c)

	// get bookings offline button
	case callback.GetBookingsOffline:
		c.Set("eventID", data)
		return h.GetBookingsOffline(c)

	// get bookings online button
	case callback.GetBookingsOnline:
		c.Set("eventID", data)
		return h.GetBookingsOnline(c)

	// edit event button
	case callback.EditEvent:
		c.Set("eventID", data)
		return h.EditEvent(c)

	// edit format button
	case callback.EditFormat:
		c.Set("eventID", data)
		return h.EditEventFormatInit(c)

	// edit event offline paid button
	case callback.EditOfflinePaid:
		c.Set("eventID", data)
		c.Set("format", "offline")
		return h.EditEventPaidInit(c)

	// edit event online paid button
	case callback.EditOnlinePaid:
		c.Set("eventID", data)
		c.Set("format", "online")
		return h.EditEventPaidInit(c)

	// edit event title button
	case callback.EditTitle:
		c.Set("eventID", data)
		return h.EditEventTitleInit(c)

	// edit event time button
	case callback.EditTime:
		c.Set("eventID", data)
		return h.EditEventTimeInit(c)

	// edit event photo and text button
	case callback.EditPhotoText:
		c.Set("eventID", data)
		return h.EditEventPhotoTextInit(c)

	// send notification button
	case callback.SendNotification:
		c.Set("eventID", data)
		return h.SendNotification(c)

	// send notification offline button
	case callback.SendNotificationOffline:
		c.Set("eventID", data)
		c.Set("format", "offline")
		return h.SendNotificationInit(c)

		// send notification online button
	case callback.SendNotificationOnline:
		c.Set("eventID", data)
		c.Set("format", "online")
		return h.SendNotificationInit(c)

	// delete event button
	case callback.DeleteEvent:
		c.Set("eventID", data)
		return h.DeleteEventInit(c)
	}

	return fmt.Errorf("Unexpected callback unique: %s", unique)
}

func (h *AdminHandler) OnTextPhoto(c tele.Context) error {
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
			return nil
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
