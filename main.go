package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/provodnik67/judge/api"
	"github.com/provodnik67/judge/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := database.InitDB()
	defer db.Close()
	judges, err := database.GetAllJudges(db)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if strings.HasSuffix(update.Message.Text, "?") {
			results := make(chan string, len(judges))
			for _, judge := range judges {

				go func(judge database.Judge) {
					response, err := api.AskDeepSeek(update.Message.Text+" Только ответь кратко. Уложись в 1 предложение.", fmt.Sprintf("Судья %s, мировоззрение: %s, как я действую: %s, история: %s", judge.Name, judge.Worldview, judge.Personality, judge.Backstory))
					if err != nil {
						results <- fmt.Sprintf("❌ %s: %v", judge.Name, err)
					} else {
						results <- fmt.Sprintf("✅ %s: %s", judge.Name, response)
					}
				}(judge)
			}
			var responses []string
			for i := 0; i < len(judges); i++ {
				responses = append(responses, <-results) // собираем результаты
			}
			for _, response := range responses {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вердикт:"+response)
				bot.Send(msg)
			}
		}
	}
}
