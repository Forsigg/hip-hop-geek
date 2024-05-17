package bot

import (
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"hip-hop-geek/internal/models"
)

func (b *TGBot) StartCommandHandler(upd tgbotapi.Update, user *models.User) {
	msg := tgbotapi.NewMessage(user.Id, StartCommandMessageText)
	adminId, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	keyboard := mainKeyboard
	if upd.Message.From.ID == adminId {
		keyboard = adminKeyboard
	}
	msg.ReplyMarkup = keyboard
	b.mustSend(msg)
}
