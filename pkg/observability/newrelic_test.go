package observability_test

import (
	"testing"

	"github.com/fiap/postech-tc1/pkg/observability"
)

func TestNewRelic_DisabledWithoutLicenseKey(t *testing.T) {
	app, err := observability.NewRelic("workshop-api", "")
	if err != nil {
		t.Fatalf("esperava nil error, veio %v", err)
	}
	if app != nil {
		t.Error("esperava aplicacao nil quando nao ha license key")
	}
}

func TestNewRelic_RejectsInvalidLicenseKey(t *testing.T) {
	if _, err := observability.NewRelic("workshop-api", "chave-invalida"); err == nil {
		t.Error("esperava erro para license key malformada")
	}
}
