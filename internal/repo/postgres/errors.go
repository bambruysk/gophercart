package postgres

import "errors"

var ErrRecordNotFound = errors.New("postgres: record not found")
var ErrNilOptions = errors.New("postgres: received nil options")
