package message

import (
	"fmt"

	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/command"
)

const EventNotFoundMessage = "Мероприятие не найдено :("

var (
	MyEventsNotFoundMessage = fmt.Sprintf("Вы пока не записались ни на одно мероприятие.\nНажмите %s, чтобы посмотреть мероприятия аудитории.", command.Events.Text)
)
