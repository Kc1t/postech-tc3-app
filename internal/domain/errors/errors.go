package domainerrors

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrAlreadyExists         = errors.New("already exists")
	ErrInvalidDocument       = errors.New("invalid CPF/CNPJ")
	ErrInvalidPlate          = errors.New("invalid vehicle plate")
	ErrInvalidStatus         = errors.New("invalid status transition")
	ErrInvalidStatusValue    = errors.New("unknown order status")
	ErrInsufficientStock     = errors.New("insufficient stock")
	ErrOrderNotCancellable   = errors.New("order cannot be cancelled at current status")
	ErrVehicleNotFromCustomer = errors.New("vehicle does not belong to the specified customer")
)
