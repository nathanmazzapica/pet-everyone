package registry

import (
	"testing"
)

func TestHubRegistry_GetOrCreateHub(t *testing.T) {
	registry := NewHubRegistry()

	tests := []struct {
		name         string
		shouldCreate bool
	}{
		{"test1", true},
		{"abcdef", true},
		{"test1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hub, created := registry.GetOrCreateHub(tt.name)
			if hub == nil {
				t.Fatalf("Expected hub to be non-nil, got nil")
			}

			if created != tt.shouldCreate {
				t.Fatalf("Expected hub creation to be %v, got %v", tt.shouldCreate, created)
			}
		})
	}
}
