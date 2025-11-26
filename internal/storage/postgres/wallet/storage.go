package wallet

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"wallet-service/internal/model"
)

const (
	tableName = "wallets"

	operationDoOperation          = "do_operation"
	operationGetBalanceByWalletId = "get_balance_by_wallet_id"
)

type Storage struct {
	db  *sqlx.DB
	log *slog.Logger
}

func New(db *sqlx.DB, log *slog.Logger) *Storage {
	return &Storage{
		db:  db,
		log: log,
	}
}

func (s *Storage) DoOperation(ctx context.Context, body model.TransactionBody) (TransferCode, *uint64, error) {
	log := s.log.With(
		slog.String("operation", operationDoOperation),
	)

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.ErrorContext(ctx, "failed to start tx", "error", err)
		return ErrorCode, nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	sb := sq.Select("wallet_id", "balance").
		From(tableName).
		Where(sq.Eq{"wallet_id": body.WalletId}).
		Suffix("FOR NO KEY UPDATE").
		PlaceholderFormat(sq.Dollar)

	query, args, err := sb.ToSql()
	if err != nil {
		log.ErrorContext(ctx, "failed to build select query", "error", err)
		return ErrorCode, nil, err
	}

	var modelWallet model.Wallet
	err = tx.GetContext(ctx, &modelWallet, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NotFoundCode, nil, nil
		}

		log.ErrorContext(ctx, "failed to exec select query", "error", err)
		return ErrorCode, nil, err
	}

	delta := s.getDeltaValue(body)
	if int64(modelWallet.Balance)+delta < 0 {
		return InsufficientFundsCode, lo.ToPtr(modelWallet.Balance), nil // не ходим лишний раз в бд
	}

	ub := sq.Update(tableName).
		Set("balance", sq.Expr("balance + ?", delta)).
		Where(sq.Eq{"wallet_id": body.WalletId}).
		Suffix("RETURNING balance").
		PlaceholderFormat(sq.Dollar)

	query, args, err = ub.ToSql()
	if err != nil {
		log.ErrorContext(ctx, "failed to build update query", "error", err)
		return ErrorCode, nil, err
	}

	var newBalance uint64
	err = tx.GetContext(ctx, &newBalance, query, args...)
	if err != nil {
		log.ErrorContext(ctx, "failed to exec update query", "error", err)
		return ErrorCode, nil, err
	}

	if err = tx.Commit(); err != nil {
		s.log.ErrorContext(ctx, "failed to commit tx", "error", err)
		return ErrorCode, nil, err
	}

	return SuccessCode, &newBalance, nil
}

func (s *Storage) GetBalanceByWalletId(ctx context.Context, walletId string) (*uint64, error) {
	log := s.log.With(
		slog.String("operation", operationGetBalanceByWalletId),
	)

	sb := sq.Select("balance").
		From(tableName).
		Where(sq.Eq{"wallet_id": walletId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sb.ToSql()
	if err != nil {
		log.ErrorContext(ctx, "failed to build query", "error", err)
		return nil, err
	}

	var b uint64
	err = s.db.GetContext(ctx, &b, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.ErrorContext(ctx, "failed to exec query", "error", err)
		return nil, err
	}

	return &b, nil
}

func (s *Storage) getDeltaValue(transactionBody model.TransactionBody) int64 {
	switch transactionBody.OperationType {
	case model.DepositType:
		return int64(transactionBody.Amount)
	case model.WithDrawType:
		return -int64(transactionBody.Amount)
	default:
		return 0
	}
}
