package tbot

import (
"fmt"
)

type Bot struct {
	Token string
}

func (bot *Bot) getMethodURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.Token, method)
}


func (bot *Bot) getMe() string{
	
}