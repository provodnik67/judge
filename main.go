package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/provodnik67/judge/api"
	"github.com/provodnik67/judge/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("TELEGRAM_BOT_TOKEN"))
	db := database.InitDB()
	defer db.Close()
	judges, err := database.GetAllJudges(db)
	if err != nil {
		log.Fatal(err)
	}
	results := make(chan string, len(judges))
	for _, judge := range judges {

		go func(judge database.Judge) {
			response, err := api.AskDeepSeek("Кому на Руси жить хорошо? Только ответь кратко. Уложись в 5 предложений.", fmt.Sprintf("Судья %s, мировоззрение: %s, как я действую: %s, история: %s", judge.Name, judge.Worldview, judge.Personality, judge.Backstory))
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
		log.Println("DeepSeek ответил:", response)
	}

}
