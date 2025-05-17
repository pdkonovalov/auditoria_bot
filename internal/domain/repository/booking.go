package repository

import "github.com/pdkonovalov/auditoria_bot/internal/domain/entity"

type BookingRepository interface {
	Get(userID int64, eventID string) (*entity.Booking, bool, error)
	GetByUserID(userID int64) ([]*entity.Booking, error)
	GetByEventID(eventID string) ([]*entity.Booking, error)
	Create(*entity.Booking) (bool, error)
	Update(*entity.Booking) (bool, error)
	Delete(userID int64, eventID string) (bool, error)
}
