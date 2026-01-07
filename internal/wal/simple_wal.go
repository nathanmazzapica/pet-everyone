package wal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type WAL interface {
	WriteEntry(userID string, count uint64) error
	Recover() (map[string]uint64, error)
	Close() error
}

type PetCountWAL struct {
	filename string
	mu       sync.Mutex
	file     *os.File
	writes   int
}

type RowErr struct {
	Err error
	Row string
}

func (e RowErr) Error() string {
	return fmt.Sprintf("Error on row: %s, %s", e.Row, e.Err.Error())
}

var (
	ErrMalformedWALEntry = fmt.Errorf("malformed WAL entry")
	ErrMalformedUUID     = fmt.Errorf("malformed UUID")
)

func NewPetCountWAL(petID string) (*PetCountWAL, error) {
	filename := fmt.Sprintf("pet_%s.wal", petID)
	w := &PetCountWAL{filename: filename}

	// Open file once on creation for speed
	f, err := os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	w.file = f

	return w, nil
}

func (w *PetCountWAL) WriteEntry(userID string, count uint64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	const FORMAT = "%s,%d\n"
	entry := fmt.Sprintf(FORMAT, userID, count)
	_, err := w.file.WriteString(entry)
	if err != nil {
		return err
	}

	w.writes++

	return nil
}

func (w *PetCountWAL) Recover() (map[string]uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	rebuiltState := make(map[string]uint64)

	f, err := os.Open(w.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewScanner(f)
	for buf.Scan() {
		line := buf.Text()

		// format: userID,count\n
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			log.Printf("discarding malformed WAL entry: %s\nerr:%s", line, ErrMalformedWALEntry)
			continue
		}

		userID := parts[0]

		// check for invalid uuid and skip them
		_, err := uuid.Parse(userID)
		if err != nil {
			log.Printf("discarding malformed WAL entry: %s\nerr:%s", line, ErrMalformedUUID)
			continue
		}

		count, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return rebuiltState, err
		}

		rebuiltState[userID] = count
	}
	return rebuiltState, nil
}

// Close closes the WAL file. Errors are logged, but the WAL file is always closed.
func (w *PetCountWAL) Close() {
	// ensure all writes are flushed to disk
	defer w.file.Close()
	err := w.file.Sync()
	if err != nil {
		log.Println("failed to sync WAL file:", err)
	}

	if w.writes == 0 {
		err = os.Remove(w.filename)
		if err != nil {
			log.Println("failed to remove empty WAL file:", err)
		}
	}
}
