package repo

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"

	"gophercart/internal/repo/inmemory"
	"gophercart/internal/repo/postgres"
)

type repoType int

const (
	repoTypePostgres = iota
	repoTypeInMemory
)

var ErrUnknownTypeStorage = errors.New("repo: unknown type storage")

func NewRepository(ctx context.Context, opts *Options, log *logrus.Logger) (Repo, error) {
	switch opts.getType() {
	case repoTypePostgres:
		return postgres.NewStorage(ctx, log, opts.PostgresOptions)
	case repoTypeInMemory:
		return inmemory.NewStorage(log), nil
	default:
		return nil, ErrUnknownTypeStorage
	}
}
