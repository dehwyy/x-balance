package user_test

import (
	"testing"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUser_ValueObjects(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		userName       string
		overdraftLimit string
		expectValid    bool
	}{
		{
			name:           "Valid user with positive overdraft",
			id:             "user-123",
			userName:       "John Doe",
			overdraftLimit: "1000",
			expectValid:    true,
		},
		{
			name:           "Valid user with zero overdraft",
			id:             "user-456",
			userName:       "Jane Smith",
			overdraftLimit: "0",
			expectValid:    true,
		},
		{
			name:           "Valid user with decimal overdraft",
			id:             "user-789",
			userName:       "Bob Johnson",
			overdraftLimit: "500.50",
			expectValid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overdraft, err := decimal.NewFromString(tt.overdraftLimit)
			assert.NoError(t, err)

			u := &user.User{
				ID:             user.ID{Value: tt.id},
				Name:           user.Name{Value: tt.userName},
				OverdraftLimit: user.OverdraftLimit{Value: overdraft},
			}

			// Test value object access
			assert.Equal(t, tt.id, u.ID.Value)
			assert.Equal(t, tt.userName, u.Name.Value)
			assert.Equal(t, overdraft.String(), u.OverdraftLimit.Value.String())
		})
	}
}

func TestOverdraftLimit_NegativeValues(t *testing.T) {
	// Test that negative overdraft limits are handled correctly
	negativeOverdraft, _ := decimal.NewFromString("-100")

	u := &user.User{
		ID:             user.ID{Value: "test-user"},
		Name:           user.Name{Value: "Test User"},
		OverdraftLimit: user.OverdraftLimit{Value: negativeOverdraft},
	}

	// Negative overdraft limit should still be stored as provided
	// Business logic validation should happen at the service layer
	assert.Equal(t, "-100", u.OverdraftLimit.Value.String())
}
