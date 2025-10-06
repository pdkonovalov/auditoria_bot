package message

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	tele "gopkg.in/telebot.v4"
)

const (
	BookingAlredyBookedMessage = "Вы уже записаны на это мероприятие."
)

const (
	BookingInitMessage = `Обратите внимание, что процесс регистрации завершится, когда вы получите сообщение об успешной регистрации и ваша запись появится в разделе "мои мероприятия". В этом разделе вы сможете в любой момент изменить данные записи или отменить её.`
)

const (
	BookingWaitInputContactInfoMessage                   = "Для записи на мероприятия пришлите удобный способ связи. Вы можете указать телефон, почту или ссылку на соцсеть."
	BookingWaitInputContactInfoInvalidInputMessage       = "Укажите телефон, почту или ссылку на соцсеть."
	BookingWaitInputContactInfoReplyKeyboardTelegramText = "Мой телеграм"
	BookingWaitInputContactInfoTelegramNotExists         = "Ваше имя пользователя не задано или не доступно, пожалуйста, укажите другой способ связи."
	BookingContactInfoSuccessMessage                     = "Данные добавлены."
)

func BookingWaitInputContactInfoReplyKeyboard(username string) *tele.ReplyMarkup {
	var keyboard [][]tele.ReplyButton
	if len(username) != 0 {
		keyboard = [][]tele.ReplyButton{
			{
				{
					Text: BookingWaitInputContactInfoReplyKeyboardTelegramText,
				},
				{
					Text:    "Мой телефон",
					Contact: true,
				},
			},
		}
	} else {
		keyboard = [][]tele.ReplyButton{
			{
				{
					Text:    "Мой телефон",
					Contact: true,
				},
			},
		}
	}
	return &tele.ReplyMarkup{
		ReplyKeyboard:   keyboard,
		OneTimeKeyboard: true,
	}
}

const (
	BookingWaitInputFormatMessage                  = "Укажите формат, на который хотите записаться."
	BookingWaitInputFormatInvalidInputMessage      = "Укажите формат мероприятия, с помощью кнопки снизу."
	BookingWaitInputFormatReplyKeyboardOfflineText = "Офлайн"
	BookingWaitInputFormatReplyKeyboardOnlineText  = "Онлайн"
)

var (
	BookingWaitInputFormatReplyKeyboard = &tele.ReplyMarkup{
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

func BookingWaitInputPaymentMessageContent(event *entity.Event) []any {
	parts := make([]string, 0)
	entities := make(tele.Entities, 0)
	curLen := 0

	messageFirstPart := "Пришлите подтверждение оплаты мероприятия в формате фото или pdf.\n\nРеквизиты:\n"
	parts = append(parts, messageFirstPart)
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	name := fmt.Sprintf("%s %s\n", event.PaymentDetailsFirstName, event.PaymentDetailsLastName)
	parts = append(parts, name)
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityItalic,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(name))),
		},
	)
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	account := event.PaymentDetailsAccount
	parts = append(parts, account)
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityCode,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(account))),
		},
	)

	text := strings.Join(parts, "")
	return []any{text, entities, BookingWaitInputPaymentReplyKeyboard}
}

const (
	BookingWaitInputPaymentInvalidInputMessage = "Пришлите подтверждение оплаты или укажите, что заплатите потом, с помощью кнопки снизу."
	BookingWaitInputPaymentReplyKeyboardText   = "Заплачу потом"
)

var (
	BookingWaitInputPaymentReplyKeyboard = &tele.ReplyMarkup{
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
	BookingWaitInputAdditionalInfoMessage           = "Укажите дополнительную информацию для организатора, если нужно."
	BookingWaitInputAdditionalInfoReplyKeyboardText = "Не нужно"
)

var BookingWaitInputAdditionalInfoReplyKeyboard = &tele.ReplyMarkup{
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
	BookingSuccessMessage = "Вы записаны на мероприятие."
)
