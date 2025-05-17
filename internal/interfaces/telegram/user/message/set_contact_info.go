package message

import (
	"fmt"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/callback"

	tele "gopkg.in/telebot.v4"
)

func SetContactInfoMessageContent(user *entity.User) []any {
	text := fmt.Sprintf("Контакт для связи: %s", user.ContactInfo)
	entities := tele.Entities{
		{
			Type:   tele.EntityBold,
			Offset: 0,
			Length: len(utf16.Encode([]rune("Контакт для связи:"))),
		},
	}
	return []any{text, entities, setContactInfoInlineKeyboard}
}

var setContactInfoInlineKeyboard = &tele.ReplyMarkup{
	InlineKeyboard: [][]tele.InlineButton{
		{
			{
				Text:   "Редактировать",
				Unique: callback.EditContactInfo,
			},
		},
	},
}

const (
	EditContactInfoWaitInputMessage                   = "Для записи на мероприятия пришлите удобный способ связи. Вы можете указать телефон, почту или ссылку на соцсеть."
	EditContactInfoInvalidInputMessage                = "Укажите телефон, почту или ссылку на соцсеть."
	EditContactInfSuccessMessage                      = "Данные обновлены."
	EditContactInfoWaitInputReplyKeyboardTelegramText = "Мой телеграм"
	EditContactInfoWaitInputTelegramNotExists         = "Ваше имя пользователя не задано или не доступно, пожалуйста, укажите другой способ связи"
)

func EditContactInfoWaitInputReplyKeyboard(username string) *tele.ReplyMarkup {
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
