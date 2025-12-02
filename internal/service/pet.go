package service

import (
	"log"
	"pet-everyone/internal/db/models"
	"sync"
)

type PetService struct {
	petID    string
	petCount uint64
	mu       *sync.RWMutex
	events   chan<- Event
	db       *models.PetModel
}

func NewPetService(petID string, model *models.PetModel, out chan<- Event) *PetService {
	return &PetService{
		petID:    petID,
		petCount: 0, // will get from db later
		mu:       &sync.RWMutex{},
		events:   out,
		db:       model,
	}
}

// petEvent is used to send pet count to websocket clients
type petEvent struct {
	Ack byte `json:"c"`
}

func (s *PetService) IncrementPetCount() error {
	s.mu.Lock()
	s.petCount++
	s.mu.Unlock()
	petEvent := Event{Type: "pet", Data: petEvent{Ack: 1}}
	s.events <- petEvent
	return nil
}

func (s *PetService) BroadcastPetCount() {
	count := s.GetPetCount()
	petEvent := Event{Type: "petcount", Data: count}
	s.events <- petEvent
}

func (s *PetService) GetPetCount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.petCount
}

func (s *PetService) init() {
	c, err := s.db.GetPetCount(&s.petID)
	if err != nil {
		// TODO: pet service failed to init, figure out how to handle
		log.Fatal(err)
	}
	s.petCount = uint64(c)
}

func (s *PetService) Is() {}
