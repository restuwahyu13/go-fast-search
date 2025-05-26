package entitie

import (
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/uptrace/bun"
)

type (
	UsersEntitie struct {
		bun.BaseModel `bun:"table:users,alias:users"`
		ID            string    `json:"id" bun:"id,pk,default:uuid_generate_v4()"`
		Name          string    `json:"name" bun:"name,notnull"`
		Email         string    `json:"email" bun:"email,unique,notnull"`
		Phone         string    `json:"phone" bun:"phone,notnull"`
		DateOfBirth   string    `json:"date_of_birth" bun:"date_of_birth,notnull"`
		Age           string    `json:"age" bun:"age,notnull"`
		Address       string    `json:"address" bun:"address,notnull"`
		City          string    `json:"city" bun:"city,notnull"`
		State         string    `json:"state" bun:"state,notnull"`
		Direction     string    `json:"direction" bun:"direction,notnull"`
		Country       string    `json:"country" bun:"country,notnull"`
		PostalCode    string    `json:"postal_code" bun:"postal_code,notnull"`
		CreatedAt     time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
		UpdatedAt     zero.Time `json:"updated_at" bun:"updated_at,nullzero"`
		DeletedAt     zero.Time `json:"deleted_at" bun:"deleted_at,nullzero"`
	}

	UsersDocument struct {
		ID            string         `json:"id"`
		Name          string         `json:"name" `
		Email         string         `json:"email" `
		Phone         string         `json:"phone"`
		DateOfBirth   string         `json:"date_of_birth"`
		Age           string         `json:"age"`
		Address       string         `json:"address"`
		City          string         `json:"city"`
		State         string         `json:"state"`
		Direction     string         `json:"direction"`
		Country       string         `json:"country"`
		PostalCode    string         `json:"postal_code"`
		CreatedAt     int64          `json:"created_at"`
		UpdatedAt     int64          `json:"updated_at,omitempty"`
		DeletedAt     int64          `json:"deleted_at,omitempty"`
		Formatted     map[string]any `json:"_formatted,omitempty"`
		MatchPosition map[string]any `json:"_matchesPosition,omitempty"`
	}
)
