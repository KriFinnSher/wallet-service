package wallet_service

import (
	"context"

	"wallet-service/internal/model"
	storage "wallet-service/internal/storage/postgres/wallet"
)

type walletStorage interface {
	DoOperation(ctx context.Context, body model.TransactionBody) (storage.TransferCode, *uint64, error)
	GetBalanceByWalletId(ctx context.Context, walletId string) (*uint64, error)
}
