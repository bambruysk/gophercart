package postgres

import (
	"context"

	"github.com/jackc/tern/migrate"
	"github.com/pkg/errors"
)

func (st *storage) Migrate(ctx context.Context) error {
	if !st.opts.MigrationOptions.Enable {
		return nil
	}

	migrator, err := migrate.NewMigrator(ctx, st.db, "version")
	if err != nil {
		return errors.Wrap(err, "postgres: unable to create migrator")
	}

	err = migrator.LoadMigrations(st.opts.MigrationOptions.Path)
	if err != nil {
		return errors.Wrap(err, "postgres: unable to load migrations")
	}

	ver, err := migrator.GetCurrentVersion(ctx)

	if ver == 0 {
		st.log.Debugln("PostgreSQL: trying to migrate to last version")

		err = migrator.Migrate(ctx)
		if err != nil {
			return errors.Wrap(err, "postgres: unable to migrate")
		}

		st.log.Debugln("PostgreSQL: migration run successfully")

		return nil
	}

	if st.opts.MigrationOptions.Version > 0 && st.opts.MigrationOptions.Version != ver {
		st.log.Debugf("PostgreSQL: trying to migrate from %d to %d version", ver, st.opts.MigrationOptions.Version)

		err = migrator.MigrateTo(ctx, st.opts.MigrationOptions.Version)
		if err != nil {
			return errors.Wrap(err, "postgres: unable to migrate")
		}

		return nil
	}

	return nil
}
