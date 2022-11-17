package models

import (
	"time"

	"github.com/google/uuid"
)

type User uuid.UUID

type Good struct {
	ID        uuid.UUID
	Count     int
	CreatedAt time.Time
}

type Cart struct {
	User  User
	Goods []Good
}
