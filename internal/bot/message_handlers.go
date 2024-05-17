package bot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"hip-hop-geek/internal/db/sqlite"
	"hip-hop-geek/internal/models"
)

const (
	newReleasesPicUrl = "https://assets.sentaifilmworks.com/category-defaults/NewReleases.jpg"
	TextParseMode     = "HTML"

	NoOffset = 0
	NoLimit  = -1

	StandardReleasesLimit = 10
)

func (b *TGBot) ReleasesHandler(upd tgbotapi.Update, user *models.User) {
	log.Println("processing /releases command")
	chatId := upd.FromChat().ID
	now := time.Now().UTC()
	var needNewMessage bool

	if user.ReleasesPageCount == 0 {
		needNewMessage = true
	}
	pageCount := user.ReleasesPageCount

	if needNewMessage {
		pageCount = 1
		releases := b.Service.GetMonthReleases(
			now.Year(),
			now.Month(),
			StandardReleasesLimit,
			NoOffset,
		)
		msg := GenerateReleasesMessage(chatId, models.ReleasesMessage, pageCount, releases)

		doneMsg, err := b.Send(msg)
		if err != nil {
			log.Fatalf("error while sending photo message: %s", err)
		}
		b.Service.SetUserState(user.Id, models.ReleasesMessage, doneMsg.MessageID, pageCount)
		// b.SetUserState(chatId, int64(doneMsg.MessageID), pageCount)
		return
	} else {
		releases := b.Service.GetMonthReleases(
			now.Year(), now.Month(),
			StandardReleasesLimit,
			pageCount*StandardReleasesLimit,
		)
		deleteMsg := tgbotapi.NewDeleteMessage(chatId, int(user.ReleasesMessageId))
		b.Send(deleteMsg)
		msg := GenerateReleasesMessage(chatId, models.ReleasesMessage, pageCount, releases)

		doneMsg, err := b.Send(msg)
		if err != nil {
			log.Fatalf("error while sending photo message: %s", err)
		}
		b.Service.SetUserState(user.Id, models.ReleasesMessage, doneMsg.MessageID, pageCount)
		return
	}
}

func (b *TGBot) TodayEventHandler(chatId int64) {
	msg := tgbotapi.NewPhoto(chatId, nil)
	event, err := b.Service.GetTodayEvent()
	if err != nil {
		b.mustSend(tgbotapi.NewMessage(int64(chatId), ErrorUserMessage))
	}
	msg.File = tgbotapi.FileURL(event.Url)
	msg.Caption = fmt.Sprintf("%s\n%s", "Today in Hip Hop Hisory:", event.Text)

	b.mustSend(msg)
}

func (b *TGBot) echoMessage(upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
	b.mustSend(msg)
}

func (b *TGBot) UnsubscribeHandler(user *models.User) {
	b.Service.SetTodaySubscribe(user.Id, false)
	msg := tgbotapi.NewMessage(int64(user.Id), SuccessUnsubscribeMessage)
	b.mustSend(msg)
}

func (b *TGBot) SubscribeHandler(user *models.User) {
	b.Service.SetTodaySubscribe(user.Id, true)
	msg := tgbotapi.NewMessage(int64(user.Id), SuccessSubscribeMessage)
	b.mustSend(msg)
}

func (b *TGBot) CheckSubscribeHandler(user *models.User) {
	user, _ = b.Service.GetUserByUsername(user.Username)
	subscribeStatus := "не активна"
	if user.IsTodaySubscribe {
		subscribeStatus = "активна"
	}
	msg := tgbotapi.NewMessage(
		int64(user.Id),
		fmt.Sprintf(CheckSubscribeMessage, subscribeStatus),
	)
	b.mustSend(msg)
}

func (b *TGBot) SendTodayEventToSubscribers() {
	allSubs, err := b.Service.GetAllSubscribers()
	if err != nil {
		if errors.Is(err, sqlite.ErrUserNotFound) {
			log.Println("subscribers not found")
		} else {
			adminId, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
			msg := tgbotapi.NewMessage(int64(adminId), fmt.Sprintf(ErrorAdminMessage, err))
			b.mustSend(msg)
		}
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(allSubs))
	log.Println("sending today event to all subscribers")
	for _, subscriber := range allSubs {
		go func(subscriber *models.User) {
			defer wg.Done()
			b.TodayEventHandler(subscriber.Id)
		}(subscriber)
	}

	wg.Wait()
}

