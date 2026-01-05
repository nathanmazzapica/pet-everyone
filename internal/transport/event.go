package transport

import "strings"

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

func (t Target) IsBroadcastExcept() bool {
	return strings.HasPrefix(string(t), "broadcast-except:")
}

func (t Target) GetExceptUserID() string {
	return strings.TrimPrefix(string(t), "broadcast-except:")
}

func (t Target) IsTargetUser() bool {
	return !t.IsBroadcast() && !t.IsBroadcastExcept()
}

func (t Target) GetUserID() string {
	return string(t)
}

type Event struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	Target Target      `json:"-"` // Don't include in JSON output
}
