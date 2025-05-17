package message

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/callback"
	"github.com/pdkonovalov/russian_time"

	tele "gopkg.in/telebot.v4"
)

func EventMessageContent(event *entity.Event, isBooked bool, filter string) []any {
	if event.Photo == nil {
		return []any{event.Text, event.Entities, eventInlineKeyboard(event, isBooked, filter)}
	}
	photo := event.Photo
	photo.Caption = event.Text
	return []any{photo, event.Entities, eventInlineKeyboard(event, isBooked, filter)}
}

func eventInlineKeyboard(event *entity.Event, isBooked bool, filter string) *tele.ReplyMarkup {
	keyboard := make([][]tele.InlineButton, 0)
	if event.Time.After(time.Now()) {
		if !isBooked {
			keyboard = append(keyboard, []tele.InlineButton{
				{
					Text:   "Записаться",
					Unique: callback.Booking,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
						"filter":  filter,
					}),
				},
			})
		} else {
			keyboard = append(keyboard,
				[]tele.InlineButton{
					{
						Text:   "Изменить данные записи",
						Unique: callback.EditBooking,
						Data: callback.Encode(map[string]string{
							"eventID": event.EventID,
							"filter":  filter,
						}),
					},
				},
				[]tele.InlineButton{
					{
						Text:   "Отменить запись",
						Unique: callback.DeleteBooking,
						Data: callback.Encode(map[string]string{
							"eventID": event.EventID,
							"filter":  filter,
						}),
					},
				},
			)
		}
	} else if isBooked {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Посмотреть данные записи",
					Unique: callback.ShowBooking,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
						"filter":  filter,
					}),
				},
			},
		)
	}
	var eventsTitle string
	if filter == "my" {
		eventsTitle = "Мои мероприятия"
	} else {
		eventsTitle = "Мероприятия"
	}
	keyboard = append(keyboard,
		[]tele.InlineButton{
			{
				Text:   fmt.Sprintf("< %s %s", eventsTitle, russian_time.DayMonth(&event.Time, russian_time.RCase2)),
				Unique: callback.EventsByDate,
				Data: callback.Encode(map[string]string{
					"date":   event.Time.Format("02.01.2006"),
					"filter": filter,
				}),
			},
		},
	)
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}
