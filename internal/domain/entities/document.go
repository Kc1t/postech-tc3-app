package entities

import (
	"fmt"
	"strings"
	"unicode"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

// DocumentType identifica o tipo de documento brasileiro.
type DocumentType string

const (
	DocumentTypeCPF  DocumentType = "cpf"
	DocumentTypeCNPJ DocumentType = "cnpj"
)

// Document e um Value Object que representa um CPF ou CNPJ valido.
type Document struct {
	value   string       // somente digitos
	docType DocumentType // "cpf" ou "cnpj"
}

// NewDocument cria um Document a partir de uma string bruta, aplicando validacao completa.
func NewDocument(raw string) (Document, error) {
	digits := stripNonDigits(raw)

	switch len(digits) {
	case 11:
		if err := validateCPF(digits); err != nil {
			return Document{}, err
		}
		return Document{value: digits, docType: DocumentTypeCPF}, nil
	case 14:
		if err := validateCNPJ(digits); err != nil {
			return Document{}, err
		}
		return Document{value: digits, docType: DocumentTypeCNPJ}, nil
	default:
		return Document{}, domainerrors.ErrInvalidDocument
	}
}

// ReconstituteDocument restaura um Document a partir de dados persistidos (sem revalidar).
func ReconstituteDocument(value string) Document {
	digits := stripNonDigits(value)
	docType := DocumentTypeCPF
	if len(digits) == 14 {
		docType = DocumentTypeCNPJ
	}
	return Document{value: digits, docType: docType}
}

func (d Document) Value() string      { return d.value }
func (d Document) Type() DocumentType { return d.docType }
func (d Document) String() string     { return d.Formatted() }
func (d Document) IsZero() bool       { return d.value == "" }

// Formatted retorna o documento com mascara.
func (d Document) Formatted() string {
	switch d.docType {
	case DocumentTypeCPF:
		if len(d.value) != 11 {
			return d.value
		}
		return fmt.Sprintf("%s.%s.%s-%s", d.value[0:3], d.value[3:6], d.value[6:9], d.value[9:11])
	case DocumentTypeCNPJ:
		if len(d.value) != 14 {
			return d.value
		}
		return fmt.Sprintf("%s.%s.%s/%s-%s", d.value[0:2], d.value[2:5], d.value[5:8], d.value[8:12], d.value[12:14])
	default:
		return d.value
	}
}

func validateCPF(digits string) error {
	if allSameDigits(digits) {
		return domainerrors.ErrInvalidDocument
	}

	// Primeiro digito verificador: pesos 10..2
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(digits[i]-'0') * (10 - i)
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if int(digits[9]-'0') != d1 {
		return domainerrors.ErrInvalidDocument
	}

	// Segundo digito verificador: pesos 11..2
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(digits[i]-'0') * (11 - i)
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	if int(digits[10]-'0') != d2 {
		return domainerrors.ErrInvalidDocument
	}

	return nil
}

func validateCNPJ(digits string) error {
	if allSameDigits(digits) {
		return domainerrors.ErrInvalidDocument
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// Primeiro digito verificador
	sum := 0
	for i, w := range weights1 {
		sum += int(digits[i]-'0') * w
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if int(digits[12]-'0') != d1 {
		return domainerrors.ErrInvalidDocument
	}

	// Segundo digito verificador
	sum = 0
	for i, w := range weights2 {
		sum += int(digits[i]-'0') * w
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	if int(digits[13]-'0') != d2 {
		return domainerrors.ErrInvalidDocument
	}

	return nil
}

func stripNonDigits(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func allSameDigits(s string) bool {
	if len(s) == 0 {
		return true
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}
