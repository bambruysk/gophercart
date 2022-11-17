package usecase

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"

	"gophercart/internal/models"
)

type processor struct {
	logger    *logrus.Logger
	repo      Repo
	warehouse Warehouse
}

func New(logger *logrus.Logger, repo Repo, warehouse Warehouse) *processor {
	return &processor{logger: logger, repo: repo, warehouse: warehouse}
}

// Repo
//
//go:generate mockery --name Repo --with-expecter --keeptree
type Repo interface {
	Create(ctx context.Context, userID models.User, good models.Good) error
	List(ctx context.Context, userID models.User) ([]models.Good, error)
	Delete(ctx context.Context, userID models.User, good models.Good) error
	DeleteCart(ctx context.Context, cart models.Cart) error
	IsNotFound(err error) bool
}

// Repo
//
//go:generate mockery --name Warehouse --with-expecter --keeptree
type Warehouse interface {
	CheckGood(ctx context.Context, good models.Good) (count int, err error)
	DeliveryGood(ctx context.Context, good models.Good) error
}

func (p *processor) Add(ctx context.Context, userID models.User, good models.Good) error {
	count, err := p.warehouse.CheckGood(ctx, good)
	if err != nil {
		return ErrWarehouseErr
	}

	if good.Count > count {
		return ErrTooLowInWarehouse
	}

	if err = p.repo.Create(ctx, userID, good); err != nil {
		return ErrDatabaseErr
	}

	return nil
}

func (p *processor) Delete(ctx context.Context, userID models.User, good models.Good) error {
	if err := p.repo.Delete(ctx, userID, good); err != nil {
		if p.repo.IsNotFound(err) {
			return ErrGoodNotFound
		}
		return ErrDatabaseErr
	}
	return nil
}

func (p *processor) BuyCart(ctx context.Context, userID models.User) error {
	cartContent, err := p.repo.List(ctx, userID)
	if err != nil {
		if p.repo.IsNotFound(err) {
			return ErrGoodNotFound
		}
		return ErrDatabaseErr
	}

	for _, good := range cartContent {
		count, err := p.warehouse.CheckGood(ctx, good)
		if err != nil {
			return ErrGoodNotFound
		}
		if count < good.Count {
			return ErrTooLowInWarehouse
		}
	}

	cart := models.Cart{
		User:  userID,
		Goods: cartContent,
	}

	for _, good := range cartContent {
		err = p.warehouse.DeliveryGood(ctx, good)
		if err != nil {
			return ErrGoodNotFound
		}
	}

	if err := p.repo.DeleteCart(ctx, cart); err != nil {
		if p.repo.IsNotFound(err) {
			return ErrGoodNotFound
		}
		return ErrDatabaseErr
	}
	return nil
}

var ErrTooLowInWarehouse = errors.New("usecase: low goods in warehouse")
var ErrDatabaseErr = errors.New("usecase: database error")
var ErrWarehouseErr = errors.New("usecase: warehouse error")
var ErrGoodNotFound = errors.New("usecase: good not found")
