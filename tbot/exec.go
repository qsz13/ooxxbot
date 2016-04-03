package tbot

import (
	"fmt"
	jd "github.com/qsz13/ooxxbot/jandan"
	"io/ioutil"
	"net/http"
	"strings"
)

func (bot *Bot) getHelp(message *Message) {
	help := "/ip to check IP\n\n/ooxx to get latest ooxx\n/pic to get latest pics\n\n/sooxx to subscribe ooxx\n/spic to subscribe pic\n\n/uooxx to unsubscribe ooxx\n/upic to unsubscribe pic"

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
	html = strings.Replace(html, "<img src", "<a href", -1)
	html = strings.Replace(html, "/>", ">查看原图</a>", -1)
	bot.ReplyHTML(message.Chat.ID, html)
	fmt.Println(html)
}

func (bot *Bot) getPic(message *Message) {
	html := jd.GetLatestPic().Content
	html = strings.Replace(html, "<img src", "<a href", -1)
	html = strings.Replace(html, "/>", ">查看原图</a>", -1)
	bot.ReplyHTML(message.Chat.ID, html)
	fmt.Println(html)

}

func (bot *Bot) subscribeOOXX(message *Message) {
	err := bot.registerUser(message)
	err = bot.subscribeOOXXInDB(message)
	if err != nil {
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *Bot) subscribePic(message *Message) {
	err := bot.registerUser(message)
	err = bot.subscribePicInDB(message)
	if err != nil {
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")

}

func (bot *Bot) unsubscribeOOXX(message *Message) {
	err := bot.registerUser(message)
	err = bot.unsubscribeOOXXInDB(message)
	if err != nil {
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *Bot) unsubscribePic(message *Message) {
	err := bot.registerUser(message)
	err = bot.unsubscribePicInDB(message)
	if err != nil {
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")

}
