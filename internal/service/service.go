package service

// I don't know if service has to be an interface yet, but maybe I'll discover it should be

type Service interface {
	Is()
}

type Event struct {
	originID string
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
}

func (e *Event) GetOriginID() string {
	return e.originID
}
