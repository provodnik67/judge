package main

import (
	"fmt"
	"log"
	"os"

	"github.com/provodnik67/judge/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("TELEGRAM_BOT_TOKEN"))
	db := config.InitDB()
	defer db.Close()
	log.Println("Bot starting with database...")
}
