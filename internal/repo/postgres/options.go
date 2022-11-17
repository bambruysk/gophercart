package postgres

import (
	"fmt"
	"net/url"
	"time"
)

type Options struct {
	ConnectTimeout   time.Duration
	MigrationOptions *MigrationOptions
	DatabaseName     string
	Host             string
	Port             int
	auth             *AuthOptions
}

type MigrationOptions struct {
	Path    string
	Version int32
	Enable  bool
}

type AuthOptions struct {
	User     string
	Password string
}

// setDefaults - устанавливает дефолтные значения для Options.
func setDefaults(o *Options) {
	if o.Port == 0 {
		o.Port = 5432
	}

	if o.ConnectTimeout == 0 {
		o.ConnectTimeout = 10 * time.Second
	}

}

func (opt *Options) ConnectString() string {

	var userData, databaseName string
	if len(opt.auth.User) > 0 {
		userData += opt.auth.User
	}

	if len(opt.auth.Password) > 0 {
		userData += fmt.Sprintf(":%s", url.QueryEscape(opt.auth.Password))
	}

	if len(userData) > 0 {
		userData += "@"
	}

	if len(opt.DatabaseName) > 0 {
		databaseName += fmt.Sprintf("/%s", opt.DatabaseName)
	}

	values := make(url.Values, 0)

	connectionString := fmt.Sprintf(
		"postgres://%s%s:%d%s?%s",
		userData,
		opt.Host,
		opt.Port,
		databaseName,
		values.Encode(),
	)

	return connectionString
}
