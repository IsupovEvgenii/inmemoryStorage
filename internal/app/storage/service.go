package storage

import (
	"encoding/binary"
	"errors"
	"fmt"
	"inmemoryStorage/config"
	"sync"
	"time"
)

type Service struct {
	cfg *config.Config
	sync.RWMutex
	cache       map[uint64]Item
	expirations map[int64][]uint64
}

type Item struct {
	Expiration int64
	Value      []byte
}

func New(cfg *config.Config, cache map[uint64]Item, expirations map[int64][]uint64) *Service {
	return &Service{
		cfg:         cfg,
		cache:       cache,
		expirations: expirations,
	}
}

func (s *Service) Set(key, value []byte, duration uint) {
	var expiration int64
	if duration == 0 {
		expiration = 0
	}
	if duration > 0 {
		expiration = time.Now().Add(time.Duration(duration) * time.Second).UnixNano()
	}

	s.Lock()
	if item, found := s.cache[encodeKey(key)]; found {
		for i, curItem := range s.expirations[item.Expiration] {
			if curItem == encodeKey(key) {
				s.expirations[item.Expiration] = append(s.expirations[item.Expiration][:i], s.expirations[item.Expiration][i+1:]...)
			}
		}
	}
	s.cache[encodeKey(key)] = Item{
		Value:      value,
		Expiration: expiration,
	}
	s.expirations[expiration] = append(s.expirations[expiration], encodeKey(key))
	s.Unlock()
}

func (s *Service) Get(key []byte) ([]byte, bool) {
	s.RLock()

	item, found := s.cache[encodeKey(key)]

	if !found {
		s.RUnlock()
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			s.RUnlock()
			return nil, false
		}

	}
	s.RUnlock()
	return item.Value, true
}

func (s *Service) Delete(key []byte) error {
	s.Lock()

	if item, found := s.cache[encodeKey(key)]; found {
		for i, curItem := range s.expirations[item.Expiration] {
			if curItem == encodeKey(key) {
				s.expirations[item.Expiration] = append(s.expirations[item.Expiration][:i], s.expirations[item.Expiration][i+1:]...)
			}
		}
	}

	if _, found := s.cache[encodeKey(key)]; !found {
		s.Unlock()
		return errors.New("key not found")
	}

	delete(s.cache, encodeKey(key))
	s.Unlock()
	return nil
}

func (s *Service) DeleteBatch(keys []uint64) error {

	for _, key := range keys {
		if item, found := s.cache[key]; found {
			for i, curItem := range s.expirations[item.Expiration] {
				if curItem == key {
					s.expirations[item.Expiration] = append(s.expirations[item.Expiration][:i], s.expirations[item.Expiration][i+1:]...)
				}
			}
		}

		if _, found := s.cache[key]; !found {
			return errors.New("key not found")
		}

		delete(s.cache, key)
	}
	return nil
}

func (s *Service) DeleteExpired() {
	now := time.Now().UnixNano()
	s.Lock()
	for k, v := range s.cache {
		fmt.Println(k, v.Value)
	}
	for expiration, keys := range s.expirations {
		if expiration > 0 && now > expiration {
			if err := s.DeleteBatch(keys); err != nil {
				return
			}
		}
	}
	for k, v := range s.cache {
		fmt.Println(k, v.Value)
	}
	s.Unlock()

}

func encodeKey(key []byte) uint64 {
	buf := make([]byte, 8)
	key = append(key, buf...)
	return binary.BigEndian.Uint64(key)
}
