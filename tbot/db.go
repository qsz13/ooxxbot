package tbot

import (
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"strconv"
)

func (bot *Bot) registerUser(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO \"user\"(id, first_name, last_name, user_name) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET first_name=excluded.first_name,last_name=excluded.last_name,user_name=excluded.user_name;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()
	if message.Chat.Type == "private" {
		user := message.From
		_, err = stmt.Exec(user.ID, user.FirstName, user.LastName, user.Username)

	} else if message.Chat.Type == "group" {
		chat := message.Chat
		_, err = stmt.Exec(chat.ID, "", "", chat.Title)

	}
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (bot *Bot) subscribeOOXXInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO subscription(\"user\", ooxx) VALUES ( $1, $2) ON CONFLICT (\"user\") DO UPDATE SET \"user\"=excluded.\"user\",ooxx=excluded.ooxx;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Chat.ID, 1)

	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil

}

func (bot *Bot) subscribePicInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO subscription(\"user\", pic) VALUES ( $1, $2) ON CONFLICT (\"user\") DO UPDATE SET \"user\"=excluded.\"user\",pic=excluded.pic;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Chat.ID, 1)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (bot *Bot) unsubscribeOOXXInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO subscription(\"user\", ooxx) VALUES ( $1, $2) ON CONFLICT (\"user\") DO UPDATE SET \"user\"=excluded.\"user\",ooxx=excluded.ooxx;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Chat.ID, 0)

	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil

}

func (bot *Bot) unsubscribePicInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO subscription(\"user\", pic) VALUES ( $1, $2) ON CONFLICT (\"user\") DO UPDATE SET \"user\"=excluded.\"user\", pic=excluded.pic;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Chat.ID, 0)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (spider *Spider) topExists(top *jd.Comment) bool {
	stmt, err := spider.db.Prepare("SELECT top from jandan where id = $1")
	if err != nil {
		logger.Error(err.Error())
		return true
	}

	var isTop bool
	err = stmt.QueryRow(top.ID).Scan(&top)
	if err != nil {
		return false
	}
	if isTop {
		return true
	}
	return false

}

func (bot *Bot) saveSentTops(tops []jd.Comment) {
	sqlStr := "INSERT INTO jandan (id, content, category, oo, xx, top) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO UPDATE SET top=true;"

	stmt, err := bot.db.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	for _, comment := range tops {
		_, err = stmt.Exec(comment.ID, comment.Content, strconv.Itoa(int(comment.Type)), strconv.Itoa(int(comment.OO)), strconv.Itoa(int(comment.XX)), false)
		if err != nil {
			logger.Error(err.Error())
		}

	}

	//format all vals at once
	logger.Debug("Top Saved!")

}

func (bot *Bot) getPicSubscriber() ([]int, error) {
	subscribers := []int{}
	rows, err := bot.db.Query("SELECT \"user\" FROM subscription where pic=true;")
	if err != nil {
		logger.Error(err.Error())
		return subscribers, err
	}

	defer rows.Close()
	sid := 0
	for rows.Next() {
		err = rows.Scan(&sid)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		subscribers = append(subscribers, sid)
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return subscribers, err
}

func (bot *Bot) getOOXXSubscriber() ([]int, error) {
	subscribers := []int{}
	rows, err := bot.db.Query("SELECT \"user\" FROM subscription where ooxx=true;")
	if err != nil {
		logger.Error(err.Error())
		return subscribers, err
	}

	defer rows.Close()
	sid := 0
	for rows.Next() {
		err = rows.Scan(&sid)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		subscribers = append(subscribers, sid)
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return subscribers, err
}

func (spider *Spider) saveCommentsToDB(comments []jd.Comment) {
	sqlStr := "INSERT INTO jandan (id, content, category, oo, xx, top) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO UPDATE SET content=excluded.content,category=excluded.category,oo=excluded.oo, xx=excluded.xx;"
	stmt, err := spider.db.Prepare(sqlStr)
	defer stmt.Close()

	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, comment := range comments {
		_, err = stmt.Exec(comment.ID, comment.Content, strconv.Itoa(int(comment.Type)), strconv.Itoa(int(comment.OO)), strconv.Itoa(int(comment.XX)), false)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	//format all vals at once
	logger.Debug("Comment Saved!")

}

func (bot *Bot) getRandomComment(jdType jd.JandanType) (string, error) {
	content := ""
	stmt, err := bot.db.Prepare("SELECT content FROM jandan WHERE category=$1 AND oo > 2*xx ORDER BY RANDOM() LIMIT 1;")
	if err != nil {
		logger.Error(err.Error())
		return content, err
	}
	err = stmt.QueryRow(jdType).Scan(&content)
	if err != nil {
		logger.Error(err.Error())
	}
	return content, err
}
