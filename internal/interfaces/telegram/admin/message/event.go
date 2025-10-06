package message

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/admin/callback"
	"github.com/pdkonovalov/russian_time"

	tele "gopkg.in/telebot.v4"
)

func EventMessageContent(
	event *entity.Event,
	createdBy *entity.User,
	updatedBy *entity.User,
	url string,
	bookingsOfflineCount int,
	bookingsOnlineCount int,
) []any {
	content := make([]any, 0)
	content = append(content, eventMessage(event, createdBy, updatedBy, url, bookingsOfflineCount, bookingsOnlineCount)...)
	content = append(content, eventInlineKeyboard(event, bookingsOfflineCount, bookingsOnlineCount))
	return content
}

func eventMessage(
	event *entity.Event,
	createdBy *entity.User,
	updatedBy *entity.User,
	url string,
	bookingsOfflineCount int,
	bookingsOnlineCount int,
) []any {
	if event == nil {
		return nil
	}

	parts := make([]string, 0)
	entities := make(tele.Entities, 0)
	curLen := 0

	formatKey := "Формат:"
	var formatValue string
	if event.Offline && !event.Online {
		formatValue = "офлайн"
	} else if !event.Offline && event.Online {
		formatValue = "онлайн"
	} else {
		formatValue = "офлайн и онлайн"
	}
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(formatKey))),
		},
	)
	parts = append(parts, fmt.Sprintf("%s %s\n", formatKey, formatValue))
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	typeKey := "Тип:"
	var typeValue string
	if event.Offline && !event.Online {
		if event.OfflinePaid {
			typeValue = "платное"
		} else {
			typeValue = "бесплатное"
		}
	} else if !event.Offline && event.Online {
		if event.OnlinePaid {
			typeValue = "платное"
		} else {
			typeValue = "бесплатное"
		}
	} else if event.OfflinePaid == event.OnlinePaid {
		if event.OfflinePaid {
			typeValue = "платное"
		} else {
			typeValue = "бесплатное"
		}
	} else {
		var typeOfflineValue string
		if event.OfflinePaid {
			typeOfflineValue = "платно"
		} else {
			typeOfflineValue = "бесплатно"
		}
		var typeOnlineValue string
		if event.OnlinePaid {
			typeOnlineValue = "платно"
		} else {
			typeOnlineValue = "бесплатно"
		}
		typeValue = fmt.Sprintf("офлайн - %s, онлайн - %s", typeOfflineValue, typeOnlineValue)
	}
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(typeKey))),
		},
	)
	parts = append(parts, fmt.Sprintf("%s %s\n", typeKey, typeValue))
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	titleKey := "Название:"
	titleValue := event.Title
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(titleKey))),
		},
	)
	parts = append(parts, fmt.Sprintf("%s %s\n", titleKey, titleValue))
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	timeKey := "Время начала:"
	timeValue := event.Time.Format("15:04 02.01.2006")
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(timeKey))),
		},
	)
	parts = append(parts, fmt.Sprintf("%s %s\n", timeKey, timeValue))
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	bookingsKey := "Кол-во записавшихся:"
	var bookingsValue string
	if event.Offline && event.Online ||
		event.Offline && bookingsOnlineCount != 0 ||
		event.Online && bookingsOfflineCount != 0 {
		bookingsValue = fmt.Sprintf("офлайн - %v, онлайн - %v", bookingsOfflineCount, bookingsOnlineCount)
	} else if event.Offline {
		bookingsValue = fmt.Sprintf("%v", bookingsOfflineCount)
	} else if event.Online {
		bookingsValue = fmt.Sprintf("%v", bookingsOnlineCount)
	}
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(bookingsKey))),
		},
	)
	parts = append(parts, fmt.Sprintf("%s %s\n", bookingsKey, bookingsValue))
	curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))

	if event.OfflinePaid || event.OnlinePaid {
		paymentDetailsKey := "Реквизиты:"
		paymentDetailsValue := fmt.Sprintf("%s %s %s", event.PaymentDetailsFirstName, event.PaymentDetailsLastName, event.PaymentDetailsAccount)
		entities = append(entities,
			tele.MessageEntity{
				Type:   tele.EntityBold,
				Offset: curLen,
				Length: len(utf16.Encode([]rune(paymentDetailsKey))),
			},
		)
		parts = append(parts, fmt.Sprintf("%s %s\n", paymentDetailsKey, paymentDetailsValue))
		curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))
	}

	if createdBy != nil {
		createdKey := "Создано:"
		createdValue := fmt.Sprintf(
			"%s %s %s",
			createdBy.FirstName,
			createdBy.LastName,
			event.CreatedAt.Format("15:04 02.01.2006"),
		)
		entities = append(entities,
			tele.MessageEntity{
				Type:   tele.EntityBold,
				Offset: curLen,
				Length: len(utf16.Encode([]rune(createdKey))),
			},
		)
		parts = append(parts, fmt.Sprintf("%s %s\n", createdKey, createdValue))
		curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))
	}

	if event.UpdatedAt != nil && updatedBy != nil {
		updatedKey := "Обновлено:"
		updatedValue := fmt.Sprintf(
			"%s %s %s",
			updatedBy.FirstName,
			updatedBy.LastName,
			event.UpdatedAt.Format("15:04 02.01.2006"),
		)
		entities = append(entities,
			tele.MessageEntity{
				Type:   tele.EntityBold,
				Offset: curLen,
				Length: len(utf16.Encode([]rune(updatedKey))),
			},
		)
		parts = append(parts, fmt.Sprintf("%s %s\n", updatedKey, updatedValue))
		curLen += len(utf16.Encode([]rune(parts[len(parts)-1])))
	}

	urlKey := "Ссылка:"
	urlValue := url
	entities = append(entities,
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: curLen,
			Length: len(utf16.Encode([]rune(urlKey))),
		},
	)
	parts = append(parts, fmt.Sprintf("%s %s", urlKey, urlValue))

	text := strings.Join(parts, "")
	if event.Photo == nil {
		return []any{text, entities}
	}
	photo := event.Photo
	photo.Caption = text
	return []any{photo, entities}
}

