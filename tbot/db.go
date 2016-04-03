package tbot

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func (bot *Bot) registerUser(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.user(id, first_name, last_name, user_name) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE first_name=VALUES(first_name),last_name=VALUES(last_name),user_name=VALUES(user_name);")
	if err != nil {
		fmt.Println(err)
		return err
	}
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
	if message.Chat.Type == "private" {
		_, err = stmt.Exec(message.From.ID, 1)

	} else if message.Chat.Type == "group" {
		_, err = stmt.Exec(message.Chat.ID, 1)

	}

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
	if message.Chat.Type == "private" {
		_, err = stmt.Exec(message.From.ID, 1)

	} else if message.Chat.Type == "group" {
		_, err = stmt.Exec(message.Chat.ID, 1)

	}

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
	if message.Chat.Type == "private" {
		_, err = stmt.Exec(message.From.ID, 1)

	} else if message.Chat.Type == "group" {
		_, err = stmt.Exec(message.Chat.ID, 1)

	}

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
	if message.Chat.Type == "private" {
		_, err = stmt.Exec(message.From.ID, 1)

	} else if message.Chat.Type == "group" {
		_, err = stmt.Exec(message.Chat.ID, 1)

	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
