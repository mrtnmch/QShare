package repository

import (
	"time"
)

type Envelope struct {
	Key     string
	Content []byte
	Created time.Time
	Timeout time.Duration
}

func (envelope *Envelope) hasTimedOut() bool {
	return envelope.Created.Before(time.Now().Add(-envelope.Timeout))
}
