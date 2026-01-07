package wal

import (
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestNewSimpleWAL(t *testing.T) {
	_, err := NewPetCountWAL("wal.log", "")
	if err != nil {
		t.Error(err)
	}
}

func TestSimpleWAL_WriteEntry(t *testing.T) {
	WAL, _ := NewPetCountWAL("wal.log", "")
	defer os.Remove("wal.log")

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
	WAL, _ := NewPetCountWAL("wal.log", "")

	for id, count := range dummyData() {
		WAL.WriteEntry(id, count)
	}

	rebuilt, err := WAL.Recover()
	if rebuilt == nil || err != nil {
		t.Error(err)
	}

	WAL.Close()
}

func TestSimpleWAL_Close(t *testing.T) {
	WAL, _ := NewPetCountWAL("wal.log", "")
	err := WAL.Close()
	if err != nil {
		t.Error(err)
	}
}
