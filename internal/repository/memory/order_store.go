package memory

import (
	"context"
	"sync"

	"github.com/theBaffo/subito-backend-challenge/internal/domain"
)

// OrderStore is a thread-safe, in-memory implementation of OrderRepository.
type OrderStore struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

// NewOrderStore creates an empty OrderStore.
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[string]domain.Order),
	}
}

func (s *OrderStore) Save(_ context.Context, order *domain.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.ID] = *order
	return nil
}

func (s *OrderStore) FindByID(_ context.Context, id string) (*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	o, ok := s.orders[id]

	if !ok {
		return nil, domain.ErrOrderNotFound
	}
	
	return &o, nil
}
