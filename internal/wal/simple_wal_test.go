package wal

import (
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestNewSimpleWAL(t *testing.T) {
	WAL, err := NewPetCountWAL("test-pet")
	if err != nil {
		t.Error(err)
	}
	os.Remove(WAL.filename)
}

func TestSimpleWAL_WriteEntry(t *testing.T) {
	WAL, _ := NewPetCountWAL("test-pet")

	entries := dummyData()

	for id, count := range entries {
		err := WAL.WriteEntry(id, count)
		if err != nil {
			t.Error(err)
		}
	}
	os.Remove(WAL.filename)
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
	WAL, _ := NewPetCountWAL("test-pet")

	for id, count := range dummyData() {
		WAL.WriteEntry(id, count)
	}

	_, err := WAL.Recover()
	if err != nil {
		t.Error(err)
	}
	os.Remove(WAL.filename)
	WAL.Close()
}
