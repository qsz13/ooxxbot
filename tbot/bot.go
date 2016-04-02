package tbot

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	jd "github.com/qsz13/ooxxbot/jandan"
	rc "github.com/qsz13/ooxxbot/requestclient"
	"io/ioutil"
	"net/http"
	"strings"
)

type Bot struct {
	Token    string
	client   *http.Client
	db       *sql.DB
	Messages chan *Message
	Queries  chan *InlineQuery
}

func NewBot(token string, clientProxy *rc.ClientProxy, db_dsn []string) *Bot {
	bot := new(Bot)
	bot.Token = token
	bot.client, _ = rc.GetClient(clientProxy)
	bot.db, _ = initDBConn(db_dsn)
	return bot
}

func initDBConn(db_dsn []string) (*sql.DB, error) {
	fmt.Println("Init DB connection")
	var dberr error
	for retry := 3; retry >= 0; retry-- {
		for _, d := range db_dsn {
			db, err := sql.Open("mysql", d)
			if err != nil {
				dberr = err
				continue
			}
			err = db.Ping()
			if err != nil {
				dberr = err
				continue
			}
			fmt.Println("DB connection success.")
			return db, nil
		}
	}
	fmt.Println(dberr)

	return nil, dberr
}

func (bot *Bot) Start() {
	bot.loop(bot.Messages, bot.Queries)
}

func (bot *Bot) loop(messages chan *Message, queries chan *InlineQuery) {
	lastUpdate := 0
	for {

		updates, err := bot.getUpdates(lastUpdate+1, 0, 1000) //TBD
		if err != nil {
			fmt.Println(err)
			continue
		}
		maxid := lastUpdate
		for _, update := range updates {
			if maxid < update.ID {
				maxid = update.ID
			}
			if update.Message != nil {
				messages <- update.Message
			} else if update.InlineQuery != nil {

			} else {
				continue
			}

		}
		lastUpdate = maxid
	}
}

func (bot *Bot) ReplyText(ChatID int, Text string) {

	bot.sendMessage(ChatID, Text, "", false, false, -1)

}

func (bot *Bot) ReplyHTML(ChatID int, html string) {
	bot.sendMessage(ChatID, html, "HTML", false, false, -1)

}

func (bot *Bot) ExecCmd(message *Message) {
	cmd := strings.ToLower(message.Text)
	switch cmd {
	case "/start":
		go bot.getHelp(message)
	case "/ip":
		go bot.getIP(message)
		break
	case "/ooxx":
		go bot.getOOXX(message)
	case "/pic":
		go bot.getPic(message)
	case "/sooxx":
		go bot.subscribeOOXX(message)
	case "/spic":
		go bot.subscribePic(message)
	case "/uooxx":
		go bot.unsubscribeOOXX(message)
	case "/upic":
		go bot.unsubscribePic(message)
	default:
		bot.ReplyText(message.Chat.ID, "Incorrect, idiot!")
	}

}

func (bot *Bot) getHelp(message *Message) {
	help := "/ip to check IP\n\n/ooxx to get latest ooxx\n/pic to get latest pics\n\n/sooxx to subscribe ooxx\n/spic to subscribe pic\n\n/uooxx to unsubscribe ooxx\n/upic to unsubscribe pic"

	bot.ReplyText(message.Chat.ID, help)
}

func (bot *Bot) getIP(message *Message) {
	res, err := http.Get("http://wtfismyip.com/text")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	ip := string(body)
	bot.ReplyText(message.Chat.ID, ip)

}

func (bot *Bot) getOOXX(message *Message) {
	html := jd.GetLatestOOXX().Content
	html = strings.Replace(html, "<img src", "<a href", -1)
	html = strings.Replace(html, "/>", ">查看原图</a>", -1)
	bot.ReplyHTML(message.Chat.ID, html)
	fmt.Println(html)
}

func (bot *Bot) getPic(message *Message) {
	html := jd.GetLatestPic().Content
	html = strings.Replace(html, "<img src", "<a href", -1)
	html = strings.Replace(html, "/>", ">查看原图</a>", -1)
	bot.ReplyHTML(message.Chat.ID, html)
	fmt.Println(html)

}

func (bot *Bot) subscribeOOXX(message *Message) {
	err := bot.registerUser(message.From)
	err = bot.subscribeOOXXInDB(message)
	if err != nil {
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *Bot) subscribePic(message *Message) {
	err := bot.registerUser(message.From)
	err = bot.subscribePicInDB(message)
	if err != nil {
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")

}

func (bot *Bot) unsubscribeOOXX(message *Message) {
	err := bot.registerUser(message.From)
	err = bot.unsubscribeOOXXInDB(message)
	if err != nil {
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")
}

func (bot *Bot) unsubscribePic(message *Message) {
	err := bot.registerUser(message.From)
	err = bot.unsubscribePicInDB(message)
	if err != nil {
		bot.ReplyError(message, err)
		return
	}
	bot.ReplyText(message.Chat.ID, "Success!")

}

func (bot *Bot) registerUser(user *User) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.user(id, first_name, last_name, user_name) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE first_name=VALUES(first_name),last_name=VALUES(last_name),user_name=VALUES(user_name);")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec(user.ID, user.FirstName, user.LastName, user.Username)
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

	_, err = stmt.Exec(message.From.ID, 1)
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

	_, err = stmt.Exec(message.From.ID, 1)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (bot *Bot) ReplyError(message *Message, err error) {
	m := "sorry, sth wrong: " + err.Error()
	bot.ReplyText(message.Chat.ID, m)
}

func (bot *Bot) unsubscribeOOXXInDB(message *Message) error {
	stmt, err := bot.db.Prepare("INSERT INTO ooxxbot.subscription(user, ooxx) VALUES ( ?, ?) ON DUPLICATE KEY UPDATE user=VALUES(user),ooxx=VALUES(ooxx);")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(message.From.ID, 0)
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

	_, err = stmt.Exec(message.From.ID, 0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
