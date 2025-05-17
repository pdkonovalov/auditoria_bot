package repository

import (
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
)

type EventRepository interface {
	Get(eventID string) (*entity.Event, bool, error)
	GetFirstAfter(time.Time, *int64) (*entity.Event, bool, error)
	GetFirstBefore(time.Time, *int64) (*entity.Event, bool, error)
	GetByDate(time.Time, *int64) ([]*entity.Event, error)
	Create(*entity.Event) (eventID string, err error)
	Update(*entity.Event) (bool, error)
	Delete(eventID string) (bool, error)
}
