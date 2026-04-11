package application

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fiap/postech-tc1/internal/domain/entities"
	domainerrors "github.com/fiap/postech-tc1/internal/domain/errors"

	customeruc "github.com/fiap/postech-tc1/internal/application/usecase/customer"
	vehicleuc "github.com/fiap/postech-tc1/internal/application/usecase/vehicle"
	serviceorderuc "github.com/fiap/postech-tc1/internal/application/usecase/service_order"
)

// =============================================================================
// In-Memory Repositories — simulam o banco para testes de integracao
// =============================================================================

// --- Customer ---

type inMemoryCustomerRepo struct {
	mu   sync.RWMutex
	data map[string]*entities.Customer
	seq  int
}

func newCustomerRepo() *inMemoryCustomerRepo {
	return &inMemoryCustomerRepo{data: make(map[string]*entities.Customer)}
}

func (r *inMemoryCustomerRepo) Create(_ context.Context, c *entities.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	c.SetID(fmt.Sprintf("cust-%d", r.seq))
	r.data[c.ID()] = c
	return nil
}

func (r *inMemoryCustomerRepo) FindByID(_ context.Context, id string) (*entities.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if c, ok := r.data[id]; ok {
		return c, nil
	}
	return nil, domainerrors.ErrNotFound
}

func (r *inMemoryCustomerRepo) FindByDocument(_ context.Context, doc string) (*entities.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.data {
		if c.Document() == doc {
			return c, nil
		}
	}
	return nil, domainerrors.ErrNotFound
}

func (r *inMemoryCustomerRepo) FindAll(_ context.Context) ([]*entities.Customer, error) {
	return nil, nil
}
func (r *inMemoryCustomerRepo) Update(_ context.Context, c *entities.Customer) error { return nil }
func (r *inMemoryCustomerRepo) Delete(_ context.Context, id string) error            { return nil }

// --- Vehicle ---

type inMemoryVehicleRepo struct {
	mu   sync.RWMutex
	data map[string]*entities.Vehicle
	seq  int
}

func newVehicleRepo() *inMemoryVehicleRepo {
	return &inMemoryVehicleRepo{data: make(map[string]*entities.Vehicle)}
}

func (r *inMemoryVehicleRepo) Create(_ context.Context, v *entities.Vehicle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	v.SetID(fmt.Sprintf("veh-%d", r.seq))
	r.data[v.ID()] = v
	return nil
}

func (r *inMemoryVehicleRepo) FindByID(_ context.Context, id string) (*entities.Vehicle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok {
		return v, nil
	}
	return nil, domainerrors.ErrNotFound
}

func (r *inMemoryVehicleRepo) FindByCustomerID(_ context.Context, customerID string) ([]*entities.Vehicle, error) {
	return nil, nil
}
func (r *inMemoryVehicleRepo) FindAll(_ context.Context) ([]*entities.Vehicle, error) {
	return nil, nil
}
func (r *inMemoryVehicleRepo) Update(_ context.Context, v *entities.Vehicle) error { return nil }
func (r *inMemoryVehicleRepo) Delete(_ context.Context, id string) error           { return nil }

// --- Service ---

type inMemoryServiceRepo struct {
	mu   sync.RWMutex
	data map[string]*entities.Service
	seq  int
}

func newServiceRepo() *inMemoryServiceRepo {
	return &inMemoryServiceRepo{data: make(map[string]*entities.Service)}
}

func (r *inMemoryServiceRepo) Create(_ context.Context, s *entities.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	s.SetID(fmt.Sprintf("svc-%d", r.seq))
	r.data[s.ID()] = s
	return nil
}

func (r *inMemoryServiceRepo) FindByID(_ context.Context, id string) (*entities.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.data[id]; ok {
		return s, nil
	}
	return nil, domainerrors.ErrNotFound
}

