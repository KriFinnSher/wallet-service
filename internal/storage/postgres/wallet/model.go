package wallet

type TransferCode uint8

const (
	SuccessCode TransferCode = iota
	ErrorCode
	InsufficientFundsCode
	NotFoundCode
)
