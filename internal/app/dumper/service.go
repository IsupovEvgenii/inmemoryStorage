package dumper

import (
	"fmt"
	"inmemoryStorage/internal/app/storage"
	"time"
)

type Service struct {
	interval time.Duration
	stop     chan bool
	cache    *storage.Service
}

func New(duration time.Duration, stop chan bool, cache *storage.Service) *Service {
	return &Service{
		interval: duration,
		stop:     stop,
		cache:    cache,
	}
}

func (s *Service) Run() {
	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-ticker.C:
			s.cache.Dump()
		case <-s.stop:
			ticker.Stop()
			return
		}
	}
}

func (s *Service) Stop() {
	fmt.Println("dumper stop")
	s.stop <- true
}
