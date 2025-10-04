package message

import (
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/callback"

	tele "gopkg.in/telebot.v4"
)

func EditEventMessageContent(
	event *entity.Event,
	createdBy *entity.User,
	updatedBy *entity.User,
	url string,
	bookingsOfflineCount int,
	bookingsOnlineCount int,
) []any {
	content := make([]any, 0)
	content = append(content, eventMessage(event, createdBy, updatedBy, url, bookingsOfflineCount, bookingsOnlineCount)...)
	content = append(content, editEventInlineKeyboard(event))
	return content
}

func editEventInlineKeyboard(event *entity.Event) *tele.ReplyMarkup {
	keyboard := make([][]tele.InlineButton, 0)
	keyboard = append(keyboard,
		[]tele.InlineButton{
			{
				Text:   "Редактировать формат",
				Unique: callback.EditFormat,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
	)
	if event.Offline && !event.Online {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Редактировать тип",
					Unique: callback.EditOfflinePaid,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
					}),
				},
			},
		)
	} else if !event.Offline && event.Online {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Редактировать тип",
					Unique: callback.EditOnlinePaid,
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
					Text:   "Редактировать тип оффлайн",
					Unique: callback.EditOfflinePaid,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
					}),
				},
			},
			[]tele.InlineButton{
				{
					Text:   "Редактировать тип онлайн",
					Unique: callback.EditOnlinePaid,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
					}),
				},
			},
		)
	}
	keyboard = append(keyboard,
		[]tele.InlineButton{
			{
				Text:   "Редактировать название",
				Unique: callback.EditTitle,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
		[]tele.InlineButton{
			{
				Text:   "Редактировать время начала",
				Unique: callback.EditTime,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
	)
	if event.OfflinePaid || event.OnlinePaid {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Редактировать реквизиты",
					Unique: callback.EditPaymentDetails,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
					}),
				},
			},
		)
	}
	keyboard = append(keyboard,
		[]tele.InlineButton{
			{
				Text:   "Редактировать пост",
				Unique: callback.EditPhotoText,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
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
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}

var EditEventSuccessMessage = "Данные обновлены"
