package tbot

import (
	"fmt"
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"time"
)

func (bot *Bot) jandanSpider(interval time.Duration) {
	firstTime := true
	for {
		logger.Info().Println("Jandan Spider is working!")
		hots, err := jd.GetHot()
		if err != nil {
			logger.Error().Println("Jandan Spider get hot failed: " + err.Error())
		} else {
			bot.filterHot(&hots)
			if len(hots) > 0 {
				if !firstTime {
					bot.sendHot(hots)
				}
				bot.saveSent(hots)
			} else {
				logger.Info().Println("Jandan Spider got nothing new.")
			}
			firstTime = false

		}

		time.Sleep(interval)
	}
}

func (bot *Bot) apiSpider(interval time.Duration) {

	for {
		logger.Info().Println("API Spider is working!")
		comments, err := jd.GetAllComment()
		if err != nil {
			logger.Error().Println("API Spider failed to get comment: " + err.Error())
		} else {
			bot.saveCommentsToDB(comments)
		}
		time.Sleep(interval)
	}

}

func (bot *Bot) filterHot(hots *[]jd.Hot) {
	logger.Info().Println("Filtering hot...")
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
	logger.Info().Println("Sending Hots: " + fmt.Sprintf("%v", hots))

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
