package model

import "github.com/google/uuid"

type TransactionBody struct {
	WalletId      uuid.UUID           `json:"walletId" db:"wallet_id"`
	OperationType WalletOperationType `json:"operationType"`
	Amount        uint64              `json:"amount"`
}

type WalletOperationType string

const (
	DepositType  WalletOperationType = "DEPOSIT"
	WithDrawType WalletOperationType = "WITHDRAW"
)

func (w WalletOperationType) Validate() bool {
	switch w {
	case DepositType, WithDrawType:
		return true
	default:
		return false
	}
}
