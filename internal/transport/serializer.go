package transport

import (
	"encoding/json"
	"log"
)

type HubBroadcaster interface {
	Broadcast(message []byte)
	BroadcastExcept(message []byte, exceptUserID string)
	SendToUser(message []byte, userID string)
}

type JSONSerializer struct {
	in  chan Event
	hub HubBroadcaster
}

func NewJSONSerializer(hub HubBroadcaster) *JSONSerializer {
	s := &JSONSerializer{
		in:  make(chan Event, 1024),
		hub: hub,
	}

	go s.readPump()

	return s
}

// Marshal serializes the given interface to json
func (s *JSONSerializer) Marshal(v interface{}) ([]byte, error) {
	log.Printf("serializing: %+v", v)
	data, err := json.Marshal(v)
	log.Println("marshalled:", string(data))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *JSONSerializer) In() chan Event {
	return s.in
}

func (s *JSONSerializer) readPump() {
	for event := range s.in {
		payload, err := s.Marshal(event)
		if err != nil {
			log.Println("marshal error:", err)
			continue
		}

		// Route based on Target
		switch {
		case event.Target.IsBroadcast():
			s.hub.Broadcast(payload)

		case event.Target.IsBroadcastExcept():
			exceptUserID := event.Target.GetExceptUserID()
			s.hub.BroadcastExcept(payload, exceptUserID)

		case event.Target.IsTargetUser():
			userID := event.Target.GetUserID()
			s.hub.SendToUser(payload, userID)

		default:
			log.Printf("unknown target type: %s", event.Target)
		}
	}
}
