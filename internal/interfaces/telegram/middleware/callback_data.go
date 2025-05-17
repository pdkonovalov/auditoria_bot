package middleware

import (
	"strings"

	tele "gopkg.in/telebot.v4"
)

func CallbackData(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Callback() == nil {
			return next(c)
		}
		parts := strings.Split(c.Callback().Data, "|")
		if len(parts) >= 1 && len(parts[0]) >= 1 {
			c.Set("callback_unique", parts[0][1:])
		}
		if len(parts) >= 2 {
			c.Set("callback_data", parts[1])
		}
		return next(c)
	}
}
