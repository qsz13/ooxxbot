package tbot

import (
	rc "github.com/qsz13/ooxxbot/requestclient"
	"net/http"
)

type Bot struct {
	Token  string
	client *http.Client
}

func NewBot(token string, clientProxy *rc.ClientProxy) *Bot {
	bot := new(Bot)
	bot.Token = token
	bot.client, _ = rc.GetClient(clientProxy)
	return bot
}
