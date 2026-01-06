package transport

import (
	"strings"

	"github.com/google/uuid"
)

type Target string

const (
	TargetBroadcast = Target("broadcast")
)

func TargetExcept(userID string) Target {
	return Target("broadcast-except:" + userID)
}

func TargetUser(userID string) Target {
	return Target(userID)
}

func (t Target) IsBroadcast() bool {
	return t == TargetBroadcast
}

// IsBroadcastExcept checks if the target string is prefixed with "broadcast-except:".
func (t Target) IsBroadcastExcept() bool {
	return strings.HasPrefix(string(t), "broadcast-except:")
}

// GetExceptUserID extracts and returns the user ID from a target string prefixed with "broadcast-except:".
func (t Target) GetExceptUserID() string {
	return strings.TrimPrefix(string(t), "broadcast-except:")
}

// IsTargetUser returns true if the target is a user ID
func (t Target) IsTargetUser() bool {
	err := uuid.Validate(string(t))
	return err == nil
}

func (t Target) GetUserID() string {
	return string(t)
}

type Event struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	Target Target      `json:"-"` // Don't include in JSON output
}