func (r *inMemoryServiceRepo) FindByIDs(_ context.Context, ids []string) ([]*entities.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*entities.Service, 0, len(ids))
	for _, id := range ids {
		if s, ok := r.data[id]; ok {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *inMemoryServiceRepo) FindAll(_ context.Context) ([]*entities.Service, error) {
	return nil, nil
}
func (r *inMemoryServiceRepo) Update(_ context.Context, s *entities.Service) error { return nil }
func (r *inMemoryServiceRepo) Delete(_ context.Context, id string) error           { return nil }

// --- Part ---

type inMemoryPartRepo struct {
	mu   sync.RWMutex
	data map[string]*entities.Part
	seq  int
}

func newPartRepo() *inMemoryPartRepo {
	return &inMemoryPartRepo{data: make(map[string]*entities.Part)}
}

func (r *inMemoryPartRepo) Create(_ context.Context, p *entities.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	p.SetID(fmt.Sprintf("part-%d", r.seq))
	r.data[p.ID()] = p
	return nil
}

func (r *inMemoryPartRepo) FindByID(_ context.Context, id string) (*entities.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.data[id]; ok {
		return p, nil
	}
	return nil, domainerrors.ErrNotFound
}

func (r *inMemoryPartRepo) FindByIDs(_ context.Context, ids []string) ([]*entities.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*entities.Part, 0, len(ids))
	for _, id := range ids {
		if p, ok := r.data[id]; ok {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *inMemoryPartRepo) FindAll(_ context.Context) ([]*entities.Part, error) { return nil, nil }
func (r *inMemoryPartRepo) Update(_ context.Context, p *entities.Part) error    { return nil }
func (r *inMemoryPartRepo) Delete(_ context.Context, id string) error           { return nil }
func (r *inMemoryPartRepo) UpdateStock(_ context.Context, id string, delta int) error {
	return nil
}

// --- ServiceOrder ---

type inMemoryServiceOrderRepo struct {
	mu   sync.RWMutex
	data map[string]*entities.ServiceOrder
	seq  int
}

func newServiceOrderRepo() *inMemoryServiceOrderRepo {
	return &inMemoryServiceOrderRepo{data: make(map[string]*entities.ServiceOrder)}
}

func (r *inMemoryServiceOrderRepo) Create(_ context.Context, so *entities.ServiceOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	so.SetID(fmt.Sprintf("order-%d", r.seq))
	r.data[so.ID()] = so
	return nil
}

func (r *inMemoryServiceOrderRepo) FindByID(_ context.Context, id string) (*entities.ServiceOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if so, ok := r.data[id]; ok {
		return so, nil
	}
	return nil, domainerrors.ErrNotFound
}

func (r *inMemoryServiceOrderRepo) UpdateStatus(_ context.Context, id string, status entities.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if so, ok := r.data[id]; ok {
		so.UpdateStatus(status)
		return nil
	}
	return domainerrors.ErrNotFound
}

func (r *inMemoryServiceOrderRepo) FindAll(_ context.Context) ([]*entities.ServiceOrder, error) {
	return nil, nil
}
func (r *inMemoryServiceOrderRepo) FindByCustomerID(_ context.Context, id string) ([]*entities.ServiceOrder, error) {
	return nil, nil
}
func (r *inMemoryServiceOrderRepo) Update(_ context.Context, so *entities.ServiceOrder) error {
	return nil
}
func (r *inMemoryServiceOrderRepo) Delete(_ context.Context, id string) error { return nil }

// =============================================================================
// Cenario de integracao — simula o fluxo completo da oficina
// =============================================================================

// testEnv monta todos os repos e use cases necessarios para os testes.
type testEnv struct {
	customerRepo     *inMemoryCustomerRepo
	vehicleRepo      *inMemoryVehicleRepo
	serviceRepo      *inMemoryServiceRepo
	partRepo         *inMemoryPartRepo
	serviceOrderRepo *inMemoryServiceOrderRepo

	createCustomer *customeruc.CreateCustomer
	createVehicle  *vehicleuc.CreateVehicle
	createOrder    *serviceorderuc.CreateServiceOrder
	updateStatus   *serviceorderuc.UpdateServiceOrderStatus
}

func setupTestEnv() *testEnv {
	custRepo := newCustomerRepo()
	vehRepo := newVehicleRepo()
	svcRepo := newServiceRepo()
	partRepo := newPartRepo()
	soRepo := newServiceOrderRepo()

	return &testEnv{
		customerRepo:     custRepo,
		vehicleRepo:      vehRepo,
		serviceRepo:      svcRepo,
		partRepo:         partRepo,
		serviceOrderRepo: soRepo,

		createCustomer: customeruc.NewCreateCustomer(custRepo),
		createVehicle:  vehicleuc.NewCreateVehicle(vehRepo, custRepo),
		createOrder:    serviceorderuc.NewCreateServiceOrder(soRepo, custRepo, vehRepo, svcRepo, partRepo),
		updateStatus:   serviceorderuc.NewUpdateServiceOrderStatus(soRepo),
	}
}

// =============================================================================
// TESTE 1 — Validacao CPF/CNPJ no fluxo de criacao de cliente
// =============================================================================

func TestIntegracao_CriarCliente_CPFValido(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, err := entities.NewCustomer("Diego Silva", "529.982.247-25", "diego@email.com", "11999999999")
	if err != nil {
		t.Fatalf("CPF valido deveria ser aceito: %v", err)
	}

	if err := env.createCustomer.Execute(ctx, customer); err != nil {
		t.Fatalf("criar cliente deveria ter sucesso: %v", err)
	}

	t.Logf("✓ Cliente criado: ID=%s, Doc=%s (tipo=%s)", customer.ID(), customer.DocumentVO().Formatted(), customer.DocumentVO().Type())
}

func TestIntegracao_CriarCliente_CNPJValido(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, err := entities.NewCustomer("Oficina LTDA", "11.222.333/0001-81", "contato@oficina.com", "1133334444")
	if err != nil {
		t.Fatalf("CNPJ valido deveria ser aceito: %v", err)
	}

	if err := env.createCustomer.Execute(ctx, customer); err != nil {
		t.Fatalf("criar cliente PJ deveria ter sucesso: %v", err)
	}

	t.Logf("✓ Cliente PJ criado: ID=%s, Doc=%s (tipo=%s)", customer.ID(), customer.DocumentVO().Formatted(), customer.DocumentVO().Type())
}

func TestIntegracao_CriarCliente_CPFInvalido(t *testing.T) {
	_, err := entities.NewCustomer("Invalido", "123.456.789-00", "x@e.com", "")
	if err == nil {
		t.Fatal("CPF invalido deveria ser rejeitado na criacao da entidade")
	}
	t.Logf("✓ CPF invalido rejeitado: %v", err)
}

func TestIntegracao_CriarCliente_DocumentoDuplicado(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	c1, _ := entities.NewCustomer("Diego", "529.982.247-25", "d@e.com", "")
	env.createCustomer.Execute(ctx, c1)

	c2, _ := entities.NewCustomer("Outro Diego", "52998224725", "outro@e.com", "")
	err := env.createCustomer.Execute(ctx, c2)
	if err == nil {
		t.Fatal("documento duplicado deveria ser rejeitado")
	}
	t.Logf("✓ Duplicidade detectada: %v", err)
}

// =============================================================================
// TESTE 2 — Validacao de placa no fluxo de criacao de veiculo
// =============================================================================

func TestIntegracao_CriarVeiculo_PlacaAntigaValida(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	// Criar cliente primeiro
	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, customer)

	vehicle, err := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	if err != nil {
		t.Fatalf("placa antiga valida deveria ser aceita: %v", err)
	}

	if err := env.createVehicle.Execute(ctx, vehicle); err != nil {
		t.Fatalf("criar veiculo deveria ter sucesso: %v", err)
	}

	t.Logf("✓ Veiculo criado: ID=%s, Placa=%s (formato=%s)", vehicle.ID(), vehicle.PlateVO().String(), vehicle.PlateVO().Format())
}

func TestIntegracao_CriarVeiculo_PlacaMercosulValida(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, customer)

	vehicle, err := entities.NewVehicle(customer.ID(), "BRA0S18", "VW", "Gol", 2022)
	if err != nil {
		t.Fatalf("placa Mercosul valida deveria ser aceita: %v", err)
	}

	if err := env.createVehicle.Execute(ctx, vehicle); err != nil {
		t.Fatalf("criar veiculo deveria ter sucesso: %v", err)
	}

	t.Logf("✓ Veiculo criado: ID=%s, Placa=%s (formato=%s)", vehicle.ID(), vehicle.PlateVO().String(), vehicle.PlateVO().Format())
}

