package tbot

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	jd "github.com/qsz13/ooxxbot/jandan"
	"strconv"
)

func (bot *Bot) registerUser(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.user(id, first_name, last_name, user_name) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE first_name=VALUES(first_name),last_name=VALUES(last_name),user_name=VALUES(user_name);")
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return err
	}
	return nil
}

func (bot *Bot) subscribeOOXXInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, ooxx) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),ooxx=VALUES(ooxx);")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Chat.ID, 1)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func (bot *Bot) subscribePicInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, pic) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),pic=VALUES(pic);")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Chat.ID, 1)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (bot *Bot) unsubscribeOOXXInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, ooxx) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),ooxx=VALUES(ooxx);")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Chat.ID, 1)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func (bot *Bot) unsubscribePicInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, pic) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),pic=VALUES(pic);")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Chat.ID, 1)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (bot *Bot) hotExists(hot *jd.Hot) bool {
	stmt, err := bot.db.Prepare("SELECT count(*) from ooxxbot.hot where url = ?")
	if err != nil {
		fmt.Println(err)
		return true
	}
	var count int
	err = stmt.QueryRow(hot.URL).Scan(&count)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	//format all vals at once
	_, err = stmt.Exec(vals...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Saved!")

}
