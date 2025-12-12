package service

import (
	"fmt"
	"log"
	"pet-everyone/internal/data/model"
	"sync"
	"time"
)

type PetService struct {
	petID    string
	petCount uint64
	dbQueue  map[string]uint64
	mu       *sync.RWMutex
	events   chan<- Event
	db       *model.PetModel
}

func NewPetService(petID string, model *model.PetModel, out chan<- Event) *PetService {
	service := &PetService{
		petID:    petID,
		petCount: 0, // will get from db later
		dbQueue:  make(map[string]uint64),
		mu:       &sync.RWMutex{},
		events:   out,
		db:       model,
	}
	service.init()
	return service
}

// petEvent is used to send pet count to websocket clients
type petEvent struct {
	Ack byte `json:"c"`
}

func (s *PetService) IncrementPetCount(userID string) error {
	s.mu.Lock()
	s.petCount++
	s.dbQueue[userID]++
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
	go s.persistCountsToDatabase()
}

func (s *PetService) persistCountsToDatabase() {
	for {
		time.Sleep(4 * time.Second)
		fmt.Println("persisting counts to database")
		s.mu.RLock()
		defer s.mu.RUnlock()
		for userID, count := range s.dbQueue {
			fmt.Println("persisting count for user", userID, "to database")
			err := s.db.UpdatePetCount(s.petID, userID, count)
			if err != nil {
				fmt.Println("PETID: ", s.petID)
				fmt.Println("error persisting count to database:", err)
			}
		}
	}
}

func (s *PetService) Is() {}
