package tbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qsz13/ooxxbot/logger"
	"github.com/qsz13/ooxxbot/model"
	"io/ioutil"
	"net/url"
	"strconv"
)

func (bot *TBot) getMe() (*model.User, error) {
	params := make(map[string]string)
	body, err := bot.sendGET("getMe", params)

	if err != nil {
		logger.Error("Request API getMe failed: " + err.Error())
		return nil, err
	}
	var bs BotStatus

	bot.parseResult(body, &bs)
	if bs.Ok {
		return bs.Result, nil
	} else {
		err = errors.New("Bot status is not OK, reason: " + bs.Description)
		logger.Error(err.Error())
		return nil, err
	}
}

func (bot *TBot) getUpdates(offset, limit, timeout int) ([]Update, error) {
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
		logger.Debug("Request API getUpdates failed: " + err.Error())
		return nil, err
	}
	var ur UpdateResult

	bot.parseResult(body, &ur)
	if ur.Ok {
		return ur.Result, nil
	} else {
		err = errors.New("Bot status is not OK, reason: " + ur.Description)
		logger.Error(err.Error())
		return nil, err
	}

}

//TODO reply_markup
func (bot *TBot) sendMessage(
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
		logger.Error("Request API sendMessage failed: " + err.Error())
		return nil, err
	}
	var mr MessageResult
	bot.parseResult(body, &mr)
	if mr.Ok {
		return mr.Result, nil
	} else {
		err = errors.New("Message failed, reason: " + mr.Description + ". " + text)
		logger.Error(err.Error())
		return nil, err
	}

}

func (bot *TBot) getMethodURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.token, method)
}

func (bot *TBot) sendGET(method string, params map[string]string) ([]byte, error) {
	urladdr := bot.getMethodURL(method)
	if len(params) > 0 {
		urladdr += "?"
	}
	for k, v := range params {
		urladdr = fmt.Sprintf("%s%s=%s&", urladdr, k, v)
	}
	//logger.Debug("Get request to: " + urladdr)
	res, err := bot.client.Get(urladdr)
	if err != nil {
		logger.Debug("Request for " + urladdr + " failed: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Read response body failed: " + err.Error())
		return nil, err
	}
	return body, err
}

func (bot *TBot) parseResult(body []byte, result interface{}) error {
	err := json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("Parse json failed: " + err.Error() + ", Content:\n" + string(body))
	}
	return err
}

func (bot *TBot) sendPOST(method string, params map[string]string) ([]byte, error) {
	urladdr := bot.getMethodURL(method)
	form := make(url.Values)
	for k, v := range params {
		form.Set(k, v)
	}
	//logger.Debug("POST request to: " + urladdr)
	res, err := bot.client.PostForm(urladdr, form)
	if err != nil {
		logger.Error("Request for " + urladdr + " failed: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Read response body failed: " + err.Error())
	}
	return body, err
}
