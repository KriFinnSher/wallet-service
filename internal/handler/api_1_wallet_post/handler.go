package api_1_wallet_post

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"wallet-service/internal/model"
)

type response struct {
	Code    int64   `json:"code"`
	Message string  `json:"message"`
	Balance *uint64 `json:"balance,omitempty"`
}

type Handler struct {
	log           *slog.Logger
	walletService walletService
}

func New(log *slog.Logger, walletService walletService) *Handler {
	return &Handler{
		log:           log,
		walletService: walletService,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var req model.TransactionBody
	reqCtx := ctx.Request().Context()

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	if !req.OperationType.Validate() {
		return ctx.JSON(http.StatusBadRequest, response{
			Code:    http.StatusBadRequest,
			Message: "invalid operation type",
		})
	}

	if !(req.Amount > 0) {
		return ctx.JSON(http.StatusBadRequest, response{
			Code:    http.StatusBadRequest,
			Message: "balance must be positive number",
		})
	}

	log := h.log.With(
		slog.String("wallet_id", req.WalletId.String()),
		slog.String("operation_type", string(req.OperationType)),
	)

	code, balance, err := h.walletService.MakeWalletOperation(reqCtx, req)
	if err != nil {
		log.ErrorContext(reqCtx, "failed to make wallet operation")
		return ctx.JSON(http.StatusInternalServerError, response{
			Code:    http.StatusInternalServerError,
			Message: "server error while wallet operation",
		})
	}

	switch code {
	case 0:
		return ctx.JSON(http.StatusOK, response{
			Code:    http.StatusOK,
			Message: "operation done successfully",
			Balance: balance,
		})
	case 1:
		return ctx.JSON(http.StatusInternalServerError, response{
			Code:    http.StatusInternalServerError,
			Message: "server error while wallet operation",
		})
	case 2:
		return ctx.JSON(http.StatusConflict, response{
			Code:    http.StatusConflict,
			Message: "insufficient balance",
			Balance: balance,
		})
	case 3:
		return ctx.JSON(http.StatusNotFound, response{
			Code:    http.StatusNotFound,
			Message: "wallet_id not found",
		})
	default:
		return nil
	}
}
