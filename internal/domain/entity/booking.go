package entity

import (
	"time"

	tele "gopkg.in/telebot.v4"
)

type Booking struct {
	EventID         string
	UserID          int64
	Payment         bool
	PaymentPhoto    *tele.Photo
	PaymentDocument *tele.Document
	Text            string
	Offline         bool
	Online          bool
	CheckIn         bool
	CheckInAt       *time.Time
	CheckInBy       *int64
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	Draft           bool
}
