package bot

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"hip-hop-geek/internal/models"
)

func (b *TGBot) PreviousReleasesCallbackHandler(
	upd tgbotapi.Update,
	user *models.User,
	messageType models.MessageIdType,
) {
	now := time.Now().UTC()
	chatId := upd.CallbackQuery.From.ID
	// pageCount := user.ReleasesPageCount

	var releases []models.Release
	var pageCount int

	switch messageType {
	case models.ReleasesMessage:
		pageCount = user.ReleasesPageCount
	case models.TodayReleasesMessage:
		pageCount = user.TodayReleasesPageCount
	}
	var msgId int64

	if pageCount > 1 {
		pageCount -= 1
	} else {
		log.Fatalf("cant get previous page - its 1")
	}

	switch messageType {
	case models.ReleasesMessage:
		releases = b.Service.GetMonthReleases(
			now.Year(), now.Month(),
			StandardReleasesLimit,
			pageCount*StandardReleasesLimit,
		)
		msgId = user.ReleasesMessageId
	case models.TodayReleasesMessage:
		releases = b.Service.GetReleasesByDay(
			now.Year(), now.Month(), now.Day(),
			StandardReleasesLimit,
			pageCount*StandardReleasesLimit,
		)
		msgId = user.TodayReleasesMessageId
	}

	msg := GenerateReleasesEditMessage(releases)
	inlineKeyboard := GenerateInlineReleasesKeyboard(messageType, pageCount, releases)
	msgEdit := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatId,
			MessageID:   int(msgId),
			ReplyMarkup: &inlineKeyboard,
		},
		Media: msg,
	}

	doneMsg, err := b.Send(msgEdit)
	if err != nil {
		log.Fatalf("error while sending photo message: %s", err)
	}
	b.Service.SetUserState(user.Id, int(messageType), doneMsg.MessageID, pageCount)
}

func (b *TGBot) NextReleasesCallbackHandler(
	upd tgbotapi.Update,
	user *models.User,
	messageType models.MessageIdType,
) {
	now := time.Now().UTC()
	chatId := upd.CallbackQuery.From.ID
	var releases []models.Release
	var msgId int64
	var pageCount int

	switch messageType {
	case models.ReleasesMessage:
		pageCount = user.ReleasesPageCount
		msgId = user.ReleasesMessageId
	case models.TodayReleasesMessage:
		pageCount = user.TodayReleasesPageCount
		msgId = user.TodayReleasesMessageId
	}

	pageCount = pageCount + 1

	switch messageType {
	case models.ReleasesMessage:
		releases = b.Service.GetMonthReleases(
			now.Year(), now.Month(),
			StandardReleasesLimit,
			pageCount*StandardReleasesLimit,
		)
	case models.TodayReleasesMessage:
		releases = b.Service.GetReleasesByDay(
			now.Year(), now.Month(), now.Day(),
			StandardReleasesLimit,
			pageCount*StandardReleasesLimit,
		)
	}

	msg := GenerateReleasesEditMessage(releases)
	inlineKeyboard := GenerateInlineReleasesKeyboard(messageType, pageCount, releases)
	msgEdit := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatId,
			MessageID:   int(msgId),
			ReplyMarkup: &inlineKeyboard,
		},
		Media: msg,
	}

	doneMsg, err := b.Send(msgEdit)
	if err != nil {
		log.Fatalf("error while sending photo message: %s", err)
	}
	b.Service.SetUserState(user.Id, int(messageType), doneMsg.MessageID, pageCount)
}
