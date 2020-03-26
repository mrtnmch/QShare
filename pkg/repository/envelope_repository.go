package repository

import (
	"sync"
	"time"
)

const (
	CleanupSignalInterval = 100 * time.Millisecond
)

type EnvelopeRepository interface {
	Add(envelope *Envelope)
	Get(key string) *Envelope
	GetAll() []*Envelope
	Remove(key string)
}

type InMemoryEnvelopeRepository struct {
	envelopes     map[string]*Envelope
	envelopesLock sync.Mutex
}

func NewInMemoryEnvelopeRepository() *InMemoryEnvelopeRepository {
	return &InMemoryEnvelopeRepository{
		envelopes:     map[string]*Envelope{},
		envelopesLock: sync.Mutex{},
	}
}

func (manager *InMemoryEnvelopeRepository) Add(envelope *Envelope) {
	manager.envelopesLock.Lock()
	manager.envelopes[envelope.Key] = envelope
	manager.envelopesLock.Unlock()
}

func (manager *InMemoryEnvelopeRepository) Get(key string) *Envelope {
	manager.envelopesLock.Lock()
	defer manager.envelopesLock.Unlock()

	for envelopeKey, envelope := range manager.envelopes {
		if envelopeKey == key {
			return envelope
		}
	}

	return nil
}

func (manager *InMemoryEnvelopeRepository) GetAll() []*Envelope {
	var envelopes []*Envelope

	for _, envelope := range manager.envelopes {
		envelopes = append(envelopes, envelope)
	}

	return envelopes
}


func (manager *InMemoryEnvelopeRepository) Remove(key string) {
	manager.envelopesLock.Lock()
	delete(manager.envelopes, key)
	manager.envelopesLock.Unlock()
}
