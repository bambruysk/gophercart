package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"

	"gophercart/internal/models"
)

func (s *storage) Create(ctx context.Context, userID models.User, good models.Good) error {
	_, err := s.db.Exec(ctx, "INSERT INTO carts (user_id, good_id, count) VALUES ($1, $2, $3)", userID, good.ID, good.Count)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) List(ctx context.Context, userID models.User) (res []models.Good, err error) {
	rows, err := s.db.Query(ctx, "SELECT good_id, count, created_at FROM carts WHERE user_id= $1", userID)
	if err != nil {
		return nil, err
	}

	var goodID pgtype.UUID
	var count pgtype.Int4
	var createdAt pgtype.Timestamptz

	for rows.Next() {
		if err = rows.Scan(&goodID, &count, &createdAt); err != nil {
			return nil, err
		}

		res = append(res, models.Good{
			ID:        goodID.Get().(uuid.UUID),
			Count:     count.Get().(int),
			CreatedAt: createdAt.Time,
		})
	}

	return res, nil
}

func (s storage) Delete(ctx context.Context, userID models.User, good models.Good) error {
	tags, err := s.db.Exec(ctx, "DELETE FROM carts WHERE user_id = ($1) and good_id = ($2)", userID, good.ID)
	if err != nil {
		return err
	}

	if tags.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (s storage) DeleteCart(ctx context.Context, cart models.Cart) error {
	tags, err := s.db.Exec(ctx, "DELETE FROM carts WHERE user_id = ($1) )", cart.User)
	if err != nil {
		return err
	}

	if tags.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (s storage) IsNotFound(err error) bool {
	return err == ErrRecordNotFound
}
