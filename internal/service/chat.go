package service

import (
	"log"
)

type ChatService struct {
}

func NewChatService() *ChatService {
	return &ChatService{}
}

func moderate(msg string) string {
	// TODO: moderate message
	return msg
}

type ChatMessage struct {
	Msg    string
	Author string
}

func (c *ChatService) ProcessMessage(msg string, userID string) (ChatMessage, error) {
	msg = moderate(msg)
	log.Println("Processing message:", msg)
	return ChatMessage{
		Msg:    msg,
		Author: userID,
	}, nil
}
