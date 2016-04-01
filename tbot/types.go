package tbot

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
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

type Message struct {
	ID             int      `json:"message_id"`
	From           *User    `json:"from"`
	Date           int      `json:"date"`
	Chat           *Chat    `json:"chat"`
	ForwardFrom    *User    `json:"forward_from"`
	ForwardDate    int      `json:"forward_date"`
	ReplyToMessage *Message `json:"reply_to_message"`
	Text           string   `json:"text"`
	Audio          *Audio   `json:"Audio"`
}

type InlineQuery struct {
}

type ChosenInlineResult struct {
}

type Update struct {
	ID                 int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
}

type ResultBase struct {
	Ok          bool
	Description string
	ErrorCode   int `json:"error_code"`
}

type BotStatus struct {
	ResultBase
	Result *User
}

type UpdateResult struct {
	ResultBase
	Result []Update
}

type MessageResult struct {
	ResultBase
	Result *Message
}
