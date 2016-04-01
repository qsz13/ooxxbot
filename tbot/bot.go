package tbot

import (
	"fmt"
	jd "github.com/qsz13/ooxxbot/jandan"
	rc "github.com/qsz13/ooxxbot/requestclient"
	"io/ioutil"
	"net/http"
	"strings"
)

type Bot struct {
	Token    string
	client   *http.Client
	Messages chan *Message
	Queries  chan *InlineQuery
}

func NewBot(token string, clientProxy *rc.ClientProxy) *Bot {
	bot := new(Bot)
	bot.Token = token
	bot.client, _ = rc.GetClient(clientProxy)
	return bot
}

func (bot *Bot) Start() {
	bot.loop(bot.Messages, bot.Queries)
}

func (bot *Bot) loop(messages chan *Message, queries chan *InlineQuery) {
	lastUpdate := 0
	for {

		updates, err := bot.getUpdates(lastUpdate, 0, 1000) //TBD
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
		lastUpdate = maxid + 1
	}
}

func (bot *Bot) ReplyText(ChatID int, Text string) {

	bot.sendMessage(ChatID, Text, "", false, false, -1)

}

func (bot *Bot) ReplyHTML(ChatID int, html string) {
	bot.sendMessage(ChatID, html, "HTML", false, false, -1)

}

func (bot *Bot) ExecCmd(message *Message) {
	cmd := strings.ToLower(message.Text)
	switch cmd {
	case "/start":
		go bot.getHelp(message)
	case "/ip":
		go bot.getIP(message)
		break
	case "/ooxx":
		go bot.getOOXX(message)
	default:
		bot.ReplyText(message.Chat.ID, "Incorrect, idiot!")
	}

}

func (bot *Bot) getHelp(message *Message) {
	help := "/ip to check IP\n/ooxx to get latest ooxx"

	bot.ReplyText(message.Chat.ID, help)
}

func (bot *Bot) getIP(message *Message) {
	res, err := http.Get("http://wtfismyip.com/text")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	ip := string(body)
	bot.ReplyText(message.Chat.ID, ip)

}

func (bot *Bot) getOOXX(message *Message) {
	html := jd.GetLatestOOXX().Content
	bot.ReplyText(message.Chat.ID, html)
	fmt.Println(html)
}
