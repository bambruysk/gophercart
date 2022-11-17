package repo

import "gophercart/internal/repo/postgres"

type Options struct {
	PostgresOptions *postgres.Options `toml:"postgres"`
}

func (o *Options) getType() repoType {
	return repoTypePostgres
}
