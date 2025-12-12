package transport

import (
	"encoding/json"
	"log"
	"pet-everyone/internal/service"
)

// router takes in commands and routes them to the appropriate service

type Envelope struct {
	Sender string
	Data   []byte
}
type Command struct {
	Type string `json:"type"`

	// AuthorID string unused for now

	Data interface{} `json:"data"`
}

type Router struct {
	in          chan Envelope
	petService  *service.PetService
	chatService *service.ChatService
}

func NewRouter(petService *service.PetService, chatService *service.ChatService) *Router {
	return &Router{
		in:          make(chan Envelope, 2048),
		petService:  petService,
		chatService: chatService,
	}
}

func (r *Router) Route() {
	for {
		env, ok := <-r.in
		dat := env.Data
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
			user := env.Sender
			r.petService.IncrementPetCount(user)
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

func (r *Router) In() chan Envelope {
	return r.in
}
