package tbot

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	jd "github.com/qsz13/ooxxbot/jandan"
	rc "github.com/qsz13/ooxxbot/requestclient"
	"net/http"
	"strings"
	"time"
)

type Bot struct {
	Token    string
	client   *http.Client
	db       *sql.DB
	Messages chan *Message
	Queries  chan *InlineQuery
}

func NewBot(token string, clientProxy *rc.ClientProxy, db_dsn []string) *Bot {
	bot := new(Bot)
	bot.Token = token
	bot.client, _ = rc.GetClient(clientProxy)
	bot.db, _ = initDBConn(db_dsn)
	go bot.jandanSpider(1800 * time.Second)
	go bot.apiSpider(60 * time.Second)
	return bot
}

func initDBConn(db_dsn []string) (*sql.DB, error) {
	fmt.Println("Init DB connection")
	var dberr error
	for retry := 3; retry >= 0; retry-- {
		for _, d := range db_dsn {
			db, err := sql.Open("mysql", d)
			if err != nil {
				dberr = err
				continue
			}
			err = db.Ping()
			if err != nil {
				dberr = err
				continue
			}
			fmt.Println("DB connection success.")
			return db, nil
		}
	}
	fmt.Println(dberr)

	return nil, dberr
}

func (bot *Bot) Start() {
	bot.loop(bot.Messages, bot.Queries)
}

func (bot *Bot) loop(messages chan *Message, queries chan *InlineQuery) {
	lastUpdate := 0
	for {

		updates, err := bot.getUpdates(lastUpdate+1, 0, 1000) //TBD
		if err != nil {
			fmt.Println(err)
			continue
		}
		maxid := lastUpdate
		for _, update := range updates {
			if maxid < update.ID {
				maxid = update.ID
			}
			if update.Message != nil {
				messages <- update.Message
			} else if update.InlineQuery != nil {

			} else {
				continue
			}

		}
		lastUpdate = maxid
	}
}

func (bot *Bot) ReplyText(ChatID int, Text string) {

	bot.sendMessage(ChatID, Text, "", false, false, -1)

}

func (bot *Bot) ReplyHTML(ChatID int, html string) {
	bot.sendMessage(ChatID, html, "HTML", false, false, -1)

}

func (bot *Bot) ReplyError(message *Message, err error) {
	m := "sorry, sth wrong: " + err.Error()
	bot.ReplyText(message.Chat.ID, m)
}

func (bot *Bot) sendHotMessage(userid int, hot jd.Hot) {
	bot.ReplyHTML(userid, hot.Content)
}

func (bot *Bot) ExecCmd(message *Message) {
	cmd := strings.ToLower(message.Text)
	switch cmd {
	case "/start":
		go bot.getHelp(message)
		break
	case "/ip":
		go bot.getIP(message)
		break
	case "/ooxx":
		go bot.getRandomOOXX(message)
		break
	case "/pic":
		go bot.getRandomPic(message)
		break
	case "/looxx":
		go bot.getLatestPic(message)
		break
	case "/lpic":
		go bot.getLatestPic(message)
		break
	case "/sooxx":
		go bot.subscribeOOXX(message)
		break
	case "/spic":
		go bot.subscribePic(message)
		break
	case "/uooxx":
		go bot.unsubscribeOOXX(message)
		break
	case "/upic":
		go bot.unsubscribePic(message)
		break
	case "":
		break
	default:
		bot.ReplyText(message.Chat.ID, "Incorrect, idiot!")
	}

}
