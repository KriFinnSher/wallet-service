package api_1_wallet_get

import (
	"context"
)

type walletService interface {
	GetWalletBalance(ctx context.Context, walletId string) (*uint64, error)
}
