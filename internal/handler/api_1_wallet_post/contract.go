package api_1_wallet_post

import (
	"context"

	"wallet-service/internal/model"
)

type walletService interface {
	MakeWalletOperation(ctx context.Context, trxBody model.TransactionBody) (uint8, *uint64, error)
}
