package spider

import (
	"github.com/qsz13/ooxxbot/db"
	dp "github.com/qsz13/ooxxbot/dispatcher"
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	//	"os"
	//	"strconv"
	"time"
)

type Spider struct {
	dispatcher *dp.Dispatcher
	db         *db.DB
	interval   int
}

func NewSpider(dispatcher *dp.Dispatcher, db *db.DB, interval int) *Spider {
	return &Spider{dispatcher: dispatcher, db: db, interval: interval}
}

func (spider *Spider) getInterval() int {
	//interval, err := strconv.Atoi(os.Getenv("SPIDER_INTERVAL"))
	//if err != nil {
	//		interval = 600
	//	}

	return spider.interval
}

func (spider *Spider) Start() {

	go spider.topSpider()
	go spider.apiSpider()

}

func (spider *Spider) topSpider() {
	//firstTime := true
	for {
		logger.Debug("Jandan Spider is working!")
		tops, err := jd.GetTop()
		if err != nil {
			logger.Error("Jandan Spider get top failed: " + err.Error())
		} else {
			spider.filterTop(&tops)
			if len(tops) > 0 {
				//			if !firstTime {
				spider.dispatcher.SendJandanTop(tops)
				//			}
				spider.db.SaveJandanTops(tops)
			} else {
				logger.Debug("Jandan Spider got nothing new.")
			}
			//		firstTime = false

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
			spider.db.SaveJandanToDB(comments)
		}
		time.Sleep(time.Duration(spider.getInterval()) * time.Second)
	}

}

func (spider *Spider) filterTop(tops *[]jd.Comment) {
	logger.Debug("Filtering top...")
	newTops := []jd.Comment{}
	for _, top := range *tops {
		if !spider.db.TopExists(&top) {
			newTops = append(newTops, top)
		}

	}
	*tops = newTops
}
