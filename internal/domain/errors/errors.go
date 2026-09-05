package domainerrors

import "errors"

var (
	ErrNotFound                     = errors.New("not found")
	ErrAlreadyExists                = errors.New("already exists")
	ErrInvalidDocument              = errors.New("invalid CPF/CNPJ")
	ErrInvalidPlate                 = errors.New("invalid vehicle plate")
	ErrInvalidStatus                = errors.New("invalid status transition")
	ErrInvalidStatusValue           = errors.New("unknown order status")
	ErrStatusNotAllowedForRequester = errors.New("status transition not allowed for requester")
	ErrInsufficientStock            = errors.New("insufficient stock")
	ErrOrderNotCancellable          = errors.New("order cannot be cancelled at current status")
	ErrVehicleNotFromRequester      = errors.New("vehicle does not belong to the specified requester")
	ErrInvalidCredentials           = errors.New("invalid credentials")
	ErrInvalidRefreshToken          = errors.New("invalid or expired refresh token")
	ErrAccountLocked                = errors.New("account temporarily locked")
)