func (b *TGBot) SendTodayReleasesToSubscribers() {
	log.Println("sending today releases to subscribers")
	allSubs, err := b.Service.GetAllSubscribers()
	if err != nil {
		if errors.Is(err, sqlite.ErrUserNotFound) {
			log.Println("subscribers not found")
		} else {
			adminId, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
			msg := tgbotapi.NewMessage(int64(adminId), fmt.Sprintf(ErrorAdminMessage, err))
			b.mustSend(msg)
		}
		return
	}
	now := time.Now().UTC()
	releases := b.Service.GetReleasesByDay(now.Year(), now.Month(), now.Day(), NoLimit, NoOffset)

	// Если нет релизов то просто ничего не отправляем
	if len(releases) == 0 {
		log.Println("today releases not found, skip...")
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(allSubs))
	log.Println("sending today releases to all subscribers")
	for _, subscriber := range allSubs {
		go func(subscriber *models.User, releases []models.Release) {
			defer wg.Done()
			b.TodayReleasesHandler(subscriber, releases)
		}(subscriber, releases)
	}

	wg.Wait()
}

func (b *TGBot) TodayReleasesHandler(user *models.User, releases []models.Release) {
	log.Println("processing today releases")

	now := time.Now().UTC()
	var needNewMessage bool

	if user.TodayReleasesPageCount == 0 {
		needNewMessage = true
	}
	log.Println(user.TodayReleasesMessageId)
	log.Println(user.TodayReleasesPageCount)

	pageCount := user.TodayReleasesPageCount

	if needNewMessage {
		pageCount = 1
		releases := b.Service.GetReleasesByDay(
			now.Year(),
			now.Month(),
			now.Day(),
			StandardReleasesLimit,
			NoOffset,
		)
		if len(releases) == 0 {
			b.Send(tgbotapi.NewMessage(user.Id, "No today releases :("))
			return
		}
		msg := GenerateReleasesMessage(user.Id, models.TodayReleasesMessage, pageCount, releases)

		doneMsg, err := b.Send(msg)
		if err != nil {
			log.Fatalf("error while sending photo message: %s", err)
		}
		b.Service.SetUserState(user.Id, models.TodayReleasesMessage, doneMsg.MessageID, pageCount)
		// b.SetUserState(chatId, int64(doneMsg.MessageID), pageCount)
		return
	} else {
		releases := b.Service.GetReleasesByDay(
			now.Year(), now.Month(), now.Day(),
			StandardReleasesLimit,
			(pageCount-1)*StandardReleasesLimit,
		)
		if len(releases) == 0 {
			b.Send(tgbotapi.NewMessage(user.Id, "No today releases :("))
			return
		}
		deleteMsg := tgbotapi.NewDeleteMessage(user.Id, int(user.TodayReleasesMessageId))
		b.Send(deleteMsg)
		msg := GenerateReleasesMessage(user.Id, models.TodayReleasesMessage, pageCount, releases)

		doneMsg, err := b.Send(msg)
		if err != nil {
			log.Fatalf("error while sending photo message: %s", err)
		}
		b.Service.SetUserState(user.Id, models.TodayReleasesMessage, doneMsg.MessageID, pageCount)
		return
	}
	// for _, chunk := range chunkBy(releases, 10) {
	// 	media := make([]tgbotapi.InputMediaPhoto, 0)
	// 	caption := make([]string, 0)
	// 	for _, release := range chunk {
	//
	// 		imgCaption := GenerateCaptionForTodayRelease(release)
	// 		caption = append(caption, imgCaption)
	//
	// 		if release.CoverUrl.IsValid {
	// 			photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(release.CoverUrl.Value))
	// 			photo.Caption = imgCaption
	// 			photo.ParseMode = TextParseMode
	// 			media = append(media, photo)
	// 		}
	// 	}
	//
	// 	if len(media) != 0 {
	// 		media[0].Caption = fmt.Sprintf("Today releases:\n\n%s", strings.Join(caption, "\n\n"))
	// 		mediaInterfaces := make([]interface{}, len(media))
	// 		for i, photo := range media {
	// 			mediaInterfaces[i] = photo
	// 		}
	// 		mg := tgbotapi.NewMediaGroup(chatId, mediaInterfaces)
	// 		if _, err := b.SendMediaGroup(mg); err != nil {
	// 			log.Printf("error while sending media group: %s", err)
	// 		}
	//
	// 	} else {
	// 		photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(newReleasesPicUrl))
	// 		photo.ParseMode = TextParseMode
	// 		photo.Caption = fmt.Sprintf("Today releases:\n\n%s", strings.Join(caption, "\n\n"))
	// 		msg := tgbotapi.NewMediaGroup(chatId, []interface{}{photo})
	// 		if _, err := b.SendMediaGroup(msg); err != nil {
	// 			log.Printf("error while sending media group in today releases: %s", err)
	// 		}
	// 	}
	//
	// }
}

func (b *TGBot) RefreshReleasesHandler(years []int) {
	b.Updater.RefreshReleases(years)
}
