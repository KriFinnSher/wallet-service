package model

import "github.com/google/uuid"

type Wallet struct {
	ID      uuid.UUID `json:"walletId" db:"wallet_id"`
	Balance uint64    `json:"balance" db:"balance"`
}
