package entities

import (
	"regexp"
	"strings"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

// PlateFormat identifica o formato da placa brasileira.
type PlateFormat string

const (
	PlateFormatOld      PlateFormat = "old"
	PlateFormatMercosul PlateFormat = "mercosul"
)

// plateFormats mapeia cada formato de placa ao seu regex de validacao.
// Para suportar novos formatos, basta adicionar uma entrada aqui.
var plateFormats = []struct {
	format PlateFormat
	regex  *regexp.Regexp
}{
	{PlateFormatOld, regexp.MustCompile(`^[A-Z]{3}[0-9]{4}$`)},
	{PlateFormatMercosul, regexp.MustCompile(`^[A-Z]{3}[0-9][A-Z][0-9]{2}$`)},
}

// Plate e um Value Object que representa uma placa veicular brasileira valida.
type Plate struct {
	value  string      // normalizado: uppercase, sem hifen
	format PlateFormat // "old" ou "mercosul"
}

// NewPlate cria uma Plate a partir de uma string bruta, aplicando validacao.
func NewPlate(raw string) (Plate, error) {
	normalized := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(raw), "-", ""))

	if len(normalized) != 7 {
		return Plate{}, domainerrors.ErrInvalidPlate
	}

	for _, pf := range plateFormats {
		if pf.regex.MatchString(normalized) {
			return Plate{value: normalized, format: pf.format}, nil
		}
	}

	return Plate{}, domainerrors.ErrInvalidPlate
}

// ReconstitutePlate restaura uma Plate a partir de dados persistidos (sem revalidar).
func ReconstitutePlate(value string) Plate {
	normalized := strings.ToUpper(strings.ReplaceAll(value, "-", ""))
	for _, pf := range plateFormats {
		if pf.regex.MatchString(normalized) {
			return Plate{value: normalized, format: pf.format}
		}
	}
	return Plate{value: normalized, format: PlateFormatOld}
}

func (p Plate) Value() string       { return p.value }
func (p Plate) Format() PlateFormat { return p.format }
func (p Plate) IsZero() bool        { return p.value == "" }

// String retorna a placa formatada (com hifen para formato antigo).
func (p Plate) String() string {
	if p.format == PlateFormatOld && len(p.value) == 7 {
		return p.value[:3] + "-" + p.value[3:]
	}
	return p.value
}