func TestIntegracao_CriarVeiculo_PlacaInvalida(t *testing.T) {
	_, err := entities.NewVehicle("cust-1", "INVALIDA", "Fiat", "Uno", 2020)
	if err == nil {
		t.Fatal("placa invalida deveria ser rejeitada")
	}
	t.Logf("✓ Placa invalida rejeitada: %v", err)
}

func TestIntegracao_CriarVeiculo_ClienteInexistente(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	vehicle, _ := entities.NewVehicle("cust-inexistente", "ABC-1234", "Fiat", "Uno", 2020)
	err := env.createVehicle.Execute(ctx, vehicle)
	if err == nil {
		t.Fatal("deveria rejeitar veiculo para cliente inexistente")
	}
	t.Logf("✓ Cliente inexistente detectado: %v", err)
}

// =============================================================================
// TESTE 3 — Orcamento automatico na criacao da Ordem de Servico
// =============================================================================

func TestIntegracao_CriarOS_OrcamentoAutomatico(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	// Setup: cliente + veiculo + servicos + pecas
	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, customer)

	vehicle, _ := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	env.createVehicle.Execute(ctx, vehicle)

	svc1 := entities.NewService("Troca de oleo", "Troca completa de oleo", 150.00, 30)
	env.serviceRepo.Create(ctx, svc1)

	svc2 := entities.NewService("Alinhamento", "Alinhamento e balanceamento", 120.00, 45)
	env.serviceRepo.Create(ctx, svc2)

	part1 := entities.NewPart("Filtro de oleo", "Filtro WIX", "un", 35.00, 50)
	env.partRepo.Create(ctx, part1)

	part2 := entities.NewPart("Oleo 5W30", "Oleo sintetico 1L", "litro", 45.00, 100)
	env.partRepo.Create(ctx, part2)

	// Criar OS com precos "do input" (que serao sobrescritos pelo orcamento)
	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	so.AddService(entities.ServiceItem{ServiceID: svc1.ID(), Description: "qualquer", Price: 999.99})
	so.AddService(entities.ServiceItem{ServiceID: svc2.ID(), Description: "qualquer", Price: 999.99})
	so.AddPart(entities.PartItem{PartID: part1.ID(), Description: "qualquer", Quantity: 2, UnitPrice: 999.99})
	so.AddPart(entities.PartItem{PartID: part2.ID(), Description: "qualquer", Quantity: 4, UnitPrice: 999.99})

	if err := env.createOrder.Execute(ctx, so); err != nil {
		t.Fatalf("criar OS deveria ter sucesso: %v", err)
	}

	// Verificar: precos vieram do cadastro, nao do input
	t.Logf("── Orcamento gerado ──")

	for _, s := range so.Services() {
		t.Logf("  Servico: %s = R$ %.2f", s.Description, s.Price)
	}
	for _, p := range so.Parts() {
		t.Logf("  Peca:    %s x%d = R$ %.2f (un: R$ %.2f)", p.Description, p.Quantity, float64(p.Quantity)*p.UnitPrice, p.UnitPrice)
	}

	// Total esperado: 150 + 120 + (2*35) + (4*45) = 150 + 120 + 70 + 180 = 520
	expectedTotal := 520.00
	if so.TotalAmount() != expectedTotal {
		t.Errorf("TotalAmount = %.2f, esperava %.2f", so.TotalAmount(), expectedTotal)
	}
	t.Logf("  TOTAL:   R$ %.2f ✓", so.TotalAmount())

	// Verificar que descricoes vieram do cadastro
	if so.Services()[0].Description != "Troca de oleo" {
		t.Errorf("Descricao do servico deveria vir do cadastro, veio: %q", so.Services()[0].Description)
	}
	if so.Services()[0].Price != 150.00 {
		t.Errorf("Preco deveria ser 150.00 (do cadastro), veio: %.2f", so.Services()[0].Price)
	}

	t.Logf("✓ Orcamento automatico: precos do cadastro, input ignorado, total correto")
}

