package message

import (
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/callback"
	tele "gopkg.in/telebot.v4"
)

func SendNotificationMessageContent(
	event *entity.Event,
	createdBy *entity.User,
	updatedBy *entity.User,
	url string,
	sendNotificationOffline bool,
	sendNotificationOnline bool,
) []any {
	content := make([]any, 0)
	content = append(content, eventMessage(event, createdBy, updatedBy, url)...)
	content = append(content, sendNotificationInlineKeyboard(event, sendNotificationOffline, sendNotificationOnline))
	return content
}

func sendNotificationInlineKeyboard(
	event *entity.Event,
	sendNotificationOffline bool,
	sendNotificationOnline bool,
) *tele.ReplyMarkup {
	keyboard := make([][]tele.InlineButton, 0)
	if sendNotificationOffline {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Записавшимся оффлайн",
					Unique: callback.SendNotificationOffline,
					Data:   event.EventID,
				},
			},
		)
	}
	if sendNotificationOnline {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Записавшимся онлайн",
					Unique: callback.SendNotificationOnline,
					Data:   event.EventID,
				},
			},
		)
	}
	keyboard = append(keyboard,
		[]tele.InlineButton{
			{
				Text:   "< Назад",
				Unique: callback.Event,
				Data:   event.EventID,
			},
		},
	)
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}

var (
	SendNotificationWaitInputPhotoTextMessage = "Пришлите сообщение уведомления."
	SendNotificationSuccessMessage            = "Уведомление отправлено."
)
