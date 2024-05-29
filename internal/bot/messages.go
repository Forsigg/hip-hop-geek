package bot

const (
	// MESSAGES
	ErrorAdminMessage               = "Произошла ошибка: %s"
	ErrorUserMessage                = "Во время обработки сообщения произошла ошибка на сервере."
	SuccessSubscribeMessage         = "Вы успешно подписались на ежедневную рассылку истории Хип Хопа."
	SuccessUnsubscribeMessage       = "Вы успешно отписались от ежедневной рассылки  истории Хип Хопа."
	CheckSubscribeMessage           = "Ваша подписка %s"
	ReleasesNotFoundMessage         = "Релизы не найдены"
	StartCommandMessageText         = "Привет! Я бот - Хип Хоп гик, который знает обо всех релизах и событиях в жизни хип хопа. Отправляю тебе клавиатуру с нужными командами"
	RefreshReleasesStartMessageText = "Запускаю обновление релизов"
	RefreshReleasesEndMessageText   = "Обновление релизов завершено"

	SingleEmoji = "🎤"
	AlbumEmoji  = "💿"

	// BUTTONS
	TodayButtonText               = "Today in Hip Hop History"
	TodayReleasesButtonText       = "Today releases"
	MonthReleasesButtonText       = "Month releases"
	YearReleasesByMonthButtonText = "Year releases by month"

	SubscribeButtonText      = "Subscribe"
	UnsubscribeButtonText    = "Unsubscribe"
	CheckSubscribeButtonText = "Check subscribe"

	RefreshReleasesButtonText = "Manual refresh releases"
	TestButtonText            = "Test message"

	// COMMANDS
	StartCommandText = "start"

	// CALLBACKS
	PrevReleasesButtonText = "⬅️"
	NextReleasesButtonText = "➡️"

	PreviousReleasesCallbackText      = "prev_releases"
	NextReleasesCallbackText          = "next_releases"
	PreviousTodayReleasesCallbackText = "prev_today_releases"
	NextTodayReleasesCallbackText     = "next_today_releases"
	PageCountCallbackText             = "pageCount"
)

var NumbersToEmojiMapping = map[int]string{
	0: "0️⃣",
	1: "1️⃣",
	2: "2️⃣",
	3: "3️⃣",
	4: "4️⃣",
	5: "5️⃣",
	6: "6️⃣",
	7: "7️⃣",
	8: "8️⃣",
	9: "9️⃣",
}
