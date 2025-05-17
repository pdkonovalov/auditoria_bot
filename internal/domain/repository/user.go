package repository

import "github.com/pdkonovalov/auditoria_bot/internal/domain/entity"

type UserRepository interface {
	Get(userID int64) (*entity.User, bool, error)
	Create(*entity.User) (bool, error)
	Update(*entity.User) (bool, error)
}
