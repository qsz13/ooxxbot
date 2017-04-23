package tbot

import (
	"database/sql"
	"fmt"
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"os"
	"strconv"
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

func (spider *Spider) getInterval() int {
	interval, err := strconv.Atoi(os.Getenv("SPIDER_INTERVAL"))
	if err != nil {
		interval = 600
	}
	return interval
}

func (spider *Spider) Start() {

	go spider.topSpider()
	go spider.apiSpider()

}

func (spider *Spider) topSpider() {
	firstTime := true
	for {
		logger.Debug("Jandan Spider is working!")
		tops, err := jd.GetTop()
		if err != nil {
			logger.Error("Jandan Spider get top failed: " + err.Error())
		} else {
			spider.filterTop(&tops)
			if len(tops) > 0 {
				if !firstTime {
					spider.bot.sendTop(tops)
				}
				spider.bot.saveSentTops(tops)
			} else {
				logger.Debug("Jandan Spider got nothing new.")
			}
			firstTime = false

		}
		time.Sleep(time.Duration(spider.getInterval()) * time.Second)
	}
}

func (spider *Spider) apiSpider() {
	for {
		logger.Debug("API Spider is working!")
		comments, err := jd.GetAllComment()
		if err != nil {
			logger.Error("API Spider failed to get comment: " + err.Error())
		} else {
			spider.saveCommentsToDB(comments)
		}
		time.Sleep(time.Duration(spider.getInterval()) * time.Second)
	}

}

func (spider *Spider) filterTop(tops *[]jd.Comment) {
	logger.Debug("Filtering top...")
	newTops := []jd.Comment{}
	for _, top := range *tops {
		if !spider.topExists(&top) {
			newTops = append(newTops, top)
		}

	}
	*tops = newTops
}

func (bot *Bot) sendTop(tops []jd.Comment) {
	ooxxSuber, _ := bot.getOOXXSubscriber()
	picSuber, _ := bot.getPicSubscriber()
	go bot.sendOOXXSubscription(ooxxSuber, tops)
	go bot.sendPicSubscription(picSuber, tops)
	logger.Debug("Sending Tops: " + fmt.Sprintf("%v", tops))

}

func (bot *Bot) sendOOXXSubscription(suber []int, tops []jd.Comment) {
	for _, u := range suber {
		for _, h := range tops {
			if h.Type == jd.OOXX_TYPE {
				bot.sendTopMessage(u, h)
			}
		}
	}

}

func (bot *Bot) sendPicSubscription(suber []int, tops []jd.Comment) {
	for _, u := range suber {
		for _, h := range tops {
			if h.Type == jd.PIC_TYPE {
				bot.sendTopMessage(u, h)
			}
		}
	}
}
