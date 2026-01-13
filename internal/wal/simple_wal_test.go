package wal

import (
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
)

func newTestWAL(t *testing.T) *PetCountWAL {
	t.Helper()

	tmpDir := t.TempDir()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir to temp dir: %v", err)
	}

	wal, err := NewPetCountWAL("test-pet")
	if err != nil {
		t.Fatalf("create WAL: %v", err)
	}

	t.Cleanup(func() {
		wal.Close()
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	return wal
}

func TestNewSimpleWAL(t *testing.T) {
	WAL := newTestWAL(t)

	if WAL == nil {
		t.Fatal("expected WAL instance")
	}
}

func TestSimpleWAL_WriteEntry(t *testing.T) {
	WAL := newTestWAL(t)

	entries := dummyData()

	for id, count := range entries {
		err := WAL.WriteEntry(id, count)
		if err != nil {
			t.Error(err)
		}
	}
}

func dummyData() map[string]uint64 {
	entries := make(map[string]uint64)
	for i := 0; i < 1000; i++ {
		id := uuid.New().String()
		count := rand.Uint64()
		entries[id] = count
	}

	return entries
}

func TestSimpleWAL_Recover(t *testing.T) {
	WAL := newTestWAL(t)

	for id, count := range dummyData() {
		if err := WAL.WriteEntry(id, count); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := WAL.Recover(); err != nil {
		t.Error(err)
	}
}
