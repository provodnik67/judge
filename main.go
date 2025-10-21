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
	// judge := database.Judge{
	// 	Name:        "Строгий Формалист",
	// 	Worldview:   "Законник",
	// 	Personality: "Ты строго следуешь букве закона...",
	// 	Backstory:   "Бывший прокурор...",
	// 	IsActive:    true,
	// }

	// id, err := database.CreateJudge(db, judge)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("ID: %d", id)
	judges := []string{
		"соевый российский эмигрант",
		"ватник новиоп",
	}
	results := make(chan string, len(judges))
	for _, name := range judges {

		go func(name string) {
			response, err := api.AskDeepSeek("Кому на Руси жить хорошо? Только ответь кратко. Уложись в 5 предложений.", name)
			if err != nil {
				results <- fmt.Sprintf("❌ %s: %v", name, err)
			} else {
				results <- fmt.Sprintf("✅ %s: %s", name, response)
			}
		}(name)
	}
	var responses []string
	for i := 0; i < len(judges); i++ {
		responses = append(responses, <-results) // собираем результаты
	}

	for _, response := range responses {
		log.Println("DeepSeek ответил:", response)
	}

}
