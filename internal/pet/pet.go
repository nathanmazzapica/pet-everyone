package pet

import (
	"sync/atomic"
)

type Instance struct {
	petID         string
	petClickCount atomic.Uint64
}

func NewPetInstance(petID string) *Instance {
	return &Instance{
		petID:         petID,
		petClickCount: atomic.Uint64{}, // will get from db later
	}
}

func (p *Instance) Pat() {
	p.petClickCount.Add(1)
}
