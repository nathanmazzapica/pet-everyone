package service

import (
	"context"
	"fmt"
	"testing"
)

func TestFlushBuffer(t *testing.T) {
	mock := &SaboteurDatabase{ShouldFail: true}
	testService := NewPetService(context.Background(), "1", mock)
	for i := 0; i < 6000; i++ {
		mockID := fmt.Sprintf("%d", i)
		testService.dbQueue[actorKey{id: mockID, isGuest: false}] = 1
	}

	t.Log("Flushing buffer...")
	testService.flushBuffer()

	if len(testService.pendingUpdates) != 5000 {
		t.Errorf("Expected 5000 pending updates, got %d", len(testService.pendingUpdates))
		t.FailNow()
	}

	t.Log("Success!")
}

func TestFlushBufferAddsToPendingOnFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock := &SaboteurDatabase{ShouldFail: true}
	svc := NewPetService(ctx, "pet-1", mock)
	cancel()

	actor := actorKey{id: "user-1", isGuest: false}

	svc.mu.Lock()
	svc.pendingUpdates = map[actorKey]uint64{actor: 3}
	svc.dbQueue[actor] = 2
	svc.mu.Unlock()

	svc.flushBuffer()

	svc.mu.RLock()
	pending := svc.pendingUpdates[actor]
	svc.mu.RUnlock()

	if pending != 5 {
		t.Fatalf("expected pending updates to total 5, got %d", pending)
	}
}
