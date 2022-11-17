package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/sirupsen/logrus"
)

type storage struct {
	db   *pgx.Conn
	log  *logrus.Logger
	opts *Options
}

func NewStorage(ctx context.Context, log *logrus.Logger, opts *Options) (*storage, error) {
	if opts == nil {
		return nil, ErrNilOptions
	}

	setDefaults(opts)

	if log == nil {
		log = logrus.StandardLogger()
		log.Warnln("postgres: received nil logger, use standard logger ")
	}

	_ctx, cancel := context.WithTimeout(ctx, opts.ConnectTimeout)
	defer cancel()

	db, err := pgx.Connect(_ctx, opts.ConnectString())
	if err != nil {
		return nil, err
	}

	st := &storage{db: db, log: log, opts: opts}

	if err := st.Migrate(ctx); err != nil {
		return nil, err
	}

	return st, nil
}
