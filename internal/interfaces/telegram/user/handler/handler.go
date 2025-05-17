package handler

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/config"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/repository"
	"github.com/pdkonovalov/auditoria_bot/internal/infrastructure/validator"
	admin_message "github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/message"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/callback"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/command"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/state"

	tele "gopkg.in/telebot.v4"
)

type UserHandler struct {
	userRepository    repository.UserRepository
	eventRepository   repository.EventRepository
	bookingRepository repository.BookingRepository
	validator         *validator.Validator
	location          *time.Location
}

func NewUserHandler(
	cfg *config.Config,
	userRepository repository.UserRepository,
	eventRepository repository.EventRepository,
	bookingRepository repository.BookingRepository,
	validator *validator.Validator,
) (*UserHandler, error) {
	location, err := time.LoadLocation(cfg.TelegramBotTimezone)
	if err != nil {
		return nil, err
	}
	return &UserHandler{
		userRepository:    userRepository,
		eventRepository:   eventRepository,
		bookingRepository: bookingRepository,
		validator:         validator,
		location:          location,
	}, nil
}

func (h *UserHandler) OnCallback(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}
	user.State = state.Init
	if eventID, ok := user.Context["eventID"].(string); ok {
		booking, exists, err := h.bookingRepository.Get(user.UserID, eventID)
		if err != nil {
			return fmt.Errorf("Failed get booking: %s", err)
		}
		if exists && booking.Draft {
			exists, err := h.bookingRepository.Delete(user.UserID, eventID)
			if err != nil {
				return fmt.Errorf("Failed delete draft booking: %s", err)
			}
			if !exists {
				return fmt.Errorf("Failed delete draft booking, booking not exists")
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
	data, _ := c.Get("callback_data").(string)

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
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.Event(c)

	// date button
	case callback.EventsByDate:
		date, ok := data_decoded["date"]
		if !ok {
			return fmt.Errorf("Failed get date from callback data")
		}
		c.Set("date", date)
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.EventsByDate(c)

	// booking button
	case callback.Booking:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.BookingInit(c)

	// edit booking button
	case callback.EditBooking:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.EditBooking(c)

	// edit format button
	case callback.EditFormat:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.EditBookingFormatInit(c)

	// edit payment button
	case callback.EditPayment:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.EditBookingPaymentInit(c)

	// edit additional info button
	case callback.EditAdditionalInfo:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.EditBookingAdditionalInfoInit(c)

	// delete booking button
	case callback.DeleteBooking:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.DeleteBookingInit(c)

	// show booking button
	case callback.ShowBooking:
		eventID, ok := data_decoded["eventID"]
		if !ok {
			return fmt.Errorf("Failed get event id from callback data")
		}
		c.Set("eventID", eventID)
		filter, ok := data_decoded["filter"]
		if !ok {
			return fmt.Errorf("Failed get filter from callback data")
		}
		c.Set("filter", filter)
		return h.ShowBooking(c)

	// edit contact info button
	case callback.EditContactInfo:
		return h.EditContactInfoInit(c)
	}

	return fmt.Errorf("Unexpected callback unique: %s", unique)
}

func (h *UserHandler) OnTextPhoto(c tele.Context) error {
	user, ok := c.Get("user").(entity.User)
	if !ok {
		return fmt.Errorf("Failed get user from context")
	}

	if len(c.Message().Text) != 0 && c.Message().Text[0] == '/' {
		user.State = state.Init
		if eventID, ok := user.Context["eventID"].(string); ok {
			booking, exists, err := h.bookingRepository.Get(user.UserID, eventID)
			if err != nil {
				return fmt.Errorf("Failed get booking: %s", err)
			}
			if exists && booking.Draft {
				exists, err := h.bookingRepository.Delete(user.UserID, eventID)
				if err != nil {
					return fmt.Errorf("Failed delete draft booking: %s", err)
				}
				if !exists {
					return fmt.Errorf("Failed delete draft booking, booking not exists")
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
			c.Set("filter", "all")
			return h.Event(c)
		case c.Message().Text == command.Events.Text:
			return h.Events(c)
		case c.Message().Text == command.MyEvents.Text:
			return h.MyEvents(c)
		case c.Message().Text == command.SetContactInfo.Text:
			return h.SetContactInfo(c)
		case c.Message().Text == command.Cancel.Text:
			return nil
		case c.Message().Text == command.Admin.Text:
			is_admin, ok := c.Get("is_admin").(bool)
			if !ok {
				return fmt.Errorf("Failed get is admin from context")
			}
			if is_admin {
				user.Admin = true
				exists, err := h.userRepository.Update(&user)
				if err != nil {
					return fmt.Errorf("Failed update user: %s", err)
				}
				if !exists {
					return fmt.Errorf("Failed update user, user not exists")
				}
				return c.Send(admin_message.HelloMessage, admin_message.HelloEntities)
			}
		}
		return h.Hello(c)
	}

	switch user.State {
	case state.Init:
		return h.Hello(c)

	// edit contact info
	case state.EditContactInfoWaitInput:
		return h.EditContactInfoInput(c)

	// booking
	case state.BookingWaitInputContactInfo:
		return h.BookingContactInfoInput(c)
	case state.BookingWaitInputFormat:
		return h.BookingFormatInput(c)
	case state.BookingWaitInputPayment:
		return h.BookingPaymentInput(c)
	case state.BookingWaitInputAdditionalInfo:
		return h.BookingAdditionalInfoInput(c)

	// edit booking
	case state.EditBookingWaitInputFormat:
		return h.EditBookingFormatInput(c)
	case state.EditBookingWaitInputPayment:
		return h.EditBookingPaymentInput(c)
	case state.EditBookingWaitInputAdditionalInfo:
		return h.EditBookingAdditionalInfoInput(c)

	// delete booking
	case state.DeleteBookingWaitInputAreYouSure:
		return h.DeleteBookingAreYouSureInput(c)
	}

	return fmt.Errorf("Unexpected user state: %s", user.State)
}
