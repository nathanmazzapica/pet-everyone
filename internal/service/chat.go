package service

import (
	"log"
)

type ChatService struct {
	events chan<- Event
}

func NewChatService(events chan<- Event) *ChatService {
	return &ChatService{events: events}
}

func moderate(msg string) string {
	// TODO: moderate message
	return msg
}

func (c *ChatService) Send(msg string) {
	msg = moderate(msg)
	log.Println("Sending message:", msg)
	type msgData struct {
		Msg    string `json:"msg"`
		Author string `json:"author"`
	}
	c.events <- Event{Type: "chat", Data: msgData{Msg: msg, Author: "User"}}
}
