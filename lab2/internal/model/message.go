package model

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type CreateMessageInput struct {
	Text string `json:"text"`
}
