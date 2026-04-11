package entities

import (
	"testing"

	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"
)

func TestNewDocument_CPFValido(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"CPF sem mascara", "52998224725", "52998224725"},
		{"CPF com mascara", "529.982.247-25", "52998224725"},
		{"CPF com espacos", " 529.982.247-25 ", "52998224725"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewDocument(tt.input)
			if err != nil {
				t.Fatalf("esperava sucesso, obteve erro: %v", err)
			}
			if doc.Value() != tt.want {
				t.Errorf("Value() = %q, esperava %q", doc.Value(), tt.want)
			}
			if doc.Type() != DocumentTypeCPF {
				t.Errorf("Type() = %q, esperava %q", doc.Type(), DocumentTypeCPF)
			}
		})
	}
}

func TestNewDocument_CNPJValido(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"CNPJ sem mascara", "11222333000181", "11222333000181"},
		{"CNPJ com mascara", "11.222.333/0001-81", "11222333000181"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewDocument(tt.input)
			if err != nil {
				t.Fatalf("esperava sucesso, obteve erro: %v", err)
			}
			if doc.Value() != tt.want {
				t.Errorf("Value() = %q, esperava %q", doc.Value(), tt.want)
			}
			if doc.Type() != DocumentTypeCNPJ {
				t.Errorf("Type() = %q, esperava %q", doc.Type(), DocumentTypeCNPJ)
			}
		})
	}
}

func TestNewDocument_CPFInvalido(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"digitos repetidos 111", "11111111111"},
		{"digitos repetidos 000", "00000000000"},
		{"digitos repetidos 999", "99999999999"},
		{"digito verificador errado", "52998224726"},
		{"menos de 11 digitos", "1234567890"},
		{"mais de 11 menos de 14", "123456789012"},
		{"vazio", ""},
		{"somente letras", "abcdefghijk"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDocument(tt.input)
			if err == nil {
				t.Fatal("esperava erro, obteve sucesso")
			}
			if err != domainerrors.ErrInvalidDocument {
				t.Errorf("erro = %v, esperava ErrInvalidDocument", err)
			}
		})
	}
}

func TestNewDocument_CNPJInvalido(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"digitos repetidos", "11111111111111"},
		{"digito verificador errado", "11222333000182"},
		{"mais de 14 digitos", "112223330001811"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDocument(tt.input)
			if err == nil {
				t.Fatal("esperava erro, obteve sucesso")
			}
			if err != domainerrors.ErrInvalidDocument {
				t.Errorf("erro = %v, esperava ErrInvalidDocument", err)
			}
		})
	}
}

func TestDocument_Formatted(t *testing.T) {
	t.Run("CPF formatado", func(t *testing.T) {
		doc, _ := NewDocument("52998224725")
		want := "529.982.247-25"
		if got := doc.Formatted(); got != want {
			t.Errorf("Formatted() = %q, esperava %q", got, want)
		}
	})

	t.Run("CNPJ formatado", func(t *testing.T) {
		doc, _ := NewDocument("11222333000181")
		want := "11.222.333/0001-81"
		if got := doc.Formatted(); got != want {
			t.Errorf("Formatted() = %q, esperava %q", got, want)
		}
	})
}

func TestDocument_String(t *testing.T) {
	doc, _ := NewDocument("52998224725")
	if doc.String() != doc.Formatted() {
		t.Error("String() deveria retornar o mesmo que Formatted()")
	}
}

func TestReconstituteDocument(t *testing.T) {
	t.Run("reconstitui CPF sem revalidar", func(t *testing.T) {
		doc := ReconstituteDocument("52998224725")
		if doc.Value() != "52998224725" {
			t.Errorf("Value() = %q, esperava %q", doc.Value(), "52998224725")
		}
		if doc.Type() != DocumentTypeCPF {
			t.Errorf("Type() = %q, esperava %q", doc.Type(), DocumentTypeCPF)
		}
	})

	t.Run("reconstitui CNPJ sem revalidar", func(t *testing.T) {
		doc := ReconstituteDocument("11222333000181")
		if doc.Value() != "11222333000181" {
			t.Errorf("Value() = %q, esperava %q", doc.Value(), "11222333000181")
		}
		if doc.Type() != DocumentTypeCNPJ {
			t.Errorf("Type() = %q, esperava %q", doc.Type(), DocumentTypeCNPJ)
		}
	})
}

func TestDocument_IsZero(t *testing.T) {
	var zero Document
	if !zero.IsZero() {
		t.Error("Document zero-value deveria retornar IsZero() = true")
	}

	doc, _ := NewDocument("52998224725")
	if doc.IsZero() {
		t.Error("Document valido nao deveria retornar IsZero() = true")
	}
}
