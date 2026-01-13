package service

import (
	"fmt"
)

type SaboteurDatabase struct {
	ShouldFail bool
}

func (s *SaboteurDatabase) GetPetCount(id *string) (int64, error) {
	return 0, nil
}

func (s *SaboteurDatabase) UpdatePetCount(petID string, userID string, count uint64) error {
	if s.ShouldFail {
		return fmt.Errorf("database error")
	}
	return nil
}
