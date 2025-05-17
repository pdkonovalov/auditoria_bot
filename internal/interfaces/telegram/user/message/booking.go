package message

import (
	tele "gopkg.in/telebot.v4"
)

const (
	BookingAlredyBookedMessage = "Вы уже записаны на это мероприятие."
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

const (
	BookingWaitInputPaymentMessage             = "Пришлите скриншот оплаты мероприятия."
	BookingWaitInputPaymentInvalidInputMessage = "Пришлите скриншот оплаты или укажите, что заплатите потом, с помощью кнопки снизу."
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
