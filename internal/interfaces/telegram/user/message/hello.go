package message

import (
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram/user/command"

	tele "gopkg.in/telebot.v4"
)

const HelloText = "Привет! Это бот аудитории."

var (
	HelloMessage = func() string {
		parts := make([]string, len(command.List)+2, len(command.List)+2)
		parts[0] = HelloText
		parts[1] = ""
		for index, cmd := range command.List {
			parts[index+2] = fmt.Sprintf("%s - %s", cmd.Text, cmd.Description)
		}
		return strings.Join(parts, "\n")
	}()
	HelloAdminMessage = func() string {
		parts := make([]string, len(command.AdminList)+2, len(command.AdminList)+2)
		parts[0] = HelloText
		parts[1] = ""
		for index, cmd := range command.AdminList {
			parts[index+2] = fmt.Sprintf("%s - %s", cmd.Text, cmd.Description)
		}
		return strings.Join(parts, "\n")
	}()
	HelloEntities = tele.Entities{
		tele.MessageEntity{
			Type:   tele.EntityBold,
			Offset: 0,
			Length: len(utf16.Encode([]rune(HelloText))),
		},
	}
)
