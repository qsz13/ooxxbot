package main

import (
	"fmt"
	"github.com/qsz13/ooxxbot/logger"
	"github.com/qsz13/ooxxbot/tbot"
	"os"
)

func main() {
	fmt.Printf("Welcome to OOXXBot\n")
	DB_DSN := os.Getenv("DATABASE_URL")
	TOKEN := os.Getenv("TOKEN")
	bot := tbot.NewBot(TOKEN, nil, DB_DSN)
	//bot := tbot.NewBot(TOKEN, &rc.ClientProxy{ProxyType: rc.ENV_PROXY}, DB_DSN)
	//bot := tbot.NewBot(TOKEN, &rc.ClientProxy{URL: "http://proxy.phl.sap.corp:8080", ProxyType: rc.MANUAL_PROXY}, DB_DSN)
	// fmt.Println(bot.GetMe())
	// fmt.Println(bot.GetUpdates(100, 100, 10))
	bot.Messages = make(chan *tbot.Message, 1000)

	go handleMessages(bot)

	bot.Start()
}

func handleMessages(bot *tbot.Bot) {
	for message := range bot.Messages {
		logger.Info().Println("Message from " + message.From.FirstName + " " + message.From.LastName + ": " + message.Text)
		bot.ExecCmd(message)
	}
}
