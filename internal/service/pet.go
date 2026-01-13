package service

import (
	"context"
	"errors"
	"log"
	"pet-everyone/internal/wal"
	"sync"
	"time"
)

type PetDatabase interface {
	GetPetCount(id *string) (int64, error)
	UpdatePetCountUser(petID string, userID string, count uint64) error
	UpdatePetCountGuest(petID string, guestID string, count uint64) error
}

type Actor struct {
	ID      string
	IsGuest bool
}

type actorKey struct {
	id      string
	isGuest bool
}

func keyForActor(a Actor) actorKey {
	return actorKey{id: a.ID, isGuest: a.IsGuest}
}

func encodeActorKey(k actorKey) string {
	if k.isGuest {
		return "g:" + k.id
	}
	return "u:" + k.id
}

func decodeActorKey(s string) (actorKey, error) {
	if len(s) < 3 {
		return actorKey{}, errors.New("malformed actor key")
	}
	prefix := s[:2]
	id := s[2:]
	return actorKey{id: id, isGuest: prefix == "g:"}, nil
}

type PetService struct {
	petID          string
	petCount       uint64
	dbQueue        map[actorKey]uint64
	pendingUpdates map[actorKey]uint64
	mu             *sync.RWMutex
	db             PetDatabase
	wal            *wal.PetCountWAL
}

func NewPetService(ctx context.Context, petID string, model PetDatabase) *PetService {
	wal, err := wal.NewPetCountWAL(petID)
	if err != nil {
		log.Printf("[PET SERVICE %s]: failed to create pet count WAL. err: %s", petID, err)
	}
	service := &PetService{
		petID:    petID,
		petCount: 0,
		dbQueue:  make(map[actorKey]uint64),
		mu:       &sync.RWMutex{},
		db:       model,
		wal:      wal,
	}
	service.init(ctx)
	return service
}

func (s *PetService) IncrementPetCount(actor Actor) error {
	key := keyForActor(actor)
	s.mu.Lock()
	s.petCount++
	s.dbQueue[key]++
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

func (s *PetService) shutdown() error {
	log.Printf("[PET SERVICE %s]: shutting down", s.petID)
	// save pending updates to DB
	s.flushBuffer()
	// save failed updates to WAL
	err := s.flushPendingUpdatesToWAL()
	if err != nil {
		// disk/quota is full, we can't recover from this
		log.Printf("[PET SERVICE %s]: could not flush to WAL: %s", s.petID, err)
		return err
	}

	s.wal.Close()
	return nil
}

func (s *PetService) persistCountsToDatabase(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.flushBuffer()
		case <-ctx.Done():
			log.Println(s.shutdown())
			return
		}
	}
}

func (s *PetService) flushBuffer() {
	s.mu.Lock()
	snapshot := s.dbQueue
	s.dbQueue = make(map[actorKey]uint64)
	s.mu.Unlock()

	if s.pendingUpdates == nil {
		s.pendingUpdates = make(map[actorKey]uint64)
	}

	const MAX_RETRY_QUEUE_SIZE = 5000

	for actorKey, count := range snapshot {
		_, exists := s.pendingUpdates[actorKey]

		if exists {
			s.pendingUpdates[actorKey] += count
		} else {
			if len(s.pendingUpdates) >= MAX_RETRY_QUEUE_SIZE {
				log.Printf("[PET SERVICE %s]: QUEUE FULL dropping %d clicks for actor %+v\n", s.petID, count, actorKey)
				continue
			}
			s.pendingUpdates[actorKey] = count
		}
	}

	failures := 0
	for actorKey, count := range s.pendingUpdates {
		var err error
		if actorKey.isGuest {
			err = s.db.UpdatePetCountGuest(s.petID, actorKey.id, count)
		} else {
			err = s.db.UpdatePetCountUser(s.petID, actorKey.id, count)
		}
		if err == nil {
			delete(s.pendingUpdates, actorKey)
		} else {
			log.Printf("[PET SERVICE %s]: ERROR flushing %d clicks for actor %+v: %s. Retrying next round\n", s.petID, count, actorKey, err)
			failures++
		}
	}

	log.Printf("[PET SERVICE %s]: flushed %d click updates to DB with %d errors\n", s.petID, len(snapshot), failures)
}

func (s *PetService) flushPendingUpdatesToWAL() error {
	if s.wal == nil {
		return errors.New("WAL not initialized")
	}

	for actorKey, count := range s.pendingUpdates {
		encoded := encodeActorKey(actorKey)
		err := s.wal.WriteEntry(encoded, count)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PetService) Is() {}
