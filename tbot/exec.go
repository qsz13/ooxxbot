package tbot

import (
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"github.com/qsz13/ooxxbot/model"
	"io/ioutil"
	"net/http"
)

func (bot *TBot) getHelp(message *Message) {
	help := "/ip to check IP\n\n/ooxx to get random ooxx\n/pic to get random pics\n\n/looxx to get latest ooxx\n /lpic to get latest pic\n\n/sooxx to subscribe ooxx\n/spic to subscribe pic\n\n/uooxx to unsubscribe ooxx\n/upic to unsubscribe pic"

	bot.ReplyText(message.Chat.ID, help)
}

func (bot *TBot) getShortHelp(message *Message) {
	help := "/ip to check IP \n\n/ooxx to get random ooxx\n\n/pic to get random pic"
	bot.ReplyText(message.Chat.ID, help)
}

func (bot *TBot) getIP(message *Message) {
	res, err := http.Get("http://wtfismyip.com/text")
	if err != nil {
		logger.Error("Get IP failed: " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Read response body failed:" + err.Error())
		bot.ReplyError(message, err)

		return
	}
	ip := string(body)
	m, err := bot.ReplyText(message.Chat.ID, ip)
	if err != nil {
		logger.Error("Reply IP failed:" + err.Error())
		bot.ReplyError(message, err)
		return
	}
	logger.Debug("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ", " + m.Text)

}

func (bot *TBot) getRandomOOXX(message *Message) {

	content, err := bot.dispatcher.GetRandomJandan(jd.OOXX_TYPE)
	if err != nil {
		logger.Error("Get random comment failed: " + err.Error())
		bot.ReplyError(message, err)
	} else {
		content = "[OOXX]\n" + content
		err := bot.ReplyHTML(message.Chat.ID, content)
		if err != nil {
			logger.Error("Reply random OOXX failed:" + err.Error())
			bot.ReplyError(message, err)
			return
		}
		logger.Debug("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ", " + content)
	}
}

func (bot *TBot) getRandomPic(message *Message) {
	content, err := bot.dispatcher.GetRandomJandan(jd.PIC_TYPE)
	if err != nil {
		logger.Error("Get random Pic failed: " + err.Error())
		bot.ReplyError(message, err)
	} else {
		content = "[Pic]\n" + content
		err := bot.ReplyHTML(message.Chat.ID, content)
		if err != nil {
			logger.Error("Reply random Pic failed:" + err.Error())
			bot.ReplyError(message, err)
			return
		}
		logger.Debug("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ", " + content)

	}

}

func (bot *TBot) getLatestOOXX(message *Message) {
	html, err := jd.GetLatestOOXX()
	content := html.Content
	if err != nil {
		logger.Error("Get latest OOXX failed: " + err.Error())
		bot.ReplyError(message, err)
	} else {
		content = "[OOXX]\n" + content
		err := bot.ReplyHTML(message.Chat.ID, content)
		if err != nil {
			logger.Error("Reply Latest OOXX failed:" + err.Error())
			bot.ReplyError(message, err)

			return
		}
		logger.Debug("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ": " + content)
	}
}

func (bot *TBot) getLatestPic(message *Message) {
	html, err := jd.GetLatestPic()
	content := html.Content
	if err != nil {
		logger.Error("Get latest Pic failed: " + err.Error())
		bot.ReplyError(message, err)
	} else {
		content = "[Pic]\n" + content
		err := bot.ReplyHTML(message.Chat.ID, content)
		if err != nil {
			logger.Error("Reply Latest Pic failed:" + err.Error())
			bot.ReplyError(message, err)
			return
		}
		logger.Debug("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ": " + content)
	}
}

func (bot *TBot) subscribePic(message *Message) {
	err := bot.dispatcher.SubscribeJandanPic(bot.getUser(message))
	if err != nil {
		logger.Error("Subscribe Pic failed for " + message.From.FirstName + " " + message.From.LastName + ": " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *TBot) subscribeOOXX(message *Message) {
	err := bot.dispatcher.SubscribeJandanOOXX(bot.getUser(message))
	if err != nil {
		logger.Error("Subscribe OOXX failed for " + message.From.FirstName + " " + message.From.LastName + ": " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *TBot) unsubscribePic(message *Message) {
	err := bot.dispatcher.UnsubscribeJandanPic(bot.getUser(message))
	if err != nil {
		logger.Error("Unubscribe Pic failed for " + message.From.FirstName + " " + message.From.LastName + ": " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")

}

func (bot *TBot) unsubscribeOOXX(message *Message) {
	err := bot.dispatcher.UnsubscribeJandanOOXX(bot.getUser(message))
	if err != nil {
		logger.Error("Unubscribe OOXX failed for " + message.From.FirstName + " " + message.From.LastName + ": " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *TBot) getUser(message *Message) *model.User {
	if message.Chat.Type == "private" {
		return message.From
	} else if message.Chat.Type == "group" {
		chat := message.Chat
		return &model.User{ID: chat.ID, FirstName: "", LastName: "", Username: chat.Title}
	}
	return message.From
}
