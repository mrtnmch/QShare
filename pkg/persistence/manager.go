package persistence

import (
	"sync"
	"time"
)

const (
	CleanupSignalInterval = 100 * time.Millisecond
)

type Manager interface {
	Add(item Item)
	Get(code string) (*Item, bool)
	Remove(code string)
	Flush() error
	Close()
}

type Item struct {
	Code    string
	Content string
	Created time.Time
	Timeout time.Duration
}

type MemoryItemManager struct {
	items         map[string]Item
	itemsLock     sync.Mutex
	closeSignal   chan bool
	cleanupSignal <-chan time.Time
}

func isTimeOuted(item Item) bool {
	return item.Created.Before(time.Now().Add(-item.Timeout))
}

func NewMemoryItemManager() *MemoryItemManager {
	manager := &MemoryItemManager{
		items:         map[string]Item{},
		itemsLock:     sync.Mutex{},
		closeSignal:   make(chan bool),
		cleanupSignal: time.Tick(CleanupSignalInterval),
	}

	go func(manager *MemoryItemManager) {
	loop:
		for {
			select {
			case <-manager.cleanupSignal:
				manager.cleanup()
			case <-manager.closeSignal:
				break loop
			}
		}
	}(manager)

	return manager
}

func (manager *MemoryItemManager) cleanup() {
	var toDelete []string

	manager.itemsLock.Lock()
	defer manager.itemsLock.Unlock()

	for key, result := range manager.items {
		if isTimeOuted(result) {
			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		delete(manager.items, key)
	}
}

func (manager *MemoryItemManager) Add(item Item) {
	manager.itemsLock.Lock()
	item.Created = time.Now()
	manager.items[item.Code] = item
	manager.itemsLock.Unlock()
}

func (manager *MemoryItemManager) Get(code string) (*Item, bool) {
	manager.itemsLock.Lock()
	defer manager.itemsLock.Unlock()

	for itemCode, item := range manager.items {
		if itemCode == code && !isTimeOuted(item) {
			return &item, true
		}
	}

	return nil, false
}

func (manager *MemoryItemManager) Remove(code string) {
	manager.itemsLock.Lock()
	defer manager.itemsLock.Unlock()
	delete(manager.items, code)
}

func (manager *MemoryItemManager) Flush() error {
	return nil
}

func (manager *MemoryItemManager) Close() {
	manager.closeSignal <- true
}
