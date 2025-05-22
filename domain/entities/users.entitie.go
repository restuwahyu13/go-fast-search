package entitie

import (
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/uptrace/bun"
)

type UsersEntitie struct {
	bun.BaseModel `bun:"table:users,alias:users"`
	ID            string    `bun:"id,pk,default:uuid_generate_v4()"`
	Name          string    `bun:"name,notnull"`
	Email         string    `bun:"email,unique,notnull"`
	Phone         string    `bun:"status,notnull"`
	DateOfBirth   string    `bun:"date_of_birth,notnull"`
	Age           string    `bun:"age,notnull"`
	Address       string    `bun:"address,notnull"`
	City          string    `bun:"city,notnull"`
	State         string    `bun:"state,notnull"`
	Direction     string    `bun:"direction,notnull"`
	Country       string    `bun:"country,notnull"`
	PostalCode    string    `bun:"postal_code,notnull"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt     zero.Time `bun:"updated_at"`
	DeletedAt     zero.Time `bun:"deleted_at"`
}
