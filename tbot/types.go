package tbot

import (
	"github.com/qsz13/ooxxbot/model"
)

type InlineQuery struct {
}

type ChosenInlineResult struct {
}

type BotStatus struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	ErrorCode   int    `json:"error_code"`
	Result      *model.User
}

type UpdateResult struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	ErrorCode   int    `json:"error_code"`
	Result      []Update
}

type MessageResult struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	ErrorCode   int    `json:"error_code"`
	Result      *Message
}

type Message struct {
	ID             int         `json:"message_id"`
	From           *model.User `json:"from"`
	Date           int         `json:"date"`
	Chat           *Chat       `json:"chat"`
	ForwardFrom    *model.User `json:"forward_from"`
	ForwardDate    int         `json:"forward_date"`
	ReplyToMessage *Message    `json:"reply_to_message"`
	Text           string      `json:"text"`
	Audio          *Audio      `json:"Audio"`
}

type Chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Audio struct {
}

type Update struct {
	ID                 int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
}
