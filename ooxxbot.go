package main

import (
	"fmt"
	rc "github.com/qsz13/ooxxbot/requestclient"
	"github.com/qsz13/ooxxbot/tbot"
)

func main() {
	fmt.Printf("Welcome to OOXXBot\n")
	bot := tbot.NewBot(TOKEN, &rc.ClientProxy{URL: "127.0.0.1:1080", ProxyType: rc.SOCKS5_PROXY})
	fmt.Println(bot.GetMe())
	fmt.Println(bot.GetUpdate())
}
