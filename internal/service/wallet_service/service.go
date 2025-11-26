package wallet_service

import (
	"context"
	"log/slog"

	"wallet-service/internal/model"
)

const (
	operationMakeWalletOperation = "make_wallet_operation"
	operationGetWalletBalance    = "get_wallet_balance"
)

type Service struct {
	log           *slog.Logger
	walletStorage walletStorage
}

func New(log *slog.Logger, walletStorage walletStorage) *Service {
	return &Service{
		log:           log,
		walletStorage: walletStorage,
	}
}

func (s *Service) MakeWalletOperation(ctx context.Context, trxBody model.TransactionBody) (uint8, *uint64, error) {
	log := s.log.With(
		slog.String("operation", operationMakeWalletOperation),
	)

	code, updatedBalance, err := s.walletStorage.DoOperation(ctx, trxBody)
	if err != nil {
		log.ErrorContext(ctx, "failed to do operation in storage")
		return uint8(code), nil, err
	}

	return uint8(code), updatedBalance, nil
}

func (s *Service) GetWalletBalance(ctx context.Context, walletId string) (*uint64, error) {
	log := s.log.With(
		slog.String("operation", operationGetWalletBalance),
	)

	b, err := s.walletStorage.GetBalanceByWalletId(ctx, walletId)
	if err != nil {
		log.ErrorContext(ctx, "failed to get balance from storage")
		return nil, err
	}

	return b, nil
}
