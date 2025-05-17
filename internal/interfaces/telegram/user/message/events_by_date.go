package message

import (
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/callback"
	"github.com/pdkonovalov/russian_time"

	tele "gopkg.in/telebot.v4"
)

var EventsNotFoundMessage = "Мероприятия не найдены :("

func EventsByDateMessageContent(
	prevDate *string,
	curDate string,
	nextDate *string,
	events []*entity.Event,
	filter string,
) ([]any, error) {
	time, err := time.Parse("02.01.2006", curDate)
	if err != nil {
		return nil, err
	}
	var eventsTitle string
	if filter == "my" {
		eventsTitle = "Мои мероприятия"
	} else {
		eventsTitle = "Мероприятия"
	}
	text := fmt.Sprintf("%s %s", eventsTitle, russian_time.DayMonth(&time, russian_time.RCase2))
	keyboard, err := eventsByDateInlineKeyboard(prevDate, nextDate, events, filter)
	if err != nil {
		return nil, err
	}
	return []any{text, keyboard}, nil
}

func eventsByDateInlineKeyboard(
	prevDate *string,
	nextDate *string,
	events []*entity.Event,
	filter string,
) (*tele.ReplyMarkup, error) {
	keyboard := make([][]tele.InlineButton, len(events))
	for index, event := range events {
		keyboard[index] = []tele.InlineButton{
			{
				Text:   fmt.Sprintf("%s %s", event.Time.Format("15:04"), event.Title),
				Unique: callback.Event,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
					"filter":  filter,
				}),
			},
		}
	}
	navigationRow := make([]tele.InlineButton, 0)
	if prevDate != nil {
		time, err := time.Parse("02.01.2006", *prevDate)
		if err != nil {
			return nil, err
		}
		navigationRow = append(navigationRow,
			tele.InlineButton{
				Text:   fmt.Sprintf("< %s", russian_time.DayMonth(&time, russian_time.RCase2)),
				Unique: callback.EventsByDate,
				Data: callback.Encode(map[string]string{
					"date":   *prevDate,
					"filter": filter,
				}),
			},
		)
	}
	if nextDate != nil {
		time, err := time.Parse("02.01.2006", *nextDate)
		if err != nil {
			return nil, err
		}
		navigationRow = append(navigationRow,
			tele.InlineButton{
				Text:   fmt.Sprintf("%s >", russian_time.DayMonth(&time, russian_time.RCase2)),
				Unique: callback.EventsByDate,
				Data: callback.Encode(map[string]string{
					"date":   *nextDate,
					"filter": filter,
				}),
			},
		)
	}
	keyboard = append(keyboard, navigationRow)
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}, nil
}
