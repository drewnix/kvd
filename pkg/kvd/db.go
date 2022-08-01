package kvd

import (
	"errors"
	"sync"
)

type DB struct {
	mutex   *sync.RWMutex
	store   map[string]string
	metrics *Metrics
}

type Metrics struct {
	KeysStored int `json:"KeysStored"`
	GetOps     int `json:"GetOps"`
	SetOps     int `json:"SetOps"`
	DelOps     int `json:"DelOps"`
}

var (
	ErrInvalidKey     = errors.New("invalid key")
	ErrDatabaseClosed = errors.New("database closed")
)

func (db *DB) Init() error {
	db.mutex = &sync.RWMutex{}
	db.store = make(map[string]string, 0)
	db.metrics = &Metrics{
		KeysStored: 0,
		GetOps:     0,
		SetOps:     0,
		DelOps:     0,
	}

	return nil
}

func (db *DB) Get(key string) (string, error) {
	db.mutex.RLock()
	value, ok := db.store[key]
	db.metrics.GetOps++
	db.mutex.RUnlock()

	if !ok {
		return "", ErrInvalidKey
	}

	return value, nil
}

func (db *DB) Set(key string, value string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	_, existing := db.store[key]
	db.store[key] = value
	db.metrics.SetOps++
	if !existing {
		db.metrics.KeysStored++
	}

	return nil
}

func (db *DB) BulkSet(records []Record) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	for _, r := range records {
		_, existing := db.store[r.Key]
		db.store[r.Key] = r.Value
		db.metrics.SetOps++
		if !existing {
			db.metrics.KeysStored++
		}
	}

	return nil
}

func (db *DB) BulkGet(query []string) ([]Record, error) {
	var records []Record
	db.mutex.Lock()
	defer db.mutex.Unlock()
	for _, q := range query {
		raw, ok := db.store[q]
		if !ok {
			return nil, ErrInvalidKey
		}
		var rec = Record{
			Key:   q,
			Value: raw,
		}

		db.metrics.GetOps++
		records = append(records, rec)
	}

	return records, nil
}

func (db *DB) BulkDelete(query []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	for _, q := range query {
		delete(db.store, q)
		db.metrics.DelOps++
	}

	return nil
}

func (db *DB) Delete(key string) error {
	db.mutex.Lock()
	delete(db.store, key)
	db.metrics.DelOps++
	db.mutex.Unlock()

	return nil
}
