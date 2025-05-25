package message

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/callback"

	tele "gopkg.in/telebot.v4"
)

const (
	BookingsNotFoundMessage        = "На мероприятие пока никто не записался."
	BookingsOfflineNotFoundMessage = "На офлайн мероприятие пока никто не записался."
	BookingsOnlineNotFoundMessage  = "На онлайн мероприятие пока никто не записался."
)

func GetBookingsMessageContent(
	event *entity.Event,
	createdBy *entity.User,
	updatedBy *entity.User,
	url string,
	bookingsOfflineCount int,
	bookingsOnlineCount int,
) []any {
	content := make([]any, 0)
	content = append(content, eventMessage(event, createdBy, updatedBy, url, bookingsOfflineCount, bookingsOnlineCount)...)
	content = append(content, getBookingsInlineKeyboard(event))
	return content
}

func getBookingsInlineKeyboard(
	event *entity.Event,
) *tele.ReplyMarkup {
	keyboard := [][]tele.InlineButton{
		[]tele.InlineButton{
			{
				Text:   "Оффлайн",
				Unique: callback.GetBookingsOffline,
				Data:   event.EventID,
			},
		},
		[]tele.InlineButton{
			{
				Text:   "Онлайн",
				Unique: callback.GetBookingsOnline,
				Data:   event.EventID,
			},
		},
		[]tele.InlineButton{
			{
				Text:   "< Назад",
				Unique: callback.Event,
				Data:   event.EventID,
			},
		},
	}
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}

func BookingMessageContent(booking *entity.Booking, user *entity.User) []any {
	parts := make([]string, 0)
	entities := make(tele.Entities, 0)
	curLen := 0

	telegram := "Телеграм:"
	var telegramURL string
	if len(user.Username) != 0 {
		telegramURL = fmt.Sprintf("https://t.me/%s", user.Username)
	} else {
		telegramURL = "-"
	}
	parts = append(parts,
		fmt.Sprintf("%s %s\n", telegram, telegramURL),
	)
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: 0,
			Length: len(utf16.Encode([]rune(telegram))),
		},
	)
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	name := "Имя:"
	if len(user.FirstName) != 0 {
		parts = append(parts,
			fmt.Sprintf("%s %s\n", name, user.FirstName),
		)
	} else {
		parts = append(parts,
			fmt.Sprintf("%s %s\n", name, "-"),
		)
	}
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(name))),
		},
	)
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	lastname := "Фамилия:"
	if len(user.LastName) != 0 {
		parts = append(parts,
			fmt.Sprintf("%s %s\n", lastname, user.LastName),
		)
	} else {
		parts = append(parts,
			fmt.Sprintf("%s %s\n", lastname, "-"),
		)
	}
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(lastname))),
		},
	)
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	contact := "Контакт для связи:"
	if len(user.ContactInfo) != 0 {
		parts = append(parts,
			fmt.Sprintf("%s %s\n", contact, user.ContactInfo),
		)
	} else {
		parts = append(parts,
			fmt.Sprintf("%s %s\n", contact, "-"),
		)
	}
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(contact))),
		},
	)
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

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
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	created := "Создано:"
	parts = append(parts,
		fmt.Sprintf("%s %s\n", created, booking.CreatedAt.Format("15:04 02.01.2006")),
	)
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(created))),
		},
	)
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	if booking.UpdatedAt != nil {
		updated := "Обновлено:"
		parts = append(parts,
			fmt.Sprintf("%s %s\n", updated, booking.UpdatedAt.Format("15:04 02.01.2006")),
		)
		entities = append(entities,
			tele.MessageEntity{
				Type:   tele.EntityBold,
				Offset: curLen,
				Length: len(utf16.Encode([]rune(updated))),
			},
		)
		curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))
	}

	text := strings.Join(parts, "")
	if booking.Payment == nil {
		return []any{text, entities}
	}
	photo := booking.Payment
	photo.Caption = text
	return []any{photo, entities}
}
