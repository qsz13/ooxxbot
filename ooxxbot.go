package main

import (
	"flag"
	"fmt"
	"github.com/qsz13/ooxxbot/db"
	dp "github.com/qsz13/ooxxbot/dispatcher"
	"github.com/qsz13/ooxxbot/logger"
	"github.com/qsz13/ooxxbot/spider"
	"github.com/qsz13/ooxxbot/tbot"
)

func main() {
	fmt.Println("Welcome to OOXXBot")
	parseArg()

	db := db.NewDB()
	dispatcher := dp.NewDispatcher(db)
	bot := tbot.NewBot(dispatcher, nil)
	spider := spider.NewSpider(dispatcher, db)
	spider.Start()
	bot.Start()

}

func parseArg() {
	flag.BoolVar(&logger.DebugFlag, "debug", false, "debug output")
	flag.Parse()
}
