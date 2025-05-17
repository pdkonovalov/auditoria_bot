package entity

import (
	"time"

	tele "gopkg.in/telebot.v4"
)

type Event struct {
	EventID     string
	Title       string
	Photo       *tele.Photo
	Text        string
	Entities    tele.Entities
	Time        time.Time
	Offline     bool
	OfflinePaid bool
	Online      bool
	OnlinePaid  bool
	CreatedAt   time.Time
	CreatedBy   int64
	UpdatedAt   *time.Time
	UpdatedBy   *int64
	Draft       bool
}
