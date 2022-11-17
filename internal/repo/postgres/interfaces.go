package postgres

import (
	"context"

	"gophercart/internal/models"
)

type Repo interface {
	Create(ctx context.Context, userID models.User, good models.Good) error
	List(ctx context.Context, userID models.User) ([]models.Good, error)
	Delete(ctx context.Context, userID models.User, good models.Good) error
	DeleteCart(ctx context.Context, cart models.Cart) error
	IsNotFound(err error) bool
}
