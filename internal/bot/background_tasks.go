package bot

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"
)

func (b *TGBot) SendEventAndReleasesEveryday(ctx context.Context) {
	now := time.Now().Local()
	send_hour, _ := strconv.Atoi(os.Getenv("SEND_SUBS_HOUR"))
	send_minute, _ := strconv.Atoi(os.Getenv("SEND_SUBS_MINUTE"))
	next := time.Date(now.Year(), now.Month(), now.Day(), send_hour, send_minute, 0, 0, time.Local)
	
  if now.After(next) {
		next = next.Add(24 * time.Hour)
	}
	log.Printf("next send to subscribers: %d:%2d %2d.%2d.%d",
		next.Hour(), next.Minute(),
		next.Day(), next.Month(), next.Year(),
	)

	time.Sleep(time.Until(next))

	ticker := time.NewTicker(20 * time.Second)
	b.SendTodayEventToSubscribers()
	b.SendTodayReleasesToSubscribers()
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("bot timer goroutine closing...")
			return

		// sending message to subscribers every $DURATION
		case <-ticker.C:
			b.SendTodayEventToSubscribers()
			b.SendTodayReleasesToSubscribers()
		}
	}
}
