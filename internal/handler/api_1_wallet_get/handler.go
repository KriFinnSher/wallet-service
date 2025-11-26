package api_1_wallet_get

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
	walletId := ctx.Param("WALLET_UUID")
	err := uuid.Validate(walletId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response{
			Code:    http.StatusBadRequest,
			Message: "invalid wallet_id",
		})
	}

	reqCtx := ctx.Request().Context()

	log := h.log.With(
		slog.String("wallet_id", walletId),
	)

	balance, err := h.walletService.GetWalletBalance(reqCtx, walletId)
	if err != nil {
		log.ErrorContext(reqCtx, "failed to get wallet balance")
		return ctx.JSON(http.StatusInternalServerError, response{
			Code:    http.StatusInternalServerError,
			Message: "server error while fetching wallet balance",
		})
	}

	if balance == nil {
		log.WarnContext(reqCtx, "got unknown wallet_id")
		return ctx.JSON(http.StatusNotFound, response{
			Code:    http.StatusNotFound,
			Message: "wallet not found",
		})
	}

	return ctx.JSON(http.StatusOK, response{
		Code:    http.StatusOK,
		Message: "get balance success",
		Balance: balance,
	})
}
