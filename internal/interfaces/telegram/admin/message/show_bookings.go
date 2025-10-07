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

func ShowBookingsFormatSelectionMessageContent(
	event *entity.Event,
	createdBy *entity.User,
	updatedBy *entity.User,
	url string,
	bookingsOfflineCount int,
	bookingsOnlineCount int,
) []any {
	content := make([]any, 0)
	content = append(content, eventMessage(event, createdBy, updatedBy, url, bookingsOfflineCount, bookingsOnlineCount)...)
	content = append(content, showBookingsFormatSelectionInlineKeyboard(event))
	return content
}

func showBookingsFormatSelectionInlineKeyboard(
	event *entity.Event,
) *tele.ReplyMarkup {
	keyboard := [][]tele.InlineButton{
		{
			{
				Text:   "Оффлайн",
				Unique: callback.ShowBookings,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
					"format":  "offline",
					"page":    "0",
				}),
			},
		},
		{
			{
				Text:   "Онлайн",
				Unique: callback.ShowBookings,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
					"format":  "online",
					"page":    "0",
				}),
			},
		},
		{
			{
				Text:   "< Назад",
				Unique: callback.Event,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
	}
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}

func ShowBookingsMessageContent(
	event *entity.Event,
	page int,
	bookingsUsers []*entity.User,
	prevBookingsExists bool,
	nextBookingsExists bool,
	bookingsOfflineExists bool,
	bookingsOnlineExists bool,
	format string,
) []any {
	var text string
	if format == "offline" && (event.Online || bookingsOnlineExists) {
		text = "Список записавшихся оффлайн"
	} else if format == "online" && (event.Offline || bookingsOfflineExists) {
		text = "Список записавшихся онлайн"
	} else {
		text = "Список записавшихся"
	}

	return []any{
		text,
		showBookingsInlineKeyboard(
			event,
			page,
			bookingsUsers,
			prevBookingsExists,
			nextBookingsExists,
			bookingsOfflineExists,
			bookingsOnlineExists,
			format,
		),
	}

}

func showBookingsInlineKeyboard(
	event *entity.Event,
	page int,
	bookingsUsers []*entity.User,
	prevBookingsExists bool,
	nextBookingsExists bool,
	bookingsOfflineExists bool,
	bookingsOnlineExists bool,
	format string,
) *tele.ReplyMarkup {
	keyboard := make([][]tele.InlineButton, 0)
	for _, user := range bookingsUsers {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   fmt.Sprintf("%s %s", user.FirstName, user.LastName),
					Unique: callback.Booking,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
						"userID":  fmt.Sprintf("%v", user.UserID),
					}),
				},
			},
		)
	}
	if prevBookingsExists || nextBookingsExists {
		navigationRow := make([]tele.InlineButton, 0)
		if prevBookingsExists {
			navigationRow = append(navigationRow, tele.InlineButton{
				Text:   "<",
				Unique: callback.ShowBookings,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
					"format":  format,
					"page":    fmt.Sprintf("%v", page-1),
				}),
			})
		}
		if nextBookingsExists {
			navigationRow = append(navigationRow, tele.InlineButton{
				Text:   ">",
				Unique: callback.ShowBookings,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
					"format":  format,
					"page":    fmt.Sprintf("%v", page+1),
				}),
			})
		}
		keyboard = append(keyboard, navigationRow)
	}
	if (event.Offline || bookingsOfflineExists) &&
		(event.Online || bookingsOnlineExists) {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "< Назад",
					Unique: callback.ShowBookingsFormatSelection,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
					}),
				},
			},
		)
	} else {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "< Назад",
					Unique: callback.Event,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
					}),
				},
			},
		)
	}
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}

func BookingMessageContent(
	booking *entity.Booking,
	user *entity.User,
	page int,
) []any {
	content := make([]any, 0)
	content = append(content,
		bookingMessage(booking, user)...,
	)
	content = append(content,
		bookingInlineKeyboard(booking, page),
	)
	return content
}

func bookingMessage(
	booking *entity.Booking,
	user *entity.User,
) []any {
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
	if booking.Payment {
		if booking.PaymentPhoto != nil {
			parts = append(parts,
				fmt.Sprintf("%s %s\n", payment, "скриншот"),
			)
		} else if booking.PaymentDocument != nil {
			parts = append(parts,
				fmt.Sprintf("%s %s\n", payment, "документ"),
			)
		} else {
			parts = append(parts,
				fmt.Sprintf("%s %s\n", payment, "-"),
			)
		}
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

	if !booking.Payment {
		return []any{text, entities}
	}

	if booking.PaymentPhoto != nil {
		photo := booking.PaymentPhoto
		photo.Caption = text
		return []any{photo, entities}
	}

	if booking.PaymentDocument != nil {
		document := booking.PaymentDocument
		document.Caption = text
		return []any{document, entities}
	}

	return []any{text, entities}
}

func bookingInlineKeyboard(
	booking *entity.Booking,
	page int,
) *tele.ReplyMarkup {
	keyboard := make([][]tele.InlineButton, 0)
	var format string
	if booking.Offline {
		format = "offline"
	} else {
		format = "online"
	}
	keyboard = append(keyboard,
		[]tele.InlineButton{
			{
				Text:   "< Назад",
				Unique: callback.ShowBookings,
				Data: callback.Encode(map[string]string{
					"eventID": booking.EventID,
					"page":    fmt.Sprintf("%v", page),
					"format":  format,
				}),
			},
		},
	)
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}
