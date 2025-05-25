package message

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/config"
	tele "gopkg.in/telebot.v4"
)

var (
	WaitInputFormatMessage                          = "Укажите формат мероприятия."
	WaitInputFormatReplyKeyboardButtonOffline       = "офлайн"
	WaitInputFormatReplyKeyboardButtonOnline        = "онлайн"
	WaitInputFormatReplyKeyboardButtonOfflineOnline = "офлайн и онлайн"
	WaitInputFormatReplyKeyboard                    = &tele.ReplyMarkup{
		ReplyKeyboard: [][]tele.ReplyButton{
			{
				{
					Text: WaitInputFormatReplyKeyboardButtonOffline,
				},
				{
					Text: WaitInputFormatReplyKeyboardButtonOnline,
				},
				{
					Text: WaitInputFormatReplyKeyboardButtonOfflineOnline,
				},
			},
		},
		OneTimeKeyboard: true,
	}
	WaitInputFormatInvalidInputMessage = "Укажите формат офлайн, онлайн или офлайн и онлайн, с помощью кнопки снизу."
)

var (
	WaitInputPaidMessage                  = "Укажите тип мероприятия"
	WaitInputOfflinePaidMessage           = "Укажите тип оффлайн мероприятия"
	WaitInputOnlinePaidMessage            = "Укажите тип онлайн мероприятия"
	WaitInputPaidReplyKeyboardButtonTrue  = "Платное"
	WaitInputPaidReplyKeyboardButtonFalse = "Бесплатное"
	WaitInputPaidReplyKeyboard            = &tele.ReplyMarkup{
		ReplyKeyboard: [][]tele.ReplyButton{
			{
				{
					Text: WaitInputPaidReplyKeyboardButtonTrue,
				},
				{
					Text: WaitInputPaidReplyKeyboardButtonFalse,
				},
			},
		},
		OneTimeKeyboard: true,
	}
	WaitInputPaidInvalidInputMessage = "Укажите тип - платное/бесплатное, с помощью кнопки снизу"
)

var (
	WaitInputPaymentDetailsMessage             = "Укажите реквизиты оплаты."
	WaitInputPaymentDetailsInvalidInputMessage = "Укажите реквизиты оплаты в формате: имя фамилия номер счёта. Или выберете реквизиты по умолчанию с помощью кнопки ниже."
)

func WaitInputPaymentDetailsReplyKeyboard(defaultPaymentDetailsList config.PaymentDetailsList) *tele.ReplyMarkup {
	keyboard := make([][]tele.ReplyButton, 0)
	for _, defaultPaymentDetails := range defaultPaymentDetailsList {
		defaultPaymentDetailsStr := fmt.Sprintf("%s %s %s", defaultPaymentDetails.FirstName, defaultPaymentDetails.LastName, defaultPaymentDetails.Account)
		keyboard = append(keyboard,
			[]tele.ReplyButton{
				{
					Text: defaultPaymentDetailsStr,
				},
			},
		)
	}
	return &tele.ReplyMarkup{
		ReplyKeyboard:   keyboard,
		OneTimeKeyboard: true,
	}
}

var (
	WaitInputTitleMessage = "Укажите название мероприятия"
)

var (
	WaitInputTimeMessage             = "Укажите время и дату мероприятия"
	WaitInputTimeInvalidInputMessage = "Укажите время и дату в формате 15:04 02.01.2006"
)

var (
	WaitInputPhotoTextMessage = "Пришлите пост мероприятия"
)

func WaitInputPhotoTextInvalidInputMessage(captionLen int) string {
	return fmt.Sprintf("Длинна поста слишком большая. Колличество символов - %v. Допустимая длинна 1024 символа. Пришлите более короткий пост.", captionLen)
}

func NewEventSuccessMessage(link string) string {
	return fmt.Sprintf("Мероприятие создано.\nСсылка: %s", link)
}
