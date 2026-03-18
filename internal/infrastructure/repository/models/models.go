package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID             string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string          `gorm:"not null"`
	OverdraftLimit decimal.Decimal `gorm:"type:numeric;not null;default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}

type Event struct {
	ID              string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID          string          `gorm:"not null;index"`
	Type            string          `gorm:"not null"`
	Amount          decimal.Decimal `gorm:"type:numeric;not null"`
	TransactionID   string          `gorm:"uniqueIndex;not null"`
	SnapshotID      *string         `gorm:"type:uuid"`
	FreezeExpiresAt *time.Time
	CreatedAt       time.Time `gorm:"index"`
}

type Snapshot struct {
	ID        string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string          `gorm:"not null;index"`
	Balance   decimal.Decimal `gorm:"type:numeric;not null"`
	Version   int64           `gorm:"not null;default:0"`
	CreatedAt time.Time
}
