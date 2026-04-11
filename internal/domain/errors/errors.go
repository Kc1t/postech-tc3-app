// Package errors centraliza os erros de dominio do sistema.
// Repositories devem converter erros de infraestrutura (ex: gorm.ErrRecordNotFound)
// pra estes, mantendo a camada de dominio desacoplada do ORM/driver.
package errors

import "errors"

var (
	// ErrNotFound indica que o recurso solicitado nao existe.
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists indica violacao de constraint de unicidade.
	ErrAlreadyExists = errors.New("resource already exists")
)