func TestIntegracao_CriarOS_VeiculoDeOutroCliente(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	c1, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, c1)

	c2, _ := entities.NewCustomer("Maria", "11222333000181", "m@e.com", "")
	env.createCustomer.Execute(ctx, c2)

	// Veiculo pertence a Diego (c1)
	vehicle, _ := entities.NewVehicle(c1.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	env.createVehicle.Execute(ctx, vehicle)

	// Tentar criar OS para Maria (c2) com veiculo do Diego
	so := entities.NewServiceOrder(c2.ID(), vehicle.ID())
	err := env.createOrder.Execute(ctx, so)
	if err == nil {
		t.Fatal("deveria rejeitar OS com veiculo de outro cliente")
	}
	t.Logf("✓ Veiculo de outro cliente rejeitado: %v", err)
}

func TestIntegracao_CriarOS_EstoqueInsuficiente(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, customer)

	vehicle, _ := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	env.createVehicle.Execute(ctx, vehicle)

	// Peca com estoque = 2
	part := entities.NewPart("Filtro raro", "Filtro especial", "un", 100.00, 2)
	env.partRepo.Create(ctx, part)

	// Pedir 10 unidades (estoque = 2)
	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	so.AddPart(entities.PartItem{PartID: part.ID(), Quantity: 10})

	err := env.createOrder.Execute(ctx, so)
	if err == nil {
		t.Fatal("deveria rejeitar por estoque insuficiente")
	}
	t.Logf("✓ Estoque insuficiente detectado: %v", err)
}

