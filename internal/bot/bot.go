package bot

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"hip-hop-geek/internal/db/sqlite"
	"hip-hop-geek/internal/models"
)

type HipHopService interface {
	GetMonthReleases(year int, month time.Month, limit, offset int) []models.Release
	GetAllYearReleases(year, limit, offset int) []models.Release
	GetAllYearSingles(year int, withCover bool) []models.Release
	GetTodayEvent() (*models.TodayPost, error)
	GetReleasesByDay(year int, month time.Month, day, limit, offset int) []models.Release
	AddUser(user models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetAllSubscribers() ([]*models.User, error)
	SetTodaySubscribe(userId int64, isSubscribe bool) error
	SetUserState(userId int64, messageType, messageId int, pageCount int) error
	Close()
}

type UpdaterInterface interface {
	RefreshReleases(years []int)
}

type TGBot struct {
	*tgbotapi.BotAPI
	// redisClient *redis.Client
	Service HipHopService
	Updater UpdaterInterface
}

func NewTGBot(botToken string, service HipHopService, updater UpdaterInterface) *TGBot {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// redisPort := os.Getenv("REDIS_PORT")
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     fmt.Sprintf("redis:%s", redisPort),
	// 	Password: "",
	// 	DB:       0,
	// })
	// ctx := context.Background()
	// if err := redisClient.Ping(ctx).Err(); err != nil {
	// 	log.Fatalf("error while ping redis: %s", err)
	// }
	return &TGBot{
		bot,
		// redisClient,
		service,
		updater,
	}
}

func (b *TGBot) Start(ctx context.Context, timeout int) {
	log.Println("bot start polling...")
	updatesConfig := tgbotapi.NewUpdate(0)
	updatesConfig.Timeout = timeout

	updates := b.GetUpdatesChan(updatesConfig)

	for {
		select {
		case <-ctx.Done():
			log.Println("bot closing...")
			return

		// sending message to subscribers every $DURATION
		case <-time.After(24 * time.Hour):
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				b.SendTodayEventToSubscribers()
				b.SendTodayReleasesToSubscribers()
			}()
			wg.Wait()

		// handling all updates
		case upd := <-updates:
			// get or create user
			chat := upd.FromChat()
			user, err := b.Service.GetUserByUsername(chat.UserName)
			log.Printf("user before handle: %v", user)
			if err != nil {
				if err == sqlite.ErrUserNotFound {
					user = &models.User{
						Id:               chat.ID,
						Username:         chat.UserName,
						IsTodaySubscribe: false,
					}
					err = b.Service.AddUser(*user)

					if err != nil {
						log.Fatal(err)
					}
				} else {
					log.Fatal(err)
				}
			}

			if upd.Message != nil {
				log.Printf(
					"received message update from ID %d with text %s",
					upd.Message.From.ID,
					upd.Message.Text,
				)
				go b.messageHandler(upd, user)

			} else if upd.CallbackQuery != nil {
				log.Printf(
					"received callback update from ID %d with data %s",
					upd.CallbackQuery.Message.From.ID,
					upd.CallbackData(),
				)
				go b.callbackHandler(upd, user)
			}
		}
	}
}

func (b *TGBot) messageHandler(upd tgbotapi.Update, user *models.User) {
	adminId, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("user: %v", user)
				log.Printf("message: %v", upd.Message)
				panic(r)
			}
		}()
		deleteUserMsg := tgbotapi.NewDeleteMessage(user.Id, upd.Message.MessageID)
		b.Send(deleteUserMsg)
	}()

	if upd.Message.IsCommand() {
		b.commandHandler(upd, user)
		return
	}

	switch upd.Message.Text {
	case TodayButtonText:
		b.TodayEventHandler(upd.Message.Chat.ID)

	case TodayReleasesButtonText:
		chatId := upd.Message.Chat.ID
		now := time.Now().UTC()
		releases := b.Service.GetReleasesByDay(now.Year(), now.Month(), now.Day(), -1, 0)
		if releases == nil {
			b.mustSend(tgbotapi.NewMessage(chatId, "No releases today :("))
		} else {
			b.TodayReleasesHandler(user, releases)
		}

	case MonthReleasesButtonText:
		b.ReleasesHandler(upd, user)

	case SubscribeButtonText:
		b.SubscribeHandler(user)

	case UnsubscribeButtonText:
		b.UnsubscribeHandler(user)

	case CheckSubscribeButtonText:
		b.CheckSubscribeHandler(user)

	case RefreshReleasesButtonText:
		if user.Id != adminId {
			return
		}

		b.mustSend(tgbotapi.NewMessage(int64(adminId), RefreshReleasesStartMessageText))
		b.RefreshReleasesHandler([]int{2023, 2024})
		b.mustSend(tgbotapi.NewMessage(int64(adminId), RefreshReleasesEndMessageText))

	case TestButtonText:
		if user.Id != adminId {
			return
		}

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("test", "test"),
			),
		)

		photoMsg := tgbotapi.NewPhoto(int64(user.Id), tgbotapi.FileURL(newReleasesPicUrl))
		photoMsg.Caption = "Test Caption"
		photoMsg.ReplyMarkup = inlineKeyboard

		b.Send(photoMsg)

	}

	log.Println("handled update")
}

func (b *TGBot) commandHandler(upd tgbotapi.Update, user *models.User) {
	switch upd.Message.Command() {
	case StartCommandText:
		b.StartCommandHandler(upd, user)
	}
}

func (b *TGBot) callbackHandler(upd tgbotapi.Update, user *models.User) {
	switch upd.CallbackData() {
	case PreviousReleasesCallbackText:
		b.PreviousReleasesCallbackHandler(upd, user, models.ReleasesMessage)
	case NextReleasesCallbackText:
		b.NextReleasesCallbackHandler(upd, user, models.ReleasesMessage)
	case NextTodayReleasesCallbackText:
		b.NextReleasesCallbackHandler(upd, user, models.TodayReleasesMessage)
	case PreviousTodayReleasesCallbackText:
		b.PreviousReleasesCallbackHandler(upd, user, models.TodayReleasesMessage)
	case PageCountCallbackText:
		callback := tgbotapi.NewCallback(upd.CallbackQuery.ID, "")
		b.Send(callback)
	}
}