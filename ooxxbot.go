package main

import (
	"flag"
	"fmt"
	"github.com/qsz13/ooxxbot/logger"
	"github.com/qsz13/ooxxbot/tbot"
	"os"
)

func main() {
	fmt.Println("Welcome to OOXXBot")
	parseArg()

	DB_DSN := os.Getenv("DATABASE_URL")
	TOKEN := os.Getenv("TOKEN")

	bot := tbot.NewBot(TOKEN, nil, DB_DSN)

	bot.Messages = make(chan *tbot.Message, 1000)
	go handleMessages(bot)

	bot.Start()
}

func handleMessages(bot *tbot.Bot) {
	for message := range bot.Messages {
		logger.Debug("Message from " + message.From.FirstName + " " + message.From.LastName + ": " + message.Text)
		bot.ExecCmd(message)
	}
}

func parseArg() {
	flag.BoolVar(&logger.DebugFlag, "debug", false, "debug output")
	flag.Parse()

}
