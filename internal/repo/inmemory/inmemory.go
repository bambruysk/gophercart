package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

	"gophercart/internal/models"
)

var ErrRecordNotFound = errors.New("inmemory: record not found")

type storage struct {
	log *logrus.Logger

	data map[models.User][]models.Good
	mtx  sync.Mutex
}

func (s *storage) Create(ctx context.Context, userID models.User, good models.Good) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	cart := s.data[userID]
	cart = append(cart, good)
	s.data[userID] = cart

	return nil
}

func (s *storage) List(ctx context.Context, userID models.User) ([]models.Good, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	cart, ok := s.data[userID]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return cart, nil
}

func (s *storage) Delete(ctx context.Context, userID models.User, good models.Good) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	cart, ok := s.data[userID]
	if !ok {
		return ErrRecordNotFound
	}

	for i, c := range cart {
		if c.ID == good.ID {
			cart = append(cart[:i], cart[i+1:]...)
			break
		}
	}

	if len(cart) == 0 {
		delete(s.data, userID)
	}

	return nil
}

func (s *storage) DeleteCart(ctx context.Context, cart models.Cart) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, ok := s.data[cart.User]
	if !ok {
		return ErrRecordNotFound
	}

	delete(s.data, cart.User)

	return nil
}

func (s *storage) IsNotFound(err error) bool {
	return err == ErrRecordNotFound
}

func NewStorage(log *logrus.Logger) *storage {
	return &storage{log: log, data: make(map[models.User][]models.Good)}
}
