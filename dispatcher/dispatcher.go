package dispatcher

import (
	"fmt"
	"github.com/qsz13/ooxxbot/db"
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"github.com/qsz13/ooxxbot/model"
)

type ITBot interface {
	ReplyHTMLWithTitle(ChatID int, title, html string) error
}

type Dispatcher struct {
	Bot ITBot
	db  *db.DB
}

func NewDispatcher(db *db.DB) *Dispatcher {
	return &Dispatcher{db: db}
}

func (dp *Dispatcher) GetRandomJandan(jdType jd.JandanType) (string, error) {
	return dp.db.GetRandomComment(jdType)
}

func (dp *Dispatcher) SubscribeJandanPic(user *model.User) error {
	err := dp.db.RegisterUser(user)
	if err != nil {
		return err
	}
	return dp.db.SubscribeJandanPic(user.ID)
}

func (dp *Dispatcher) SubscribeJandanOOXX(user *model.User) error {
	err := dp.db.RegisterUser(user)
	if err != nil {
		return err
	}
	return dp.db.SubscribeJandanOOXX(user.ID)
}

func (dp *Dispatcher) UnsubscribeJandanPic(user *model.User) error {
	err := dp.db.RegisterUser(user)
	if err != nil {
		return err
	}
	return dp.db.UnsubscribeJandanPic(user.ID)

}

func (dp *Dispatcher) UnsubscribeJandanOOXX(user *model.User) error {
	err := dp.db.RegisterUser(user)
	if err != nil {
		return err
	}
	return dp.db.UnsubscribeJandanOOXX(user.ID)

}

func (dp *Dispatcher) SendJandanTop(tops []jd.Comment) {

	ooxxSuber, _ := dp.db.GetOOXXSubscriber()
	picSuber, _ := dp.db.GetPicSubscriber()
	go dp.sendJandanOOXXSubscription(ooxxSuber, tops)
	go dp.sendJandanPicSubscription(picSuber, tops)
	logger.Debug("Sending Tops: " + fmt.Sprintf("%v", tops))
}

func (dp *Dispatcher) sendJandanOOXXSubscription(suber []int, tops []jd.Comment) {
	for _, u := range suber {
		for _, t := range tops {
			if t.Type == jd.OOXX_TYPE {
				title := ""
				content := t.Content
				if len(t.Link) != 0 {
					title = "<a href=\"" + t.Link + "\">[OOXX]</a>\n"

				} else {
					content = "[OOXX]\n" + content
				}

				logger.Debug("Sending: " + t.Content)
				dp.Bot.ReplyHTMLWithTitle(u, title, content)
			}
		}
	}

}

func (dp *Dispatcher) sendJandanPicSubscription(suber []int, tops []jd.Comment) {
	for _, u := range suber {
		for _, t := range tops {
			if t.Type == jd.PIC_TYPE {
				title := ""
				content := t.Content
				if len(t.Link) != 0 {
					title = "<a href=\"" + t.Link + "\">[Pic]</a>\n"

				} else {
					content = "[Pic]\n" + content
				}
				logger.Debug("Sending: " + t.Content)
				dp.Bot.ReplyHTMLWithTitle(u, title, content)
			}
		}
	}
}
