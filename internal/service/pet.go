package service

import (
	"context"
	"log"
	"sync"
	"time"
)

type PetDatabase interface {
	GetPetCount(id *string) (int, error)
	UpdatePetCount(petID string, userID string, count uint64) error
}

type PetService struct {
	petID          string
	petCount       uint64
	dbQueue        map[string]uint64
	pendingUpdates map[string]uint64
	mu             *sync.RWMutex
	db             PetDatabase
}

func NewPetService(ctx context.Context, petID string, model PetDatabase) *PetService {
	service := &PetService{
		petID:    petID,
		petCount: 0,
		dbQueue:  make(map[string]uint64),
		mu:       &sync.RWMutex{},
		db:       model,
	}
	service.init(ctx)
	return service
}

func (s *PetService) IncrementPetCount(userID string) error {
	s.mu.Lock()
	s.petCount++
	s.dbQueue[userID]++
	s.mu.Unlock()
	return nil
}

func (s *PetService) BroadcastPetCount() uint64 {
	return s.GetPetCount()
}

func (s *PetService) GetPetCount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.petCount
}

func (s *PetService) init(ctx context.Context) {
	c, err := s.db.GetPetCount(&s.petID)
	if err != nil {
		// TODO: pet service failed to init, figure out how to handle
		log.Fatal(err)
	}
	s.petCount = uint64(c)
	go s.persistCountsToDatabase(ctx)
}

func (s *PetService) persistCountsToDatabase(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.flushBuffer()
		case <-ctx.Done():
			log.Printf("[PET SERVICE %s]: shutting down", s.petID)
			s.flushBuffer()
			return
		}
	}
}

func (s *PetService) flushBuffer() {
	s.mu.Lock()
	snapshot := s.dbQueue
	s.dbQueue = make(map[string]uint64)
	s.mu.Unlock()

	if s.pendingUpdates == nil {
		s.pendingUpdates = make(map[string]uint64)
	}

	const MAX_RETRY_QUEUE_SIZE = 5000

	for userID, count := range snapshot {
		_, exists := s.pendingUpdates[userID]

		if exists {
			s.pendingUpdates[userID] += count
		} else {
			if len(s.pendingUpdates) >= MAX_RETRY_QUEUE_SIZE {
				log.Printf("[PET SERVICE %s]: QUEUE FULL flushing %d clicks for user %s\n", s.petID, count, userID)
				continue
			}
			s.pendingUpdates[userID] = count
		}
	}

	for userID, count := range s.pendingUpdates {
		err := s.db.UpdatePetCount(s.petID, userID, count)
		if err == nil {
			delete(s.pendingUpdates, userID)
		} else {
			log.Printf("[PET SERVICE %s]: ERROR flushing %d clicks for user %s: %s. Retrying next round\n", s.petID, count, userID, err)
		}
	}
}

func (s *PetService) Is() {}
