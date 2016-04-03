package tbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
)

func (bot *Bot) getMe() (*User, error) {
	params := make(map[string]string)
	body, err := bot.sendGET("getMe", params)

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

func (bot *Bot) getUpdates(offset, limit, timeout int) ([]Update, error) {
	params := make(map[string]string)
	params["offset"] = strconv.Itoa(offset)
	if limit != 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if timeout != 0 {
		params["timeout"] = strconv.Itoa(timeout)
	}
	body, err := bot.sendGET("getUpdates", params)
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

//TODO reply_markup
func (bot *Bot) sendMessage(
	chat_id int, text, parse_mode string,
	disable_web_page_preview, disable_notification bool,
	reply_to_message_id int) (*Message, error) {

	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(chat_id)
	params["text"] = text
	if parse_mode != "" {
		params["parse_mode"] = parse_mode

	}
	params["disable_web_page_preview"] = strconv.FormatBool(disable_web_page_preview)
	params["disable_notification"] = strconv.FormatBool(disable_notification)
	if reply_to_message_id > 0 {
		params["reply_to_message_id"] = strconv.Itoa(reply_to_message_id)
	}

	body, err := bot.sendPOST("sendMessage", params)
	if err != nil {
		return nil, err
	}
	var mr MessageResult
	bot.parseResult(body, &mr)
	if mr.Ok {
		return mr.Result, nil
	} else {
		fmt.Println("Message failed, reason: " + mr.Description)
		return nil, fmt.Errorf("Message failed, reason: " + mr.Description)
	}

}

func (bot *Bot) getMethodURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.Token, method)
}

func (bot *Bot) sendGET(method string, params map[string]string) ([]byte, error) {
	urladdr := bot.getMethodURL(method)
	if len(params) > 0 {
		urladdr += "?"
	}
	for k, v := range params {
		urladdr = fmt.Sprintf("%s%s=%s&", urladdr, k, v)
	}
	fmt.Println("sending GET request to " + urladdr)
	res, err := bot.client.Get(urladdr)
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
		fmt.Println("error")
		fmt.Println(string(body))
		return fmt.Errorf("tbot: invalid result")
	}
	return err
}

func (bot *Bot) sendPOST(method string, params map[string]string) ([]byte, error) {
	urladdr := bot.getMethodURL(method)
	form := make(url.Values)
	for k, v := range params {
		form.Set(k, v)
	}
	fmt.Println("sending POST request to " + urladdr)
	res, err := bot.client.PostForm(urladdr, form)
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
