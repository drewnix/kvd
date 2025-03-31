package kvd

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Common errors
var (
	ErrInvalidKey   = errors.New("invalid key")
	ErrKeyNotFound  = errors.New("key not found")
	ErrEmptyKey     = errors.New("empty key not allowed")
	ErrNilValue     = errors.New("nil value not allowed")
)

// DB represents the key-value database
type DB struct {
	mutex   *sync.RWMutex
	store   map[string]string
	metrics *Metrics
}

// Metrics tracks usage statistics for the database
type Metrics struct {
	KeysStored       int64 `json:"KeysStored"`
	ValueBytesStored int64 `json:"ValueBytesStored"`
	GetOps           int64 `json:"GetOps"`
	SetOps           int64 `json:"SetOps"`
	DelOps           int64 `json:"DelOps"`
}

// Init initializes the database
func (db *DB) Init() error {
	db.mutex = &sync.RWMutex{}
	db.store = make(map[string]string)
	db.metrics = &Metrics{
		KeysStored:       0,
		ValueBytesStored: 0,
		GetOps:           0,
		SetOps:           0,
		DelOps:           0,
	}

	return nil
}

// Get retrieves a value for a given key
func (db *DB) Get(key string) (string, error) {
	// Increment operations counter regardless of result
	atomic.AddInt64(&db.metrics.GetOps, 1)
	
	if key == "" {
		return "", ErrEmptyKey
	}

	db.mutex.RLock()
	value, ok := db.store[key]
	db.mutex.RUnlock()

	if !ok {
		return "", ErrKeyNotFound
	}

	return value, nil
}

// Set stores a key-value pair
func (db *DB) Set(key string, value string) error {
	if key == "" {
		return ErrEmptyKey
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()
	
	oldValue, existing := db.store[key]
	db.store[key] = value
	atomic.AddInt64(&db.metrics.SetOps, 1)
	
	if existing {
		// Update bytes stored (subtract old value size, add new value size)
		oldBytes := len(oldValue)
		newBytes := len(value)
		atomic.AddInt64(&db.metrics.ValueBytesStored, int64(newBytes-oldBytes))
	} else {
		// New key
		atomic.AddInt64(&db.metrics.KeysStored, 1)
		atomic.AddInt64(&db.metrics.ValueBytesStored, int64(len(value)))
	}

	return nil
}

// BulkSet sets multiple key-value pairs atomically
func (db *DB) BulkSet(records []Record) error {
	if len(records) == 0 {
		return nil
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()
	
	// Validate all keys and values first
	for _, r := range records {
		if r.Key == "" {
			return ErrEmptyKey
		}
	}
	
	// Process all records
	for _, r := range records {
		oldValue, existing := db.store[r.Key]
		db.store[r.Key] = r.Value
		
		if existing {
			// Update metrics for existing key
			oldBytes := len(oldValue)
			newBytes := len(r.Value)
			db.metrics.ValueBytesStored += int64(newBytes - oldBytes)
		} else {
			// Update metrics for new key
			db.metrics.KeysStored++
			db.metrics.ValueBytesStored += int64(len(r.Value))
		}
	}
	
	// Update operation count once for the entire batch
	atomic.AddInt64(&db.metrics.SetOps, int64(len(records)))

	return nil
}

// BulkGet retrieves multiple values by their keys
func (db *DB) BulkGet(keys []string) ([]Record, error) {
	if len(keys) == 0 {
		return []Record{}, nil
	}

	// Pre-allocate the slice for efficiency
	records := make([]Record, 0, len(keys))
	
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	for _, key := range keys {
		if key == "" {
			return nil, ErrEmptyKey
		}
		
		value, ok := db.store[key]
		if !ok {
			return nil, ErrKeyNotFound
		}
		
		records = append(records, Record{
			Key:   key,
			Value: value,
		})
	}
	
	// Update operation count once for the entire batch
	atomic.AddInt64(&db.metrics.GetOps, int64(len(keys)))

	return records, nil
}

// Delete removes a key-value pair
func (db *DB) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()
	
	value, exists := db.store[key]
	if !exists {
		return ErrKeyNotFound
	}
	
	delete(db.store, key)
	
	// Update metrics
	atomic.AddInt64(&db.metrics.DelOps, 1)
	atomic.AddInt64(&db.metrics.KeysStored, -1)
	atomic.AddInt64(&db.metrics.ValueBytesStored, -int64(len(value)))

	return nil
}

// BulkDelete removes multiple key-value pairs
func (db *DB) BulkDelete(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()
	
	// First, check if all keys exist
	for _, key := range keys {
		if key == "" {
			return ErrEmptyKey
		}
		
		_, exists := db.store[key]
		if !exists {
			return ErrKeyNotFound
		}
	}
	
	// Then delete all keys
	totalBytesRemoved := int64(0)
	for _, key := range keys {
		value := db.store[key]
		totalBytesRemoved += int64(len(value))
		delete(db.store, key)
	}
	
	// Update metrics
	atomic.AddInt64(&db.metrics.DelOps, int64(len(keys)))
	atomic.AddInt64(&db.metrics.KeysStored, -int64(len(keys)))
	atomic.AddInt64(&db.metrics.ValueBytesStored, -totalBytesRemoved)

	return nil
}
