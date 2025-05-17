package entity

import (
	"time"

	tele "gopkg.in/telebot.v4"
)

type Booking struct {
	EventID   string
	UserID    int64
	Payment   *tele.Photo
	Text      string
	Offline   bool
	Online    bool
	CreatedAt time.Time
	UpdatedAt *time.Time
	Draft     bool
}