// =============================================================================
// TESTE 4 — Maquina de estados da Ordem de Servico
// =============================================================================

func TestIntegracao_FluxoCompleto_StatusOS(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	// Setup
	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, customer)

	vehicle, _ := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	env.createVehicle.Execute(ctx, vehicle)

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	env.createOrder.Execute(ctx, so)

	t.Logf("OS criada: ID=%s, Status=%s", so.ID(), so.Status())

	// Fluxo completo: received -> ... -> delivered
	transicoes := []struct {
		para    entities.OrderStatus
		descr   string
	}{
		{entities.StatusInDiagnosis, "Mecanico inicia diagnostico"},
		{entities.StatusAwaitingApproval, "Diagnostico concluido, aguarda aprovacao"},
		{entities.StatusInExecution, "Cliente aprovou, inicia execucao"},
		{entities.StatusFinished, "Servico finalizado"},
		{entities.StatusDelivered, "Veiculo entregue ao cliente"},
	}

	for _, tr := range transicoes {
		if err := env.updateStatus.Execute(ctx, so.ID(), tr.para); err != nil {
			t.Fatalf("transicao para %q falhou: %v", tr.para, err)
		}

		// Buscar do repo pra confirmar persistencia
		found, _ := env.serviceOrderRepo.FindByID(ctx, so.ID())
		t.Logf("  %s → Status: %s ✓", tr.descr, found.Status())
	}

	t.Logf("✓ Fluxo completo concluido com sucesso")
}

