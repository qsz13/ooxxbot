package tbot

import (
	"encoding/json"
	"fmt"
	rc "github.com/qsz13/ooxxbot/requestclient"
	"io/ioutil"
	"net/http"
)

type Bot struct {
	Token  string
	client *http.Client
}

func NewBot(token string, clientProxy *rc.ClientProxy) *Bot {
	bot := new(Bot)
	bot.Token = token
	bot.client, _ = rc.GetClient(clientProxy)
	return bot
}

func (bot *Bot) getMethodURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.Token, method)
}

func (bot *Bot) sendRequest(method string) ([]byte, error) {
	fmt.Println("sending request to " + bot.getMethodURL(method))
	res, _ := bot.client.Get(bot.getMethodURL(method))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body, err
}

func (bot *Bot) parseResult(body []byte, result interface{}) error {
	err := json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("tbot: invalid token")
	}
	return err

}

func (bot *Bot) GetMe() (*User, error) {
	body, err := bot.sendRequest("getMe")

	if err != nil {
		fmt.Println(err)
	}
	var bs BotStatus
	err = json.Unmarshal(body, &bs)
	if err != nil {
		fmt.Errorf("tbot: invalid token")
	}
	if bs.Ok {
		return bs.Result, nil
	} else {
		return &User{}, fmt.Errorf("bot status is not OK, reason: " + bs.Description)
	}
}

func (bot *Bot) GetUpdate() ([]Update, error) {
	body, err := bot.sendRequest("getUpdates")
	if err != nil {
		return nil, err
	}
	var ur UpdateResult

	bot.parseResult(body, &ur)
	if ur.Ok {
		return ur.Result, nil
	} else {
		return nil, fmt.Errorf("bot status is not OK, reason: " + ur.Description)
	}

}
