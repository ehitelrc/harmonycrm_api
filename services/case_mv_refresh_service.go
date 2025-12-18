package services

import (
	"fmt"
	"harmony_api/repository"
	"sync"
	"time"
)

type CaseMVRefreshService struct {
	repo *repository.MessageRepository
}

var (
	mvMutex        sync.Mutex
	lastRefresh    time.Time
	refreshRunning bool
	minInterval    = 10 * time.Second
)

func NewCaseMVRefreshService() *CaseMVRefreshService {
	return &CaseMVRefreshService{
		repo: &repository.MessageRepository{},
	}
}

// 🔐 ESTE es el método que debe usarse desde controllers
func (s *CaseMVRefreshService) RefreshOnEvent(reason string) {
	mvMutex.Lock()

	// ya hay un refresh corriendo
	if refreshRunning {
		mvMutex.Unlock()
		return
	}

	// aún no pasa la ventana mínima
	if time.Since(lastRefresh) < minInterval {
		mvMutex.Unlock()
		return
	}

	refreshRunning = true
	lastRefresh = time.Now()

	mvMutex.Unlock()

	go func() {
		start := time.Now()
		err := s.repo.RefreshCasesMVConcurrently()

		mvMutex.Lock()
		refreshRunning = false
		mvMutex.Unlock()

		if err != nil {
			fmt.Printf("❌ MV refresh failed (%s): %v\n", reason, err)
			return
		}

		fmt.Printf("✅ MV refreshed (%s) in %s\n", reason, time.Since(start))
	}()
}
