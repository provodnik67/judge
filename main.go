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
	response, err := api.AskDeepSeek(
		"Кому на Руси жить хорошо?",
		"соевый российский эмигрант",
	)
	if err != nil {
		log.Fatal("DeepSeek error:", err)
	}

	log.Println("DeepSeek ответил:", response)

}
