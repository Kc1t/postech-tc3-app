package domainerrors

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrInvalidDocument     = errors.New("invalid CPF/CNPJ")
	ErrInvalidPlate        = errors.New("invalid plate")
	ErrInvalidStatus       = errors.New("invalid status transition")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrOrderNotCancellable = errors.New("order cannot be cancelled at current status")
)
