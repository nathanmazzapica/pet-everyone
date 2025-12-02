package transport

import (
	"encoding/json"
	"log"
	"pet-everyone/internal/service"
)

type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	Subscribe(ch chan []byte)
	In() chan interface{}
	readPump()
	writePump()
}

type JSONSerializer struct {
	in         chan service.Event
	out        chan envelope
	subscriber chan<- []byte
}

type envelope struct {
	recipient string
	data      []byte
}

func NewJSONSerializer() *JSONSerializer {
	s := &JSONSerializer{
		in:  make(chan service.Event, 1024),
		out: make(chan envelope, 256),
	}

	go s.readPump()
	go s.writePump()

	return s
}

// Subscribe adds a channel to the list of subscribers
func (s *JSONSerializer) Subscribe(ch chan []byte) {
	if s.subscriber != nil {
		log.Println("another hub attemped to subscribe, ignoring")
		return
	}
	s.subscriber = ch
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

func (s *JSONSerializer) In() chan service.Event {
	return s.in
}

func (s *JSONSerializer) readPump() {
	for {
		data, ok := <-s.in
		if !ok {
			return
		}

		payload, err := s.Marshal(data)
		envelope := envelope{recipient: "all", data: payload}
		if err != nil {
			// TODO: handle broken payload
			log.Println(err)
			return
		}
		s.out <- envelope

	}
}

// Right now this sends to all subscribers, which is not what we want.

func (s *JSONSerializer) writePump() {
	for {
		envelope, ok := <-s.out
		if !ok {
			// channel closed
			return
		}

		select {
		case s.subscriber <- envelope.data:
		default:
			// handle
		}
	}
}

// idea: envelope that stores recipient id and payload, maps to map[id]chan
// need way of propagating hub id through service and into serializer
// envelope interface? GetID() string, GetPayload() ?
// At that point serializer feels redundant, can demote to a bridging layer