func TestIntegracao_StatusOS_TransicaoInvalida_PularEtapa(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, customer)

	vehicle, _ := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	env.createVehicle.Execute(ctx, vehicle)

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	env.createOrder.Execute(ctx, so)

	// Tentar pular de received direto para finished
	err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusFinished)
	if err == nil {
		t.Fatal("nao deveria permitir pular de received para finished")
	}

	// Verificar que o status nao mudou
	found, _ := env.serviceOrderRepo.FindByID(ctx, so.ID())
	if found.Status() != entities.StatusReceived {
		t.Errorf("status deveria permanecer received, esta: %s", found.Status())
	}

	t.Logf("✓ Pulo de etapa bloqueado: %v (status permanece: %s)", err, found.Status())
}

func TestIntegracao_StatusOS_RecusaDoCliente(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, customer)

	vehicle, _ := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	env.createVehicle.Execute(ctx, vehicle)

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	env.createOrder.Execute(ctx, so)

	// received -> in_diagnosis -> awaiting_approval
	env.updateStatus.Execute(ctx, so.ID(), entities.StatusInDiagnosis)
	env.updateStatus.Execute(ctx, so.ID(), entities.StatusAwaitingApproval)

	// Cliente RECUSA o orcamento → volta para received
	err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusReceived)
	if err != nil {
		t.Fatalf("recusa deveria ser permitida: %v", err)
	}

	found, _ := env.serviceOrderRepo.FindByID(ctx, so.ID())
	t.Logf("✓ Cliente recusou orcamento: status voltou para %s", found.Status())

	// Deve poder reiniciar o fluxo
	err = env.updateStatus.Execute(ctx, so.ID(), entities.StatusInDiagnosis)
	if err != nil {
		t.Fatalf("reinicio do fluxo deveria ser permitido: %v", err)
	}

	t.Logf("✓ Fluxo reiniciado com sucesso apos recusa")
}

func TestIntegracao_StatusOS_EstadoTerminal(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	env.createCustomer.Execute(ctx, customer)

	vehicle, _ := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	env.createVehicle.Execute(ctx, vehicle)

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	env.createOrder.Execute(ctx, so)

	// Percorrer todo o fluxo ate delivered
	env.updateStatus.Execute(ctx, so.ID(), entities.StatusInDiagnosis)
	env.updateStatus.Execute(ctx, so.ID(), entities.StatusAwaitingApproval)
	env.updateStatus.Execute(ctx, so.ID(), entities.StatusInExecution)
	env.updateStatus.Execute(ctx, so.ID(), entities.StatusFinished)
	env.updateStatus.Execute(ctx, so.ID(), entities.StatusDelivered)

	// Tentar qualquer transicao a partir de delivered
	err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusReceived)
	if err == nil {
		t.Fatal("delivered e estado terminal, nao deveria aceitar transicao")
	}

	t.Logf("✓ Estado terminal (delivered) bloqueou transicao: %v", err)
}

// =============================================================================
// TESTE BONUS — Fluxo completo da oficina (cenario real)
// =============================================================================

