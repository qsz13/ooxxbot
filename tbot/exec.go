package tbot

import (
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"io/ioutil"
	"net/http"
)

func (bot *Bot) getHelp(message *Message) {
	help := "/ip to check IP\n\n/ooxx to get random ooxx\n/pic to get random pics\n\n/looxx to get latest ooxx\n /lpic to get latest pic\n\n/sooxx to subscribe ooxx\n/spic to subscribe pic\n\n/uooxx to unsubscribe ooxx\n/upic to unsubscribe pic"

	bot.ReplyText(message.Chat.ID, help)
}

func (bot *Bot) getIP(message *Message) {
	res, err := http.Get("http://wtfismyip.com/text")
	if err != nil {
		logger.Error().Println("Get IP failed: " + err.Error())
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error().Println("Read response body failed:" + err.Error())
		return
	}
	ip := string(body)
	m, err := bot.ReplyText(message.Chat.ID, ip)
	if err != nil {
		logger.Error().Println("Reply IP failed:" + err.Error())
		return
	}
	logger.Info().Println("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ", " + m.Text)

}

func (bot *Bot) sendHotMessage(userid int, hot jd.Hot) {
	content := hot.Content
	if hot.Type == jd.OOXX_TYPE {
		content = "[OOXX]\n" + content
	} else if hot.Type == jd.PIC_TYPE {
		content = "[Pic]\n" + content
	}
	bot.ReplyHTML(userid, content)
}

func (bot *Bot) getRandomOOXX(message *Message) {

	content, err := bot.getRandomComment(jd.OOXX_TYPE)
	if err != nil {
		logger.Error().Println("Get random comment failed: " + err.Error())
	} else {
		content = "[OOXX]\n" + content
		m, err := bot.ReplyHTML(message.Chat.ID, content)
		if err != nil {
			logger.Error().Println("Reply random OOXX failed:" + err.Error())
			return
		}
		logger.Info().Println("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ", " + m.Text)
	}
}

func (bot *Bot) getRandomPic(message *Message) {
	content, err := bot.getRandomComment(jd.PIC_TYPE)
	if err != nil {
		logger.Error().Println("Get random Pic failed: " + err.Error())
	} else {
		content = "[Pic]\n" + content
		m, err := bot.ReplyHTML(message.Chat.ID, content)
		if err != nil {
			logger.Error().Println("Reply random Pic failed:" + err.Error())
			return
		}
		logger.Info().Println("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ", " + m.Text)

	}

}

func (bot *Bot) getLatestOOXX(message *Message) {
	html, err := jd.GetLatestOOXX()
	content := html.Content
	if err != nil {
		logger.Error().Println("Get latest OOXX failed: " + err.Error())
	} else {
		content = "[OOXX]\n" + content
		m, err := bot.ReplyHTML(message.Chat.ID, content)
		if err != nil {
			logger.Error().Println("Reply Latest OOXX failed:" + err.Error())
			return
		}
		logger.Info().Println("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ": " + m.Text)
	}
}

func (bot *Bot) getLatestPic(message *Message) {
	html, err := jd.GetLatestPic()
	content := html.Content
	if err != nil {
		logger.Error().Println("Get latest Pic failed: " + err.Error())
	} else {
		content = "[Pic]\n" + content
		m, err := bot.ReplyHTML(message.Chat.ID, content)
		if err != nil {
			logger.Error().Println("Reply Latest Pic failed:" + err.Error())
			return
		}
		logger.Info().Println("Message Sent to " + message.From.FirstName + " " + message.From.LastName + ": " + m.Text)
	}
}

func (bot *Bot) subscribeOOXX(message *Message) {
	err := bot.registerUser(message)
	err = bot.subscribeOOXXInDB(message)
	if err != nil {
		logger.Error().Println("Subscribe OOXX failed for " + message.From.FirstName + " " + message.From.LastName + ": " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *Bot) subscribePic(message *Message) {
	err := bot.registerUser(message)
	err = bot.subscribePicInDB(message)
	if err != nil {
		logger.Error().Println("Subscribe Pic failed for " + message.From.FirstName + " " + message.From.LastName + ": " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")

}

func (bot *Bot) unsubscribeOOXX(message *Message) {
	err := bot.registerUser(message)
	err = bot.unsubscribeOOXXInDB(message)
	if err != nil {
		logger.Error().Println("Unubscribe OOXX failed for " + message.From.FirstName + " " + message.From.LastName + ": " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *Bot) unsubscribePic(message *Message) {
	err := bot.registerUser(message)
	err = bot.unsubscribePicInDB(message)
	if err != nil {
		logger.Error().Println("Unubscribe Pic failed for " + message.From.FirstName + " " + message.From.LastName + ": " + err.Error())
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")

}
