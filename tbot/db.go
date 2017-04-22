package tbot

import (
	_ "github.com/go-sql-driver/mysql"
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"strconv"
)

func (bot *Bot) registerUser(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.user(id, first_name, last_name, user_name) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE first_name=VALUES(first_name),last_name=VALUES(last_name),user_name=VALUES(user_name);")
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
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, ooxx) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),ooxx=VALUES(ooxx);")
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
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, pic) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),pic=VALUES(pic);")
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
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, ooxx) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),ooxx=VALUES(ooxx);")
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
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, pic) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),pic=VALUES(pic);")
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

func (bot *Bot) hotExists(hot *jd.Hot) bool {
	stmt, err := bot.db.Prepare("SELECT count(*) from ooxxbot.hot where url = ?")
	if err != nil {
		logger.Error(err.Error())
		return true
	}
	var count int
	err = stmt.QueryRow(hot.URL).Scan(&count)
	if err != nil {
		logger.Error(err.Error())
		return true
	}
	if count > 0 {
		return true
	}
	return false

}

func (bot *Bot) saveSent(hots []jd.Hot) {
	sqlStr := "INSERT INTO ooxxbot.hot(url, content, type) VALUES "
	vals := []interface{}{}

	for _, row := range hots {
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row.URL, row.Content, strconv.Itoa(int(row.Type)))
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := bot.db.Prepare(sqlStr)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer stmt.Close()

	//format all vals at once
	_, err = stmt.Exec(vals...)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Debug("Hot Saved!")

}

func (bot *Bot) getPicSubscriber() ([]int, error) {
	subscribers := []int{}
	rows, err := bot.db.Query("SELECT user FROM ooxxbot.subscription where pic=1;")
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
	rows, err := bot.db.Query("SELECT user FROM ooxxbot.subscription where ooxx=1;")
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
	sqlStr := "INSERT INTO jandan (id, content, category, oo, xx) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (id) DO UPDATE SET content=excluded.content,category=excluded.category,oo=excluded.oo, xx=excluded.xx;"
	stmt, err := spider.db.Prepare(sqlStr)
	defer stmt.Close()

	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, comment := range comments {
		sqlStr += " (unnest(?),unnest(?),unnest(?),unnest(?),unnest(?)),"
		_, err = stmt.Exec(comment.ID, comment.Content, strconv.Itoa(int(comment.Type)), strconv.Itoa(int(comment.OO)), strconv.Itoa(int(comment.XX)))
		if err != nil {
			logger.Error(err.Error())
		}
	}
	//format all vals at once
	logger.Debug("Comment Saved!")

}

func (bot *Bot) getRandomComment(jdType jd.JandanType) (string, error) {
	content := ""
	stmt, err := bot.db.Prepare("SELECT content FROM ooxxbot.jandan WHERE category = ? AND oo*2 > xx ORDER BY RAND() LIMIT 1;")
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