func eventInlineKeyboard(
	event *entity.Event,
	bookingsOfflineCount int,
	bookingsOnlineCount int,
) *tele.ReplyMarkup {
	keyboard := make([][]tele.InlineButton, 0)
	keyboard = append(keyboard,
		[]tele.InlineButton{
			{
				Text:   "Пост",
				Unique: callback.EventPhotoText,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
	)
	getBookingsOffline := event.Offline || bookingsOfflineCount != 0
	getBookingsOnline := event.Online || bookingsOnlineCount != 0
	if getBookingsOffline && !getBookingsOnline {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Список записавшихся",
					Unique: callback.ShowBookings,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
						"format":  "offline",
						"page":    "0",
					}),
				},
			},
		)
	} else if !getBookingsOffline && getBookingsOnline {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Список записавшихся",
					Unique: callback.ShowBookings,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
						"format":  "online",
						"page":    "0",
					}),
				},
			},
		)
	} else {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Список записавшихся",
					Unique: callback.ShowBookingsFormatSelection,
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
				Text:   "Редактировать",
				Unique: callback.EditEvent,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
	)
	sendNotificationOffline := getBookingsOffline
	sendNotificationOnline := getBookingsOnline
	if sendNotificationOffline && !sendNotificationOnline {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Отправить уведомление",
					Unique: callback.SendNotification,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
						"format":  "offline",
					}),
				},
			},
		)
	} else if !sendNotificationOffline && sendNotificationOnline {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Отправить уведомление",
					Unique: callback.SendNotification,
					Data: callback.Encode(map[string]string{
						"eventID": event.EventID,
						"format":  "online",
					}),
				},
			},
		)
	} else {
		keyboard = append(keyboard,
			[]tele.InlineButton{
				{
					Text:   "Отправить уведомление",
					Unique: callback.SendNotificationFormatSelection,
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
				Text:   "Удалить",
				Unique: callback.DeleteEvent,
				Data: callback.Encode(map[string]string{
					"eventID": event.EventID,
				}),
			},
		},
		[]tele.InlineButton{
			{
				Text:   fmt.Sprintf("< Мероприятия %s", russian_time.DayMonth(&event.Time, russian_time.RCase2)),
				Unique: callback.EventsByDate,
				Data: callback.Encode(map[string]string{
					"date": event.Time.Format("02.01.2006"),
				}),
			},
		},
	)
	return &tele.ReplyMarkup{InlineKeyboard: keyboard}
}
