package service

// TODO: implement LRU chat history for newly joined players

import (
	"log"
)

type ChatService struct {
	resolver DisplayNameResolver
}

func NewChatService(resolver DisplayNameResolver) *ChatService {
	return &ChatService{resolver: resolver}
}

func moderate(msg string) string {
	// TODO: moderate message
	return msg
}

type ChatMessage struct {
	Msg    string `json:"msg"`
	Author string `json:"author"`
}

func (c *ChatService) ProcessMessage(msg string, actor Actor) (ChatMessage, error) {
	msg = moderate(msg)

	author, err := c.resolver.Resolve(actor)
	if err != nil {
		log.Println("failed to resolve display name:", err)
		author = "Unknown"
	}

	return ChatMessage{
		Msg:    msg,
		Author: author,
	}, nil
}
