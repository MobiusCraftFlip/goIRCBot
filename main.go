package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"scoutdev.org/m/v2/goIrcBot/src/bot"
	"scoutdev.org/m/v2/goIrcBot/src/db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <bot-id>\n", os.Args[0])
		os.Exit(1)
	}

	botID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid bot ID: %s\n", os.Args[1])
		os.Exit(1)
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	cfg, err := db.Load(database, botID)
	if err != nil {
		log.Fatalf("Failed to load config for bot %d: %v", botID, err)
	}

	b := bot.New(cfg, database)
	if err := b.Run(); err != nil {
		log.Fatalf("Bot exited with error: %v", err)
	}
}
