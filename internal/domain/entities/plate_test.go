package entities

import (
	"testing"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

func TestNewPlate_FormatoAntigo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"com hifen", "ABC-1234", "ABC1234"},
		{"sem hifen", "ABC1234", "ABC1234"},
		{"minusculo", "abc-1234", "ABC1234"},
		{"com espacos", " ABC-1234 ", "ABC1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plate, err := NewPlate(tt.input)
			if err != nil {
				t.Fatalf("esperava sucesso, obteve erro: %v", err)
			}
			if plate.Value() != tt.want {
				t.Errorf("Value() = %q, esperava %q", plate.Value(), tt.want)
			}
			if plate.Format() != PlateFormatOld {
				t.Errorf("Format() = %q, esperava %q", plate.Format(), PlateFormatOld)
			}
		})
	}
}

func TestNewPlate_FormatoMercosul(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"mercosul padrao", "ABC1D23", "ABC1D23"},
		{"mercosul minusculo", "abc1d23", "ABC1D23"},
		{"mercosul com espaco", " BRA0S18 ", "BRA0S18"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plate, err := NewPlate(tt.input)
			if err != nil {
				t.Fatalf("esperava sucesso, obteve erro: %v", err)
			}
			if plate.Value() != tt.want {
				t.Errorf("Value() = %q, esperava %q", plate.Value(), tt.want)
			}
			if plate.Format() != PlateFormatMercosul {
				t.Errorf("Format() = %q, esperava %q", plate.Format(), PlateFormatMercosul)
			}
		})
	}
}

func TestNewPlate_Invalida(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"curta demais", "ABC123"},
		{"longa demais", "ABC12345"},
		{"somente numeros", "1234567"},
		{"somente letras", "ABCDEFG"},
		{"formato invalido", "AB12CD3"},
		{"caracteres especiais", "ABC-@#$"},
		{"vazia", ""},
		{"espacos", "       "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPlate(tt.input)
			if err == nil {
				t.Fatal("esperava erro, obteve sucesso")
			}
			if err != domainerrors.ErrInvalidPlate {
				t.Errorf("erro = %v, esperava ErrInvalidPlate", err)
			}
		})
	}
}

func TestPlate_String(t *testing.T) {
	t.Run("formato antigo com hifen", func(t *testing.T) {
		plate, _ := NewPlate("ABC1234")
		want := "ABC-1234"
		if got := plate.String(); got != want {
			t.Errorf("String() = %q, esperava %q", got, want)
		}
	})

	t.Run("formato mercosul sem hifen", func(t *testing.T) {
		plate, _ := NewPlate("ABC1D23")
		want := "ABC1D23"
		if got := plate.String(); got != want {
			t.Errorf("String() = %q, esperava %q", got, want)
		}
	})
}

func TestReconstitutePlate(t *testing.T) {
	t.Run("reconstitui formato antigo", func(t *testing.T) {
		plate := ReconstitutePlate("ABC1234")
		if plate.Value() != "ABC1234" {
			t.Errorf("Value() = %q, esperava %q", plate.Value(), "ABC1234")
		}
		if plate.Format() != PlateFormatOld {
			t.Errorf("Format() = %q, esperava %q", plate.Format(), PlateFormatOld)
		}
	})

	t.Run("reconstitui mercosul", func(t *testing.T) {
		plate := ReconstitutePlate("ABC1D23")
		if plate.Value() != "ABC1D23" {
			t.Errorf("Value() = %q, esperava %q", plate.Value(), "ABC1D23")
		}
		if plate.Format() != PlateFormatMercosul {
			t.Errorf("Format() = %q, esperava %q", plate.Format(), PlateFormatMercosul)
		}
	})
}

func TestPlate_IsZero(t *testing.T) {
	var zero Plate
	if !zero.IsZero() {
		t.Error("Plate zero-value deveria retornar IsZero() = true")
	}

	plate, _ := NewPlate("ABC1234")
	if plate.IsZero() {
		t.Error("Plate valida nao deveria retornar IsZero() = true")
	}
}
