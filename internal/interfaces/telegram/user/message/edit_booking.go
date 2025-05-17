package message

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/callback"

	tele "gopkg.in/telebot.v4"
)

func EditBookingMessageContent(event *entity.Event, booking *entity.Booking, filter string) []any {
	parts := make([]string, 0)
	entities := make(tele.Entities, 0)
	curLen := 0

	editFormat := false

	if event.Offline && event.Online ||
		event.Offline && booking.Online ||
		event.Online && booking.Offline {
		editFormat = true
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

	editPayment := false
	if booking.Offline && event.OfflinePaid ||
		booking.Online && event.OnlinePaid ||
		booking.Payment != nil {
		editPayment = true
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
		return []any{text, entities, editBookingInlineKeyboard(booking, filter, editFormat, editPayment)}
	}
	photo := booking.Payment
	photo.Caption = text
	return []any{photo, entities, editBookingInlineKeyboard(booking, filter, editFormat, editPayment)}
}

func editBookingInlineKeyboard(
	booking *entity.Booking,
	filter string,
	editFormat bool,
	editPayment bool,
) *tele.ReplyMarkup {
	keyboard := make([][]tele.InlineButton, 0)
	if editFormat {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Изменить формат",
					Unique: callback.EditFormat,
					Data: callback.Encode(map[string]string{
						"eventID": booking.EventID,
						"filter":  filter,
					}),
				},
			},
		)
	}
	if editPayment {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Изменить скриншот оплаты",
					Unique: callback.EditPayment,
					Data: callback.Encode(map[string]string{
						"eventID": booking.EventID,
						"filter":  filter,
					}),
				},
			},
		)
	}
	keyboard = append(keyboard,
		[]tele.InlineButton{
			{
				Text:   "Изменить дополнительную информацию",
				Unique: callback.EditAdditionalInfo,
				Data: callback.Encode(map[string]string{
					"eventID": booking.EventID,
					"filter":  filter,
				}),
			},
		},
		[]tele.InlineButton{
			{
				Text:   "< Назад",
				Unique: callback.Event,
				Data: callback.Encode(map[string]string{
					"eventID": booking.EventID,
					"filter":  filter,
				}),
			},
		},
	)
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}

const (
	EditBookingWaitInputFormatMessage                  = "Укажите формат, на который хотите записаться."
	EditBookingWaitInputFormatInvalidInputMessage      = "Укажите формат мероприятия, с помощью кнопки снизу."
	EditBookingWaitInputFormatReplyKeyboardOfflineText = "Офлайн"
	EditBookingWaitInputFormatReplyKeyboardOnlineText  = "Онлайн"
)

var (
	EditBookingWaitInputFormatReplyKeyboard = &tele.ReplyMarkup{
		ReplyKeyboard: [][]tele.ReplyButton{
			{
				{
					Text: BookingWaitInputFormatReplyKeyboardOfflineText,
				},
				{
					Text: BookingWaitInputFormatReplyKeyboardOnlineText,
				},
			},
		},
		OneTimeKeyboard: true,
	}
)

const (
	EditBookingWaitInputPaymentMessage             = "Пришлите скриншот оплаты мероприятия."
	EditBookingWaitInputPaymentInvalidInputMessage = "Пришлите скриншот оплаты или укажите, что заплатите потом, с помощью кнопки снизу."
)

var (
	EditBookingWaitInputPaymentReplyKeyboardText = "Заплачу потом"
	EditBookingWaitInputPaymentReplyKeyboard     = &tele.ReplyMarkup{
		ReplyKeyboard: [][]tele.ReplyButton{
			{
				{
					Text: BookingWaitInputPaymentReplyKeyboardText,
				},
			},
		},
		OneTimeKeyboard: true,
	}
)

const (
	EditBookingWaitInputAdditionalInfoMessage           = "Укажите дополнительную информацию для организатора, если нужно."
	EditBookingWaitInputAdditionalInfoReplyKeyboardText = "Не нужно."
)

var EditBookingWaitInputAdditionalInfoReplyKeyboard = &tele.ReplyMarkup{
	ReplyKeyboard: [][]tele.ReplyButton{
		{
			{
				Text: BookingWaitInputAdditionalInfoReplyKeyboardText,
			},
		},
	},
	OneTimeKeyboard: true,
}

const (
	EditBookingSuccessMessage = "Данные обновлены."
)
