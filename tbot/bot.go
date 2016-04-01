package tbot

import (
	"fmt"
	rc "github.com/qsz13/ooxxbot/requestclient"
	"net/http"
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
