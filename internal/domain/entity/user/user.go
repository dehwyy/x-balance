package user

import "time"

type User struct {
	ID             ID
	Name           Name
	OverdraftLimit OverdraftLimit
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
