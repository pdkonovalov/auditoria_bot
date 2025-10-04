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
	bookingsOfflineCount int,
	bookingsOnlineCount int,
) []any {
	content := make([]any, 0)
	content = append(content, eventMessage(event, createdBy, updatedBy, url, bookingsOfflineCount, bookingsOnlineCount)...)
	content = append(content, sendNotificationInlineKeyboard(event))
	return content
}

func sendNotificationInlineKeyboard(
	event *entity.Event,
) *tele.ReplyMarkup {
	keyboard := [][]tele.InlineButton{
		{
			{
				Text:   "Записавшимся оффлайн",
				Unique: callback.SendNotificationOffline,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
		{
			{
				Text:   "Записавшимся онлайн",
				Unique: callback.SendNotificationOnline,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
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

var (
	SendNotificationWaitInputPhotoTextMessage = "Пришлите сообщение уведомления."
	SendNotificationSuccessMessage            = "Уведомление отправлено."
)
