package middleware

import (
	"log"

	"github.com/pdkonovalov/auditoria_bot/internal/config"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/repository"

	tele "gopkg.in/telebot.v4"
)

type AuthMiddleware struct {
	admins         map[int64]config.Admin
	userRepository repository.UserRepository
}

func NewAuthMiddleware(
	cfg *config.Config,
	userRepository repository.UserRepository,
) (*AuthMiddleware, error) {
	admins := make(map[int64]config.Admin, len(cfg.TelegramBotAdminList))
	for _, admin := range cfg.TelegramBotAdminList {
		admins[admin.UserID] = admin
	}
	return &AuthMiddleware{
		admins:         admins,
		userRepository: userRepository,
	}, nil
}

func (m *AuthMiddleware) Auth(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		userID := c.Sender().ID
		user, isExist, err := m.userRepository.Get(userID)
		if err != nil {
			log.Printf("Failed get user: %s", err)
			return err
		}
		admin, isAdmin := m.admins[userID]
		if !isExist {
			user = &entity.User{
				UserID:      userID,
				Username:    c.Sender().Username,
				ContactInfo: "",
				State:       entity.UserStateInit,
				Context:     make(map[string]any),
				Admin:       isAdmin,
			}
			if isAdmin {
				user.FirstName = admin.FirstName
				user.LastName = admin.LastName
			} else {
				user.FirstName = c.Sender().FirstName
				user.LastName = c.Sender().LastName
			}
			_, err := m.userRepository.Create(user)
			if err != nil {
				log.Printf("Failed create user: %s", err)
				return err
			}
			c.Set("user", *user)
			c.Set("is_admin", isAdmin)
			return next(c)
		}
		changed := false
		if user.Username != c.Sender().Username {
			user.Username = c.Sender().Username
			changed = true
		}
		if isAdmin {
			if user.FirstName != admin.FirstName {
				user.FirstName = admin.FirstName
				changed = true
			}
			if user.LastName != admin.LastName {
				user.LastName = admin.LastName
				changed = true
			}
		} else {
			if user.FirstName != c.Sender().FirstName {
				user.FirstName = c.Sender().FirstName
				changed = true
			}
			if user.LastName != c.Sender().LastName {
				user.LastName = c.Sender().LastName
				changed = true
			}
		}
		if changed {
			isExist, err := m.userRepository.Update(user)
			if err != nil {
				log.Printf("Failed update user: %s", err)
				return err
			}
			if !isExist {
				log.Print("Failed update user, user not exist")
				return err
			}
		}
		c.Set("user", *user)
		c.Set("is_admin", isAdmin)
		return next(c)
	}
}
