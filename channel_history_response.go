package main

type ChannelHistoryResponse struct {
	HasMore  bool      `json:"has_more"`
	Latest   string    `json:"latest"`
	Messages []Message `json:"messages"`
	Ok       bool      `json:"ok"`
	Error    string    `json:"error"`
}

type Message struct {
	Attachments []Attachment `json:"attachments"`
	BotID       string       `json:"bot_id"`
	IsStarred   bool         `json:"is_starred"`
	Reactions   []Reaction   `json:"reactions"`
	Subtype     string       `json:"subtype"`
	Text        string       `json:"text"`
	Ts          string       `json:"ts"`
	Type        string       `json:"type"`
	User        string       `json:"user"`
	Username    string       `json:"username"`
}

type Reaction struct {
	Count int64    `json:"count"`
	Name  string   `json:"name"`
	Users []string `json:"users"`
}

type Attachment struct {
	Fallback string `json:"fallback"`
	ID       int64  `json:"id"`
	Text     string `json:"text"`
}
