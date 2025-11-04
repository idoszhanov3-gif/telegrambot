package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func readStreak() int {
	data, err := os.ReadFile("streak.txt")
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(string(data))
	return n
}

func writeStreak(n int) {
	os.WriteFile("streak.txt", []byte(strconv.Itoa(n)), 0644)
}

func main() {
	bot, err := tgbotapi.NewBotAPI("ВАШ_ТОКЕН_ОТ_BOTFATHER")
	if err != nil {
		log.Panic(err)
	}

	log.Println("✅ Бот запущен и готов!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	var chatID int64
	lastActiveDate := ""

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID = update.Message.Chat.ID
		text := update.Message.Text

		switch text {
		case "/start":
			msg := "Привет, Адлан! 👋\nЯ помогу тебе учить Go и держать форму 💪\n\nКоманды:\n/plan — расписание\n/check — отметить день\n/stats — показать 🔥 серию"
			bot.Send(tgbotapi.NewMessage(chatID, msg))

		case "/plan":
			bot.Send(tgbotapi.NewMessage(chatID, "🕒 Расписание:\n19:00 — практика Go\n21:30 — спорт 💪"))

		case "/check":
			today := time.Now().Format("2006-01-02")
			if today != lastActiveDate {
				streak := readStreak() + 1
				writeStreak(streak)
				lastActiveDate = today
				msg := fmt.Sprintf("🔥 Отлично, день отмечен!\nТы держишься уже %d дней подряд!", streak)
				bot.Send(tgbotapi.NewMessage(chatID, msg))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "✅ Этот день уже отмечен. Так держать!"))
			}

		case "/stats":
			streak := readStreak()
			fire := ""
			for i := 0; i < streak; i++ {
				fire += "🔥"
			}
			if streak == 0 {
				bot.Send(tgbotapi.NewMessage(chatID, "Пока нет серии. Сегодня отличный день, чтобы начать! 💪"))
			} else {
				msg := fmt.Sprintf("Твоя серия: %d дней подряд %s", streak, fire)
				bot.Send(tgbotapi.NewMessage(chatID, msg))
			}
		}
	}
}
