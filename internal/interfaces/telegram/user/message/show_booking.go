package message

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/callback"

	tele "gopkg.in/telebot.v4"
)

func ShowBookingMessageContent(event *entity.Event, booking *entity.Booking, filter string) []any {
	parts := make([]string, 0)
	entities := make(tele.Entities, 0)
	curLen := 0

	if event.Offline && event.Online ||
		event.Offline && booking.Online ||
		event.Online && booking.Offline {
		format := "Формат:"
		if booking.Offline {
			parts = append(parts,
				fmt.Sprintf("%s %s\n", format, "офлайн"),
			)
		} else {
			parts = append(parts,
				fmt.Sprintf("%s %s\n", format, "онлайн"),
			)
		}
		entities = append(entities,
			tele.MessageEntity{
				Type:   tele.EntityBold,
				Offset: curLen,
				Length: len(utf16.Encode([]rune(format))),
			},
		)
		curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))
	}

	if booking.Offline && event.OfflinePaid ||
		booking.Online && event.OnlinePaid ||
		booking.Payment != nil {
		payment := "Оплата:"
		if booking.Payment != nil {
			parts = append(parts,
				fmt.Sprintf("%s %s\n", payment, "скриншот"),
			)
		} else {
			parts = append(parts,
				fmt.Sprintf("%s %s\n", payment, "-"),
			)
		}
		entities = append(entities,
			tele.MessageEntity{
				Type:   tele.EntityBold,
				Offset: curLen,
				Length: len(utf16.Encode([]rune(payment))),
			},
		)
		curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))
	}

	additional := "Дополнительная информация:"
	if len(booking.Text) != 0 {
		parts = append(parts,
			fmt.Sprintf("%s %s\n", additional, booking.Text),
		)
	} else {
		parts = append(parts,
			fmt.Sprintf("%s %s\n", additional, "-"),
		)
	}
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(additional))),
		},
	)

	text := strings.Join(parts, "")

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
