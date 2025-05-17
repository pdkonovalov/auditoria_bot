package command

import tele "gopkg.in/telebot.v4"

var (
	Events = tele.Command{
		Text:        "/events",
		Description: "все мерориятия",
	}
	MyEvents = tele.Command{
		Text:        "/myevents",
		Description: "мои мерориятия",
	}
	SetContactInfo = tele.Command{
		Text:        "/setcontact",
		Description: "изменить контакт для связи",
	}
	Cancel = tele.Command{
		Text:        "/cancel",
		Description: "отмена",
	}
	Help = tele.Command{
		Text:        "/help",
		Description: "помощь",
	}
	Admin = tele.Command{
		Text:        "/admin",
		Description: "переключиться в режим администратора",
	}
)

var List = []tele.Command{
	Events,
	MyEvents,
	SetContactInfo,
	Cancel,
	Help,
}

var AdminList = []tele.Command{
	Events,
	MyEvents,
	SetContactInfo,
	Admin,
	Cancel,
	Help,
}
