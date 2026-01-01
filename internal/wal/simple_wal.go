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

type SimpleWAL struct {
	filename string
	mu       sync.Mutex
	file     *os.File
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

func NewSimpleWAL(filename string) (*SimpleWAL, error) {
	w := &SimpleWAL{filename: filename}

	// Open file once on creation for speed
	f, err := os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	w.file = f

	return w, nil
}

func (w *SimpleWAL) WriteEntry(userID string, count uint64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	const FORMAT = "%s,%d\n"
	entry := fmt.Sprintf(FORMAT, userID, count)
	_, err := w.file.WriteString(entry)
	if err != nil {
		return err
	}

	return nil
}

func (w *SimpleWAL) Recover() (map[string]uint64, error) {
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

func (w *SimpleWAL) Close() error {
	return w.file.Close()
}
