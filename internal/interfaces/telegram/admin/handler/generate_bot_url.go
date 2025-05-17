package handler

import "fmt"

func (h *AdminHandler) generateBotUrl(payload string) string {
	return fmt.Sprintf("https://t.me/%s?start=%s", h.botUsername, payload)
}
