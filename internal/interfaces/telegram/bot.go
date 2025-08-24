package telegram

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/config"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/repository"
	"github.com/pdkonovalov/auditoria_bot/internal/infrastructure/validator"
	admin "github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/handler"
	custom_middleware "github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/middleware"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/command"
	user "github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/handler"

	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func New(
	cfg *config.Config,
	userRepository repository.UserRepository,
	eventRepository repository.EventRepository,
	bookingRepository repository.BookingRepository,
	validator *validator.Validator,
) (*tele.Bot, error) {
	b, err := tele.NewBot(tele.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	err = b.SetCommands(command.List)
	if err != nil {
		return nil, err
	}

	err = setupMiddleware(cfg, b, userRepository)
	if err != nil {
		return nil, err
	}

	err = setupRouter(cfg, b, userRepository, eventRepository, bookingRepository, validator)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func setupMiddleware(
	cfg *config.Config, b *tele.Bot,
	userRepository repository.UserRepository,
) error {
	b.Use(middleware.Recover())

	if cfg.LogLevel == config.LogLevelDebug {
		b.Use(middleware.Logger())
	}

	b.Use(middleware.AutoRespond())

	auth_middleware, err := custom_middleware.NewAuthMiddleware(cfg, userRepository)
	if err != nil {
		return fmt.Errorf("Failed create auth middleware: %s", err)
	}
	b.Use(auth_middleware.Auth)

	b.Use(custom_middleware.CallbackData)

	return nil
}

func setupRouter(
	cfg *config.Config,
	b *tele.Bot,
	userRepository repository.UserRepository,
	eventRepository repository.EventRepository,
	bookingRepository repository.BookingRepository,
	validator *validator.Validator,
) error {
	adminHandler, err := admin.NewAdminHandler(cfg, userRepository, eventRepository, bookingRepository)
	if err != nil {
		return fmt.Errorf("Failed create admin handler: %s", err)
	}

	userHandler, err := user.NewUserHandler(cfg, userRepository, eventRepository, bookingRepository, validator)
	if err != nil {
		return fmt.Errorf("Failed create user handler: %s", err)
	}

	b.Handle(tele.OnText, route(adminHandler.OnMessage, userHandler.OnMessage))
	b.Handle(tele.OnPhoto, route(adminHandler.OnMessage, userHandler.OnMessage))
	b.Handle(tele.OnContact, route(adminHandler.OnMessage, userHandler.OnMessage))
	b.Handle(tele.OnDocument, route(adminHandler.OnMessage, userHandler.OnMessage))

	b.Handle(tele.OnCallback, route(adminHandler.OnCallback, userHandler.OnCallback))

	return nil
}

func route(adminHandler, userHandler tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user, ok := c.Get("user").(entity.User)
		if !ok {
			return fmt.Errorf("Failed get user from context")
		}
		if user.Admin {
			return adminHandler(c)
		}
		return userHandler(c)
	}
}
