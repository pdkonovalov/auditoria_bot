package command

import tele "gopkg.in/telebot.v4"

var (
	Events = tele.Command{
		Text:        "/events",
		Description: "все мероприятия",
	}
	NewEvent = tele.Command{
		Text:        "/newevent",
		Description: "создать мероприятие",
	}
	User = tele.Command{
		Text:        "/user",
		Description: "переключиться в режим пользователя",
	}
	Cancel = tele.Command{
		Text:        "/cancel",
		Description: "отмена",
	}
	Help = tele.Command{
		Text:        "/help",
		Description: "помощь",
	}
	List = []tele.Command{
		Events,
		NewEvent,
		User,
		Cancel,
		Help,
	}
)
