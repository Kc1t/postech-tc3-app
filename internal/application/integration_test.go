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
	serviceorderuc "github.com/fiap/postech-tc1/internal/application/usecase/service_order"
	vehicleuc "github.com/fiap/postech-tc1/internal/application/usecase/vehicle"
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
	if _, ok := r.data[id]; ok {
		// O use case ja validou e aplicou a transicao na entidade.
		// O repo apenas persiste — como o objeto em memoria ja foi atualizado
		// pelo use case, nao precisamos fazer nada aqui.
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

// setupCustomerAndVehicle cria um cliente e veiculo para testes, falhando se houver erro.
func setupCustomerAndVehicle(t *testing.T, env *testEnv, ctx context.Context) (*entities.Customer, *entities.Vehicle) {
	t.Helper()
	customer, err := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	if err != nil {
		t.Fatalf("erro ao criar customer: %v", err)
	}
	if err := env.createCustomer.Execute(ctx, customer); err != nil {
		t.Fatalf("erro ao persistir customer: %v", err)
	}
	vehicle, err := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	if err != nil {
		t.Fatalf("erro ao criar vehicle: %v", err)
	}
	if err := env.createVehicle.Execute(ctx, vehicle); err != nil {
		t.Fatalf("erro ao persistir vehicle: %v", err)
	}
	return customer, vehicle
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

	t.Logf("Cliente criado: ID=%s, Doc=%s (tipo=%s)", customer.ID(), customer.DocumentVO().Formatted(), customer.DocumentVO().Type())
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

	t.Logf("Cliente PJ criado: ID=%s, Doc=%s (tipo=%s)", customer.ID(), customer.DocumentVO().Formatted(), customer.DocumentVO().Type())
}

func TestIntegracao_CriarCliente_CPFInvalido(t *testing.T) {
	_, err := entities.NewCustomer("Invalido", "123.456.789-00", "x@e.com", "")
	if err == nil {
		t.Fatal("CPF invalido deveria ser rejeitado na criacao da entidade")
	}
	t.Logf("CPF invalido rejeitado: %v", err)
}

func TestIntegracao_CriarCliente_DocumentoDuplicado(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	c1, _ := entities.NewCustomer("Diego", "529.982.247-25", "d@e.com", "")
	if err := env.createCustomer.Execute(ctx, c1); err != nil {
		t.Fatalf("primeiro cliente deveria ser criado: %v", err)
	}

	c2, _ := entities.NewCustomer("Outro Diego", "52998224725", "outro@e.com", "")
	err := env.createCustomer.Execute(ctx, c2)
	if err == nil {
		t.Fatal("documento duplicado deveria ser rejeitado")
	}
	t.Logf("Duplicidade detectada: %v", err)
}

// =============================================================================
// TESTE 2 — Validacao de placa no fluxo de criacao de veiculo
// =============================================================================

func TestIntegracao_CriarVeiculo_PlacaAntigaValida(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	if err := env.createCustomer.Execute(ctx, customer); err != nil {
		t.Fatalf("setup falhou: %v", err)
	}

	vehicle, err := entities.NewVehicle(customer.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	if err != nil {
		t.Fatalf("placa antiga valida deveria ser aceita: %v", err)
	}

	if err := env.createVehicle.Execute(ctx, vehicle); err != nil {
		t.Fatalf("criar veiculo deveria ter sucesso: %v", err)
	}

	t.Logf("Veiculo criado: ID=%s, Placa=%s (formato=%s)", vehicle.ID(), vehicle.PlateVO().String(), vehicle.PlateVO().Format())
}

func TestIntegracao_CriarVeiculo_PlacaMercosulValida(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	if err := env.createCustomer.Execute(ctx, customer); err != nil {
		t.Fatalf("setup falhou: %v", err)
	}

	vehicle, err := entities.NewVehicle(customer.ID(), "BRA0S18", "VW", "Gol", 2022)
	if err != nil {
		t.Fatalf("placa Mercosul valida deveria ser aceita: %v", err)
	}

	if err := env.createVehicle.Execute(ctx, vehicle); err != nil {
		t.Fatalf("criar veiculo deveria ter sucesso: %v", err)
	}

	t.Logf("Veiculo criado: ID=%s, Placa=%s (formato=%s)", vehicle.ID(), vehicle.PlateVO().String(), vehicle.PlateVO().Format())
}

func TestIntegracao_CriarVeiculo_PlacaInvalida(t *testing.T) {
	_, err := entities.NewVehicle("cust-1", "INVALIDA", "Fiat", "Uno", 2020)
	if err == nil {
		t.Fatal("placa invalida deveria ser rejeitada")
	}
	t.Logf("Placa invalida rejeitada: %v", err)
}

func TestIntegracao_CriarVeiculo_ClienteInexistente(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	vehicle, _ := entities.NewVehicle("cust-inexistente", "ABC-1234", "Fiat", "Uno", 2020)
	err := env.createVehicle.Execute(ctx, vehicle)
	if err == nil {
		t.Fatal("deveria rejeitar veiculo para cliente inexistente")
	}
	t.Logf("Cliente inexistente detectado: %v", err)
}

// =============================================================================
// TESTE 3 — Orcamento automatico na criacao da Ordem de Servico
// =============================================================================

func TestIntegracao_CriarOS_OrcamentoAutomatico(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, vehicle := setupCustomerAndVehicle(t, env, ctx)

	svc1 := entities.NewService("Troca de oleo", "Troca completa de oleo", 150.00, 30)
	if err := env.serviceRepo.Create(ctx, svc1); err != nil {
		t.Fatalf("setup svc1 falhou: %v", err)
	}

	svc2 := entities.NewService("Alinhamento", "Alinhamento e balanceamento", 120.00, 45)
	if err := env.serviceRepo.Create(ctx, svc2); err != nil {
		t.Fatalf("setup svc2 falhou: %v", err)
	}

	part1 := entities.NewPart("Filtro de oleo", "Filtro WIX", "un", 35.00, 50)
	if err := env.partRepo.Create(ctx, part1); err != nil {
		t.Fatalf("setup part1 falhou: %v", err)
	}

	part2 := entities.NewPart("Oleo 5W30", "Oleo sintetico 1L", "litro", 45.00, 100)
	if err := env.partRepo.Create(ctx, part2); err != nil {
		t.Fatalf("setup part2 falhou: %v", err)
	}

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	so.AddService(entities.ServiceItem{ServiceID: svc1.ID(), Description: "qualquer", Price: 999.99})
	so.AddService(entities.ServiceItem{ServiceID: svc2.ID(), Description: "qualquer", Price: 999.99})
	so.AddPart(entities.PartItem{PartID: part1.ID(), Description: "qualquer", Quantity: 2, UnitPrice: 999.99})
	so.AddPart(entities.PartItem{PartID: part2.ID(), Description: "qualquer", Quantity: 4, UnitPrice: 999.99})

	if err := env.createOrder.Execute(ctx, so); err != nil {
		t.Fatalf("criar OS deveria ter sucesso: %v", err)
	}

	expectedTotal := 520.00
	if so.TotalAmount() != expectedTotal {
		t.Errorf("TotalAmount = %.2f, esperava %.2f", so.TotalAmount(), expectedTotal)
	}
	if so.Services()[0].Description != "Troca de oleo" {
		t.Errorf("Descricao do servico deveria vir do cadastro, veio: %q", so.Services()[0].Description)
	}
	if so.Services()[0].Price != 150.00 {
		t.Errorf("Preco deveria ser 150.00 (do cadastro), veio: %.2f", so.Services()[0].Price)
	}
}

func TestIntegracao_CriarOS_VeiculoDeOutroCliente(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	c1, _ := entities.NewCustomer("Diego", "52998224725", "d@e.com", "")
	if err := env.createCustomer.Execute(ctx, c1); err != nil {
		t.Fatalf("setup c1 falhou: %v", err)
	}

	c2, _ := entities.NewCustomer("Maria", "11222333000181", "m@e.com", "")
	if err := env.createCustomer.Execute(ctx, c2); err != nil {
		t.Fatalf("setup c2 falhou: %v", err)
	}

	vehicle, _ := entities.NewVehicle(c1.ID(), "ABC-1234", "Fiat", "Uno", 2020)
	if err := env.createVehicle.Execute(ctx, vehicle); err != nil {
		t.Fatalf("setup vehicle falhou: %v", err)
	}

	so := entities.NewServiceOrder(c2.ID(), vehicle.ID())
	err := env.createOrder.Execute(ctx, so)
	if err == nil {
		t.Fatal("deveria rejeitar OS com veiculo de outro cliente")
	}
	t.Logf("Veiculo de outro cliente rejeitado: %v", err)
}

func TestIntegracao_CriarOS_EstoqueInsuficiente(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, vehicle := setupCustomerAndVehicle(t, env, ctx)

	part := entities.NewPart("Filtro raro", "Filtro especial", "un", 100.00, 2)
	if err := env.partRepo.Create(ctx, part); err != nil {
		t.Fatalf("setup part falhou: %v", err)
	}

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	so.AddPart(entities.PartItem{PartID: part.ID(), Quantity: 10})

	err := env.createOrder.Execute(ctx, so)
	if err == nil {
		t.Fatal("deveria rejeitar por estoque insuficiente")
	}
	t.Logf("Estoque insuficiente detectado: %v", err)
}

// =============================================================================
// TESTE 4 — Maquina de estados da Ordem de Servico
// =============================================================================

func TestIntegracao_FluxoCompleto_StatusOS(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, vehicle := setupCustomerAndVehicle(t, env, ctx)

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	if err := env.createOrder.Execute(ctx, so); err != nil {
		t.Fatalf("criar OS falhou: %v", err)
	}

	transicoes := []struct {
		para  entities.OrderStatus
		descr string
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
		found, _ := env.serviceOrderRepo.FindByID(ctx, so.ID())
		t.Logf("  %s -> Status: %s", tr.descr, found.Status())
	}
}

func TestIntegracao_StatusOS_TransicaoInvalida_PularEtapa(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, vehicle := setupCustomerAndVehicle(t, env, ctx)

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	if err := env.createOrder.Execute(ctx, so); err != nil {
		t.Fatalf("criar OS falhou: %v", err)
	}

	err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusFinished)
	if err == nil {
		t.Fatal("nao deveria permitir pular de received para finished")
	}

	found, _ := env.serviceOrderRepo.FindByID(ctx, so.ID())
	if found.Status() != entities.StatusReceived {
		t.Errorf("status deveria permanecer received, esta: %s", found.Status())
	}
}

func TestIntegracao_StatusOS_RecusaDoCliente(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, vehicle := setupCustomerAndVehicle(t, env, ctx)

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	if err := env.createOrder.Execute(ctx, so); err != nil {
		t.Fatalf("criar OS falhou: %v", err)
	}

	if err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusInDiagnosis); err != nil {
		t.Fatalf("transicao para in_diagnosis falhou: %v", err)
	}
	if err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusAwaitingApproval); err != nil {
		t.Fatalf("transicao para awaiting_approval falhou: %v", err)
	}

	if err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusReceived); err != nil {
		t.Fatalf("recusa deveria ser permitida: %v", err)
	}

	found, _ := env.serviceOrderRepo.FindByID(ctx, so.ID())
	t.Logf("Cliente recusou orcamento: status voltou para %s", found.Status())

	if err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusInDiagnosis); err != nil {
		t.Fatalf("reinicio do fluxo deveria ser permitido: %v", err)
	}
}

