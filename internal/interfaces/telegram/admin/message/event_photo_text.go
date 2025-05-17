package message

import (
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/callback"

	tele "gopkg.in/telebot.v4"
)

func EventPhotoTextMessageContent(event *entity.Event) []any {
	if event.Photo == nil {
		return []any{event.Text, event.Entities, eventPhotoTextInlineKeyboard(event)}
	}
	photo := event.Photo
	photo.Caption = event.Text
	return []any{photo, event.Entities, eventPhotoTextInlineKeyboard(event)}
}

func eventPhotoTextInlineKeyboard(event *entity.Event) *tele.ReplyMarkup {
	keyboard := [][]tele.InlineButton{
		{
			{
				Text:   "< Назад",
				Unique: callback.Event,
				Data:   event.EventID,
			},
		},
	}
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}
