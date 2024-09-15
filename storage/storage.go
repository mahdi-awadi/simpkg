package storage

import (
	"sync"
	"time"
)

// StorageEvent type
type StorageEvent string

const (
	EventInsert StorageEvent = "insert"
	EventRemove StorageEvent = "remove"
	EventClear  StorageEvent = "clear"
)

// IStorage interface
type IStorage interface {
	Insert(key string, value any, removeAfter ...time.Duration)
	Remove(key string)
	Clear()
	Get(key string) ([]*Row, bool)
	Has(key string) bool
	Count(key string) int
	First(key string) *Row
	Last(key string) *Row
	All() map[string][]*Row
	Size() int
	Subscribe(fn func(Event))
}

// Storage struct
type Storage struct {
	storage     map[string][]*Row
	subscribers []func(Event)
	locker      sync.RWMutex
}

// Row struct
type Row struct {
	Index int
	Value any
	Time  time.Time
}

// Event struct
type Event struct {
	Key   string
	Value any
	Event StorageEvent
}

// New returns a new Storage instance
func New() IStorage {
	return &Storage{
		storage:     make(map[string][]*Row),
		subscribers: make([]func(Event), 0),
		locker:      sync.RWMutex{},
	}
}

// Insert value to storage
func (storage *Storage) Insert(key string, value any, removeAfter ...time.Duration) {
	initialized := storage.Has(key)
	storage.locker.Lock()

	if !initialized {
		storage.storage[key] = make([]*Row, 0)
	}

	storage.storage[key] = append(storage.storage[key], &Row{Value: value, Time: time.Now()})
	storage.locker.Unlock()

	if len(removeAfter) > 0 {
		time.AfterFunc(removeAfter[0], func() {
			storage.Remove(key)
		})
	}

	storage.notify(Event{Key: key, Value: value, Event: EventInsert})
}

// Remove value from storage
func (storage *Storage) Remove(key string) {
	storage.locker.Lock()
	delete(storage.storage, key)
	storage.locker.Unlock()

	storage.notify(Event{Key: key, Event: EventRemove})
}

// Clear storage
func (storage *Storage) Clear() {
	storage.locker.Lock()
	storage.storage = make(map[string][]*Row)
	storage.locker.Unlock()

	storage.notify(Event{Event: EventClear})
}

// Get value from storage
func (storage *Storage) Get(key string) ([]*Row, bool) {
	storage.locker.RLock()
	row, exists := storage.storage[key]
	storage.locker.RUnlock()

	return row, exists
}

// First returns first value from storage
func (storage *Storage) First(key string) *Row {
	rows, exists := storage.Get(key)
	if !exists {
		return nil
	}

	return rows[0]
}

// Last returns last value from storage
func (storage *Storage) Last(key string) *Row {
	rows, exists := storage.Get(key)
	if !exists {
		return nil
	}

	return rows[len(rows)-1]
}

// Has value in storage
func (storage *Storage) Has(key string) bool {
	storage.locker.RLock()
	_, ok := storage.storage[key]
	storage.locker.RUnlock()

	return ok
}

// Count returns count of values in storage
func (storage *Storage) Count(key string) int {
	rows, exists := storage.Get(key)
	if !exists {
		return 0
	}

	return len(rows)
}

// All returns all values from storage
func (storage *Storage) All() map[string][]*Row {
	return storage.storage
}

// Size returns storage size
func (storage *Storage) Size() int {
	return len(storage.storage)
}

// Subscribe to storage changes
func (storage *Storage) Subscribe(fn func(Event)) {
	storage.subscribers = append(storage.subscribers, fn)
}

// notify to subscribers
func (storage *Storage) notify(data Event) {
	for _, fn := range storage.subscribers {
		go fn(data)
	}
}
