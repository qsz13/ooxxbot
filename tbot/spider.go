package tbot

import (
	"database/sql"
	"fmt"
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"time"
)

type Spider struct {
	bot *Bot
	db  *sql.DB
}

func NewSpider(bot *Bot) *Spider {
	spider := new(Spider)
	spider.bot = bot
	spider.db = bot.db
	return spider
}

func (spider *Spider) Start() {
	//go bot.jandanSpider(1800 * time.Second)
	go spider.apiSpider(60 * time.Second)

}

/*func (spider *Spider) jandanSpider(interval time.Duration) {
	firstTime := true
	for {
		logger.Debug("Jandan Spider is working!")
		hots, err := jd.GetHot()
		if err != nil {
			logger.Error("Jandan Spider get top failed: " + err.Error())
		} else {
			spider.filterHot(&hots)
			if len(hots) > 0 {
				if !firstTime {
					spider.bot.sendHot(hots)
				}
				spider.bot.saveSent(hots)
			} else {
				logger.Debug("Jandan Spider got nothing new.")
			}
			firstTime = false

		}

		time.Sleep(interval)
	}
}*/

func (spider *Spider) apiSpider(interval time.Duration) {
	for {
		logger.Debug("API Spider is working!")
		comments, err := jd.GetAllComment()
		if err != nil {
			logger.Error("API Spider failed to get comment: " + err.Error())
		} else {
			spider.saveCommentsToDB(comments)
		}
		time.Sleep(interval)
	}

}

func (bot *Bot) filterHot(hots *[]jd.Hot) {
	logger.Debug("Filtering hot...")
	newHots := []jd.Hot{}
	for _, hot := range *hots {
		if !bot.hotExists(&hot) {
			newHots = append(newHots, hot)
		}

	}
	*hots = newHots
}

func (bot *Bot) sendHot(hots []jd.Hot) {
	ooxxSuber, _ := bot.getOOXXSubscriber()
	picSuber, _ := bot.getPicSubscriber()
	go bot.sendOOXXSubscription(ooxxSuber, hots)
	go bot.sendPicSubscription(picSuber, hots)
	logger.Debug("Sending Hots: " + fmt.Sprintf("%v", hots))

}

func (bot *Bot) sendOOXXSubscription(suber []int, hots []jd.Hot) {
	for _, u := range suber {
		for _, h := range hots {
			if h.Type == jd.OOXX_TYPE {
				bot.sendHotMessage(u, h)
			}
		}
	}

}

func (bot *Bot) sendPicSubscription(suber []int, hots []jd.Hot) {
	for _, u := range suber {
		for _, h := range hots {
			if h.Type == jd.PIC_TYPE {
				bot.sendHotMessage(u, h)
			}
		}
	}
}
