// Package storage should be refactored
// For now it uses: https://github.com/rapidloop/skv/blob/master/skv.go
package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Storage represents a key value store
type Storage struct {
	db     *bolt.DB
	bucket []byte
}

var (
	// ErrNotFound is returned when the key supplied to a Get or Delete
	// method does not exist in the database.
	ErrNotFound = errors.New("key not found")

	// ErrBadValue is returned when the value supplied to the Put method
	// is nil.
	ErrBadValue = errors.New("bad value")
)

// New creates a new instance of the key-value store
func New(path string, bucket string) (*Storage, error) {
	bucketName := []byte(bucket)
	opts := &bolt.Options{
		Timeout: 50 * time.Millisecond,
	}
	db, err := bolt.Open(path, 0640, opts)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Storage{db: db, bucket: bucketName}, nil
}

// Put an entry into the store
func (s *Storage) Put(key string, value interface{}) error {
	if value == nil {
		return ErrBadValue
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.bucket).Put([]byte(key), buf.Bytes())
	})
}

// Get retrives a key in database
func (s *Storage) Get(key string, value interface{}) error {
	return s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(s.bucket).Cursor()
		if k, v := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		} else if value == nil {
			return nil
		} else {
			d := gob.NewDecoder(bytes.NewReader(v))
			return d.Decode(value)
		}
	})
}

// Delete the entry with the given key
func (s *Storage) Delete(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(s.bucket).Cursor()
		if k, _ := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		}
		return c.Delete()
	})
}

// Close closes the key-value store file.
func (s *Storage) Close() error {
	return s.db.Close()
}
