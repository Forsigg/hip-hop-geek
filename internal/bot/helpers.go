package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"hip-hop-geek/internal/models"
	"hip-hop-geek/internal/utils"
)

func GenerateCaption(release models.Release) string {
	emoji := SingleEmoji
	if release.Type == models.Album {
		emoji = AlbumEmoji
	}
	imgCaption := fmt.Sprintf(
		"%s <b>%s - %s</b> (<i>%d %s %d</i>)",
		emoji,
		release.Artist.Name,
		release.Title,
		release.OutDate.Day(),
		release.OutDate.Month().String(),
		release.OutDate.Year(),
	)

	return imgCaption
}

func GenerateCaptionForTodayRelease(release models.Release) string {
	emoji := SingleEmoji
	if release.Type == models.Album {
		emoji = AlbumEmoji
	}
	imgCaption := fmt.Sprintf(
		"%s <b>%s - %s</b>",
		emoji,
		release.Artist.Name,
		release.Title,
	)

	return imgCaption
}

func (b *TGBot) mustSend(msg tgbotapi.Chattable) {
	if _, err := b.Send(msg); err != nil {
		log.Fatalf("error while sending message %v: %s", msg, err)
	}
}

func GenerateReleasesMessage(
	userId int64,
	messageType models.MessageIdType,
	pageCount int,
	releases []models.Release,
) tgbotapi.PhotoConfig {
	inlineKeyboard := GenerateInlineReleasesKeyboard(messageType, pageCount, releases)

	photoUrl := newReleasesPicUrl
	for _, release := range releases {
		if release.CoverUrl.IsValid {
			photoUrl = release.CoverUrl.Value
			break
		}
	}

	caption := make([]string, 0, 10)

	for _, release := range releases {
		caption = append(caption, GenerateCaption(release))
	}

	photoMsg := tgbotapi.NewPhoto(userId, tgbotapi.FileURL(photoUrl))
	photoMsg.Caption = strings.Join(caption, "\n\n")
	photoMsg.ParseMode = tgbotapi.ModeHTML
	photoMsg.ReplyMarkup = inlineKeyboard

	return photoMsg
}

func GenerateReleasesEditMessage(
	releases []models.Release,
) tgbotapi.InputMediaPhoto {
	photoUrl := newReleasesPicUrl
	for _, release := range releases {
		if release.CoverUrl.IsValid {
			photoUrl = release.CoverUrl.Value
			break
		}
	}

	caption := make([]string, 0, 10)

	for _, release := range releases {
		caption = append(caption, GenerateCaption(release))
	}

	photoMsg := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(photoUrl))
	photoMsg.Caption = strings.Join(caption, "\n\n")
	photoMsg.ParseMode = tgbotapi.ModeHTML

	return photoMsg
}

func GenerateInlineReleasesKeyboard(
	messageType models.MessageIdType,
	pageCount int,
	releases []models.Release,
) tgbotapi.InlineKeyboardMarkup {
	inlineButtons := make([]tgbotapi.InlineKeyboardButton, 0, 3)

	if pageCount > 1 {
		switch messageType {
		case models.ReleasesMessage:
			inlineButtons = append(
				inlineButtons,
				tgbotapi.NewInlineKeyboardButtonData(
					PrevReleasesButtonText,
					PreviousReleasesCallbackText,
				),
			)
		case models.TodayReleasesMessage:
			inlineButtons = append(
				inlineButtons,
				tgbotapi.NewInlineKeyboardButtonData(
					PrevReleasesButtonText,
					PreviousTodayReleasesCallbackText,
				),
			)
		}
	}

	var emojiPage string
	pageCountStr := strconv.Itoa(pageCount)
	for _, digit := range pageCountStr {
		intDigit, _ := strconv.Atoi(string(digit))
		emojiPage += NumbersToEmojiMapping[intDigit]
	}
	inlineButtons = append(
		inlineButtons,
		tgbotapi.NewInlineKeyboardButtonData(
			emojiPage,
			PageCountCallbackText,
		),
	)

	if len(releases) == StandardReleasesLimit {
		switch messageType {
		case models.ReleasesMessage:
			inlineButtons = append(
				inlineButtons,
				tgbotapi.NewInlineKeyboardButtonData(
					NextReleasesButtonText,
					NextReleasesCallbackText,
				),
			)
		case models.TodayReleasesMessage:
			inlineButtons = append(
				inlineButtons,
				tgbotapi.NewInlineKeyboardButtonData(
					NextReleasesButtonText,
					NextTodayReleasesCallbackText,
				),
			)
		}
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(inlineButtons...),
	)
	return inlineKeyboard
}

func GenerateYearByMonthKeyboard(
	year int,
) tgbotapi.InlineKeyboardMarkup {
	buttons := make([]tgbotapi.InlineKeyboardButton, 0, 12)

	for monthStr, monthInt := range utils.AllMonthsMapping {
		log.Println(monthStr)
		log.Println(monthInt)
		buttons = append(
			buttons,
			tgbotapi.NewInlineKeyboardButtonData(
				monthStr,
				fmt.Sprintf("%d.%d", year, monthInt),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			buttons[0:3]...,
		),
		tgbotapi.NewInlineKeyboardRow(
			buttons[3:6]...,
		),
		tgbotapi.NewInlineKeyboardRow(
			buttons[6:9]...,
		),
		tgbotapi.NewInlineKeyboardRow(
			buttons[9:12]...,
		),
	)
}
