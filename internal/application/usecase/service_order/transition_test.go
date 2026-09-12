package serviceorderuc

import (
	"reflect"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
)

func TestStatusTransition_TempoNoStatusAnterior(t *testing.T) {
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusInDiagnosis, nil, nil, 0, "", ft(), ft(), nil, nil)

	got := captureTransition(so).logAttrs(ft().Add(90 * time.Second))

	want := []any{"from_status", "in_diagnosis", "seconds_in_status", 90.0}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("esperava %v, obteve %v", want, got)
	}
}

func TestStatusTransition_SemHorarioDeEntrada(t *testing.T) {
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusInExecution, nil, nil, 0, "", ft(), time.Time{}, nil, nil)

	got := captureTransition(so).logAttrs(ft())

	want := []any{"from_status", "in_execution"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sem updated_at nao deve haver seconds_in_status; obteve %v", got)
	}
}

func TestStatusTransition_RelogioAtrasado(t *testing.T) {
	so := entities.ReconstituteServiceOrder("order-1", 100, "cust-1", "veh-1",
		entities.StatusFinished, nil, nil, 0, "", ft(), ft(), nil, nil)

	got := captureTransition(so).logAttrs(ft().Add(-time.Second))

	want := []any{"from_status", "finished"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("duracao negativa nao deve ser registrada; obteve %v", got)
	}
}
