package service

import (
	"context"
	"fmt"
	"testing"
)

func TestFlushBuffer(t *testing.T) {
	mock := &SaboteurDatabase{ShouldFail: true}
	testService := NewPetService(context.Background(), "1", mock, nil)
	for i := 0; i < 6000; i++ {
		mockID := fmt.Sprintf("%d", i)
		testService.dbQueue[mockID] = 1
	}

	t.Log("Flushing buffer...")
	testService.flushBuffer()

	if len(testService.pendingUpdates) != 5000 {
		t.Errorf("Expected 5000 pending updates, got %d", len(testService.pendingUpdates))
		t.FailNow()
	}

	t.Log("Success!")
}
