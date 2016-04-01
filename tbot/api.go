package tbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

func (bot *Bot) GetMe() (*User, error) {
	params := make(map[string]string)
	body, err := bot.sendRequest("getMe", params)

	if err != nil {
		fmt.Println(err)
	}
	var bs BotStatus

	bot.parseResult(body, &bs)
	if bs.Ok {
		return bs.Result, nil
	} else {
		return nil, fmt.Errorf("bot status is not OK, reason: " + bs.Description)
	}
}

func (bot *Bot) GetUpdates(offset, limit, timeout int) ([]Update, error) {
	params := make(map[string]string)
	params["offset"] = strconv.Itoa(offset)
	if limit != 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if timeout != 0 {
		params["timeout"] = strconv.Itoa(timeout)
	}

	body, err := bot.sendRequest("getUpdates", params)
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

func (bot *Bot) getMethodURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.Token, method)
}

func (bot *Bot) sendRequest(method string, params map[string]string) ([]byte, error) {
	url := bot.getMethodURL(method)
	if len(params) > 0 {
		url += "?"
	}
	for k, v := range params {
		url = fmt.Sprintf("%s%s=%s&", url, k, v)
	}
	fmt.Println("sending request to " + url)
	res, err := bot.client.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("request failed")
	}
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
