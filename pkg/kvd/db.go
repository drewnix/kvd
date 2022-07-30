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
	KeysStored int `json:"keysStored"`
	GetOps     int `json:"getOps"`
	SetOps     int `json:"setOps"`
	DelOps     int `json:"delOps"`
}

var ErrorInvalidKey = errors.New("invalid key")

func (db *DB) Init() error {
	db.mutex = &sync.RWMutex{}
	db.store = make(map[string]string)
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
		return "", ErrorInvalidKey
	}

	return value, nil
}

func (db *DB) Set(key string, value string) error {
	db.mutex.Lock()
	db.store[key] = value
	db.metrics.SetOps++
	db.mutex.Unlock()

	return nil
}

func (db *DB) Delete(key string) error {
	db.mutex.Lock()
	delete(db.store, key)
	db.metrics.DelOps++
	db.mutex.Unlock()

	return nil
}
