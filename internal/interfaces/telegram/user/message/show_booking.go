package message

import (
	"fmt"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/callback"

	tele "gopkg.in/telebot.v4"
)

func ShowBookingMessageContent(booking *entity.Booking, filter string) []any {
	entities := tele.Entities{
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: 0,
			Length: len(utf16.Encode([]rune("Дополнительная информация:"))),
		},
	}
	text := fmt.Sprintf("Дополнительная информация: %s\n", booking.Text)

	if booking.Payment == nil {
		return []any{text, entities, showBookingInlineKeyboard(booking, filter)}
	}
	photo := booking.Payment
	photo.Caption = text
	return []any{photo, entities, showBookingInlineKeyboard(booking, filter)}
}

func showBookingInlineKeyboard(booking *entity.Booking, filter string) *tele.ReplyMarkup {
	keyboard := [][]tele.InlineButton{
		{
			{
				Text:   "< Назад",
				Unique: callback.Event,
				Data: callback.Encode(map[string]string{
					"eventID": booking.EventID,
					"filter":  filter,
				}),
			},
		},
	}
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}