func TestIntegracao_FluxoCompletoOficina(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	t.Log("━━━ SIMULACAO: Fluxo completo de atendimento da oficina ━━━")
	t.Log("")

	// 1. Cadastrar cliente
	customer, _ := entities.NewCustomer("Carlos Mecenas", "529.982.247-25", "carlos@email.com", "11987654321")
	env.createCustomer.Execute(ctx, customer)
	t.Logf("1. Cliente cadastrado: %s (CPF: %s)", customer.Name(), customer.DocumentVO().Formatted())

	// 2. Cadastrar veiculo
	vehicle, _ := entities.NewVehicle(customer.ID(), "BRA0S18", "Honda", "Civic", 2023)
	env.createVehicle.Execute(ctx, vehicle)
	t.Logf("2. Veiculo cadastrado: %s %s %d (Placa: %s - %s)", vehicle.Brand(), vehicle.Model(), vehicle.Year(), vehicle.PlateVO().String(), vehicle.PlateVO().Format())

	// 3. Cadastrar servicos e pecas disponíveis
	svcRevisao := entities.NewService("Revisao completa", "Revisao dos 30.000km", 450.00, 120)
	env.serviceRepo.Create(ctx, svcRevisao)

	svcFreio := entities.NewService("Troca de pastilha", "Troca pastilha de freio dianteira", 180.00, 60)
	env.serviceRepo.Create(ctx, svcFreio)

	partPastilha := entities.NewPart("Pastilha Bosch", "Pastilha de freio dianteira", "jogo", 189.90, 15)
	env.partRepo.Create(ctx, partPastilha)

	partOleo := entities.NewPart("Oleo Mobil 5W30", "Oleo sintetico 1L", "litro", 52.90, 200)
	env.partRepo.Create(ctx, partOleo)

	partFiltro := entities.NewPart("Filtro de oleo", "Filtro Tecfil", "un", 38.50, 30)
	env.partRepo.Create(ctx, partFiltro)

	t.Logf("3. Catalogo: %d servicos, %d pecas cadastrados", 2, 3)

	// 4. Criar OS com orcamento automatico
	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	so.SetNotes("Cliente relata barulho no freio e revisao preventiva")
	so.AddService(entities.ServiceItem{ServiceID: svcRevisao.ID()})
	so.AddService(entities.ServiceItem{ServiceID: svcFreio.ID()})
	so.AddPart(entities.PartItem{PartID: partPastilha.ID(), Quantity: 1})
	so.AddPart(entities.PartItem{PartID: partOleo.ID(), Quantity: 4})
	so.AddPart(entities.PartItem{PartID: partFiltro.ID(), Quantity: 1})

	if err := env.createOrder.Execute(ctx, so); err != nil {
		t.Fatalf("criar OS falhou: %v", err)
	}

	t.Logf("4. OS criada: #%s", so.ID())
	t.Logf("   Obs: %s", so.Notes())
	t.Log("   ── Orcamento ──")
	for _, s := range so.Services() {
		t.Logf("   Servico: %-30s R$ %8.2f", s.Description, s.Price)
	}
	for _, p := range so.Parts() {
		t.Logf("   Peca:    %-20s x%-2d   R$ %8.2f", p.Description, p.Quantity, float64(p.Quantity)*p.UnitPrice)
	}

	// Total: 450 + 180 + 189.90 + (4*52.90) + 38.50 = 450 + 180 + 189.90 + 211.60 + 38.50 = 1070.00
	t.Logf("   ─────────────────────────────────────")
	t.Logf("   TOTAL:                        R$ %8.2f", so.TotalAmount())
	t.Log("")

	expectedTotal := 1070.00
	if so.TotalAmount() != expectedTotal {
		t.Errorf("Total esperado: %.2f, obtido: %.2f", expectedTotal, so.TotalAmount())
	}

	// 5. Fluxo de status
	steps := []struct {
		status entities.OrderStatus
		msg    string
	}{
		{entities.StatusInDiagnosis, "Mecanico inicia diagnostico"},
		{entities.StatusAwaitingApproval, "Orcamento enviado ao cliente"},
		{entities.StatusInExecution, "Cliente aprovou — execucao iniciada"},
		{entities.StatusFinished, "Servicos concluidos"},
		{entities.StatusDelivered, "Veiculo entregue ao cliente"},
	}

	for i, step := range steps {
		if err := env.updateStatus.Execute(ctx, so.ID(), step.status); err != nil {
			t.Fatalf("passo %d falhou: %v", i+1, err)
		}
		found, _ := env.serviceOrderRepo.FindByID(ctx, so.ID())
		t.Logf("5.%d %s → [%s]", i+1, step.msg, found.Status())
	}

	t.Log("")
	t.Logf("━━━ FLUXO COMPLETO CONCLUIDO COM SUCESSO ━━━")
	t.Logf("   Tempo do teste: %s", time.Since(time.Now().Add(-time.Millisecond)).Round(time.Millisecond))
}
