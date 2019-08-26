package storage

import (
	"errors"
	"hash/fnv"
	"inmemoryStorage/config"
	"io"
	"os"
	"sync"
	"time"
)

type Service struct {
	cfg *config.Config
	sync.RWMutex
	cache       *HashMap
	expirations map[int64][][]byte
	file        *os.File
}

type Item struct {
	Expiration int64
	Value      []byte
}

func New(cfg *config.Config, cache *HashMap, expirations map[int64][][]byte, file *os.File) *Service {
	return &Service{
		cfg:         cfg,
		cache:       cache,
		expirations: expirations,
		file:        file,
	}
}
func (s *Service) Stop() {
	s.file.Close()
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
	if item, found := s.cache.Get(key); found {
		for i, curItem := range s.expirations[item.Expiration] {
			if Equal(curItem, key) {
				s.expirations[item.Expiration] = append(s.expirations[item.Expiration][:i], s.expirations[item.Expiration][i+1:]...)
			}
		}
	}
	s.cache.Put(key, Item{
		Value:      value,
		Expiration: expiration,
	})
	s.expirations[expiration] = append(s.expirations[expiration], key)
	s.Unlock()
}

func (s *Service) Get(key []byte) ([]byte, bool) {
	s.RLock()
	item, found := s.cache.Get(key)

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

	if item, found := s.cache.Get(key); found {
		for i, curItem := range s.expirations[item.Expiration] {
			if Equal(curItem, key) {
				s.expirations[item.Expiration] = append(s.expirations[item.Expiration][:i], s.expirations[item.Expiration][i+1:]...)
			}
		}
	}

	if _, found := s.cache.Get(key); !found {
		s.Unlock()
		return errors.New("key not found")
	}

	s.cache.Delete(key)
	s.Unlock()
	return nil
}

func (s *Service) DeleteBatch(keys [][]byte) error {
	for _, key := range keys {
		if item, found := s.cache.Get(key); found {
			for i, curItem := range s.expirations[item.Expiration] {
				if Equal(curItem, key) {
					s.expirations[item.Expiration] = append(s.expirations[item.Expiration][:i], s.expirations[item.Expiration][i+1:]...)
				}
			}
		}

		if _, found := s.cache.Get(key); !found {
			return errors.New("key not found")
		}

		s.cache.Delete(key)
	}
	return nil
}

func (s *Service) DeleteExpired() {
	now := time.Now().UnixNano()
	s.Lock()
	for expiration, keys := range s.expirations {
		if expiration > 0 && now > expiration {
			if err := s.DeleteBatch(keys); err != nil {
				return
			}
		}
	}
	s.Unlock()

}
func (s *Service) Dump() error {
	s.Lock()
	cacheJSON, err := s.cache.MarshalJSON()
	if err != nil {
		s.Unlock()
		return err
	}
	err = s.file.Truncate(0)
	if err != nil {
		s.Unlock()
		return err
	}
	_, err = s.file.Seek(0, 0)
	if err != nil {
		s.Unlock()
		return err
	}
	_, err = s.file.Write(cacheJSON)
	if err != nil {
		s.Unlock()
		return err
	}
	s.Unlock()
	return nil
}

func (s *Service) Load() error {
	buf := make([]byte, 1024)
	s.Lock()
	var cacheJSON []byte
	for {
		n, err := s.file.Read(buf)
		if err != nil && err != io.EOF {
			s.Unlock()
			return err
		}
		if n == 0 {
			break
		}
		cacheJSON = append(cacheJSON, buf[:n]...)
	}

	if string(cacheJSON) != "" {
		err := s.cache.UnmarshalJSON(cacheJSON)
		if err != nil {
			s.Unlock()
			return err
		}
	}
	s.Unlock()
	return nil
}

func (s *Service) hash(key []byte) (uint64, error) {
	h := fnv.New64()
	_, err := h.Write(key)
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}
