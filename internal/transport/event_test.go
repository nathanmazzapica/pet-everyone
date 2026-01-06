package transport

import (
	"testing"
)

func TestTarget_IsTargetUser(t *testing.T) {
	tc := []struct {
		target Target
		expect bool
	}{
		{Target(""), false},
		{Target("abcde-fghij"), false},
		{Target("75442486-0878-440c-9db1-a7006c25a39f"), true},
	}

	for _, tt := range tc {
		if tt.target.IsTargetUser() != tt.expect {
			t.Errorf("Expected %v, got %v", tt.expect, tt.target.IsTargetUser())
		}
	}
}
