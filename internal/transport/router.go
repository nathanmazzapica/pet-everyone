package transport

import (
	"encoding/json"
	"log"
	"pet-everyone/internal/service"
)

// router takes in commands and routes them to the appropriate service

type Command struct {
	Type string `json:"type"`

	// AuthorID string unused for now

	Data interface{} `json:"data"`
}

type Router struct {
	in          chan []byte
	petService  *service.PetService
	chatService *service.ChatService
}

func NewRouter(petService *service.PetService, chatService *service.ChatService) *Router {
	return &Router{
		in:          make(chan []byte, 1024),
		petService:  petService,
		chatService: chatService,
	}
}

func (r *Router) Route() {
	for {
		dat, ok := <-r.in
		if !ok {
			// Channel is closed, break loop
			return
		}

		cmd := Command{}
		err := json.Unmarshal(dat, &cmd)
		if err != nil {
			log.Println("error unmarshalling command:", err)
			continue
		}

		switch cmd.Type {
		case "pet":
			r.petService.IncrementPetCount()
		case "petcount":
			r.petService.BroadcastPetCount()
		case "chat":
			msgData, ok := cmd.Data.(map[string]interface{})
			if !ok {
				log.Println("error: chat command Data is not a map[string]interface{}", cmd.Data)
				continue
			}

			msg, ok := msgData["msg"].(string)
			if !ok {
				log.Println("error: chat command Data is missing 'msg' field", cmd.Data)
				continue
			}

			r.chatService.Send(msg)

		default:
			log.Println("unknown command:", cmd)
		}
	}
}

func (r *Router) In() chan []byte {
	return r.in
}
