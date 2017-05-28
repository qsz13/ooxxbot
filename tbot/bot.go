package tbot

import (
	dp "github.com/qsz13/ooxxbot/dispatcher"
	"github.com/qsz13/ooxxbot/logger"
	rc "github.com/qsz13/ooxxbot/requestclient"
	"net/http"
	"strconv"
	"strings"
)

type TBot struct {
	token      string
	dispatcher *dp.Dispatcher
	client     *http.Client
	Messages   chan *Message
	Queries    chan *InlineQuery
}

func NewBot(token string, dispatcher *dp.Dispatcher, clientProxy *rc.ClientProxy) *TBot {
	bot := new(TBot)
	bot.dispatcher = dispatcher
	bot.dispatcher.Bot = bot
	bot.token = token
	bot.client, _ = rc.GetClient(clientProxy)
	bot.Messages = make(chan *Message, 1000)
	return bot
}

func (bot *TBot) Start() {
	go bot.handleMessages()
	bot.loop(bot.Messages, bot.Queries)
}

func (bot *TBot) handleMessages() {
	for message := range bot.Messages {
		logger.Debug("Message from " + message.From.FirstName + " " + message.From.LastName + ": " + message.Text)
		bot.ExecCmd(message)
	}
}

func (bot *TBot) loop(messages chan *Message, queries chan *InlineQuery) {
	lastUpdate := 0
	for {

		updates, err := bot.getUpdates(lastUpdate+1, 0, 1000) //TBD
		if err != nil {
			logger.Error("Get telegram updates failed: " + err.Error())
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

func (bot *TBot) ReplyText(ChatID int, Text string) (*Message, error) {
	m, err := bot.sendMessage(ChatID, Text, "", false, false, -1)
	return m, err
}

func (bot *TBot) ReplyHTML(ChatID int, html string) error {

	texts := strings.Split(html, "\r\n")
	size := len(texts)
	var (
		err error
	)
	var messageList []string
	for i := 1; i <= size; i++ {
		text := texts[i-1]
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		messageList = append(messageList, text)
	}

	size = len(messageList)
	for i := 1; i <= size; i++ {
		text := messageList[i-1]
		if size != 1 {
			text = "(" + strconv.Itoa(i) + "/" + strconv.Itoa(size) + ")" + text
		}
		_, err = bot.sendMessage(ChatID, text, "HTML", false, false, -1)
	}
	return err
}

func (bot *TBot) ReplyHTMLWithTitle(ChatID int, title, html string) error {
	if len(title) > 0 {
		bot.sendMessage(ChatID, title, "HTML", true, false, -1)
	}

	return bot.ReplyHTML(ChatID, html)

}

func (bot *TBot) ReplyError(message *Message, err error) {
	m := "sorry, sth wrong: " + err.Error()
	bot.ReplyText(message.Chat.ID, m)
}

func (bot *TBot) ExecCmd(message *Message) {
	cmd := strings.ToLower(message.Text)
	switch cmd {
	case "/start":
		go bot.getHelp(message)
		break
	case "/s":
		go bot.getShortHelp(message)
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
		go bot.getLatestOOXX(message)
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
	case "/srss":
		bot.ReplyText(message.Chat.ID, "under construction.")
		break
	case "":
		break
	default:
		bot.ReplyText(message.Chat.ID, "Incorrect, idiot! /start")
	}

}