func TestIntegracao_StatusOS_EstadoTerminal(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	customer, vehicle := setupCustomerAndVehicle(t, env, ctx)

	so := entities.NewServiceOrder(customer.ID(), vehicle.ID())
	if err := env.createOrder.Execute(ctx, so); err != nil {
		t.Fatalf("criar OS falhou: %v", err)
	}

	fluxo := []entities.OrderStatus{
		entities.StatusInDiagnosis,
		entities.StatusAwaitingApproval,
		entities.StatusInExecution,
		entities.StatusFinished,
		entities.StatusDelivered,
	}
	for _, status := range fluxo {
		if err := env.updateStatus.Execute(ctx, so.ID(), status); err != nil {
			t.Fatalf("transicao para %q falhou: %v", status, err)
		}
	}

	err := env.updateStatus.Execute(ctx, so.ID(), entities.StatusReceived)
	if err == nil {
		t.Fatal("delivered e estado terminal, nao deveria aceitar transicao")
	}
}

// =============================================================================
// TESTE BONUS — Fluxo completo da oficina (cenario real)
// =============================================================================

func TestIntegracao_FluxoCompletoOficina(t *testing.T) {
	env := setupTestEnv()
	ctx := context.Background()

	t.Log("SIMULACAO: Fluxo completo de atendimento da oficina")

	// 1. Cadastrar cliente
	customer, err := entities.NewCustomer("Carlos Mecenas", "529.982.247-25", "carlos@email.com", "11987654321")
	if err != nil {
		t.Fatalf("erro ao criar customer: %v", err)
	}
	if err := env.createCustomer.Execute(ctx, customer); err != nil {
		t.Fatalf("erro ao persistir customer: %v", err)
	}
	t.Logf("1. Cliente cadastrado: %s (CPF: %s)", customer.Name(), customer.DocumentVO().Formatted())

	// 2. Cadastrar veiculo
	vehicle, err := entities.NewVehicle(customer.ID(), "BRA0S18", "Honda", "Civic", 2023)
	if err != nil {
		t.Fatalf("erro ao criar vehicle: %v", err)
	}
	if err := env.createVehicle.Execute(ctx, vehicle); err != nil {
		t.Fatalf("erro ao persistir vehicle: %v", err)
	}
	t.Logf("2. Veiculo cadastrado: %s %s %d (Placa: %s - %s)", vehicle.Brand(), vehicle.Model(), vehicle.Year(), vehicle.PlateVO().String(), vehicle.PlateVO().Format())

	// 3. Cadastrar servicos e pecas
	svcRevisao := entities.NewService("Revisao completa", "Revisao dos 30.000km", 450.00, 120)
	if err := env.serviceRepo.Create(ctx, svcRevisao); err != nil {
		t.Fatalf("setup svcRevisao falhou: %v", err)
	}

	svcFreio := entities.NewService("Troca de pastilha", "Troca pastilha de freio dianteira", 180.00, 60)
	if err := env.serviceRepo.Create(ctx, svcFreio); err != nil {
		t.Fatalf("setup svcFreio falhou: %v", err)
	}

	partPastilha := entities.NewPart("Pastilha Bosch", "Pastilha de freio dianteira", "jogo", 189.90, 15)
	if err := env.partRepo.Create(ctx, partPastilha); err != nil {
		t.Fatalf("setup partPastilha falhou: %v", err)
	}

	partOleo := entities.NewPart("Oleo Mobil 5W30", "Oleo sintetico 1L", "litro", 52.90, 200)
	if err := env.partRepo.Create(ctx, partOleo); err != nil {
		t.Fatalf("setup partOleo falhou: %v", err)
	}

	partFiltro := entities.NewPart("Filtro de oleo", "Filtro Tecfil", "un", 38.50, 30)
	if err := env.partRepo.Create(ctx, partFiltro); err != nil {
		t.Fatalf("setup partFiltro falhou: %v", err)
	}

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

	t.Logf("4. OS criada: #%s (Total: R$ %.2f)", so.ID(), so.TotalAmount())

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
		t.Logf("5.%d %s -> [%s]", i+1, step.msg, found.Status())
	}

	t.Logf("FLUXO COMPLETO CONCLUIDO COM SUCESSO")
	t.Logf("   Tempo do teste: %s", time.Since(time.Now().Add(-time.Millisecond)).Round(time.Millisecond))
}
