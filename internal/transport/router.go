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
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Router struct {
	in          chan Envelope
	events      chan Event
	petService  *service.PetService
	chatService *service.ChatService
}

func NewRouter(petService *service.PetService, chatService *service.ChatService) *Router {
	return &Router{
		in:          make(chan Envelope, 2048),
		events:      make(chan Event, 1024),
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
			err := r.petService.IncrementPetCount(env.Sender)
			if err != nil {
				log.Println("error incrementing pet count:", err)
				continue
			}
			// Don't broadcast to sender, they already rendered their pet
			r.events <- Event{
				Type:   "pet",
				Data:   map[string]interface{}{"c": 1},
				Target: TargetExcept(env.Sender),
			}

		case "petcount":
			count := r.petService.BroadcastPetCount()
			// Broadcast to all - sync full count
			r.events <- Event{
				Type:   "petcount",
				Data:   count,
				Target: TargetBroadcast,
			}

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

			chatMsg, err := r.chatService.ProcessMessage(msg, env.Sender)
			if err != nil {
				log.Println("error processing chat message:", err)
				continue
			}

			// Broadcast to all except sender
			r.events <- Event{
				Type: "chat",
				Data: map[string]interface{}{
					"msg":    chatMsg.Msg,
					"author": chatMsg.Author,
				},
				Target: TargetExcept(env.Sender),
			}

		default:
			log.Println("unknown command:", cmd)
		}
	}
}

func (r *Router) In() chan Envelope {
	return r.in
}

func (r *Router) Events() <-chan Event {
	return r.events
}
