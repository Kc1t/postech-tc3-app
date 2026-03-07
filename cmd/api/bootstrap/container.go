package bootstrap

import (
	"context"
	"log"
	"sync"

	"github.com/fiap/postech-tc1/config"
	handler "github.com/fiap/postech-tc1/internal/adapters/inbound/http"
	"github.com/fiap/postech-tc1/internal/adapters/outbound/mongodb"
	"github.com/fiap/postech-tc1/internal/application/usecase"
	"github.com/fiap/postech-tc1/internal/ports"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TODO: Arrumar estrutura de injecao de dependencia

// Container centraliza o acesso as dependencias da aplicacao.
// A tag `container` documenta a camada de cada dependencia.
type Container struct {
	Config *config.Config
	db     *mongodriver.Database

	// Repositories — outbound adapters (persistencia)
	CustomerRepo     ports.CustomerRepository     `container:"repository"`
	VehicleRepo      ports.VehicleRepository      `container:"repository"`
	ServiceOrderRepo ports.ServiceOrderRepository `container:"repository"`
	ServiceRepo      ports.ServiceRepository      `container:"repository"`
	PartRepo         ports.PartRepository         `container:"repository"`

	// Use Cases — logica de negocio
	CustomerUseCase     ports.CustomerUseCase     `container:"usecase"`
	VehicleUseCase      ports.VehicleUseCase      `container:"usecase"`
	ServiceOrderUseCase ports.ServiceOrderUseCase `container:"usecase"`
	ServiceUseCase      ports.ServiceUseCase      `container:"usecase"`
	PartUseCase         ports.PartUseCase         `container:"usecase"`

	// Handlers — inbound adapters (HTTP)
	CustomerHandler     *handler.CustomerHandler     `container:"handler"`
	VehicleHandler      *handler.VehicleHandler      `container:"handler"`
	ServiceOrderHandler *handler.ServiceOrderHandler `container:"handler"`
	ServiceHandler      *handler.ServiceHandler      `container:"handler"`
	PartHandler         *handler.PartHandler         `container:"handler"`
}

var (
	instance *Container
	once     sync.Once
)

// GetContainer retorna a instancia unica do Container (Singleton).
func GetContainer() *Container {
	once.Do(func() {
		instance = &Container{}
		instance.initialize()
	})
	return instance
}

func (c *Container) initialize() {
	c.Config = config.Load()
	c.setupDatabase()
	c.setupRepositories()
	c.setupUseCases()
	c.setupHandlers()
}

func (c *Container) setupDatabase() {
	ctx := context.Background()
	client, err := mongodriver.Connect(ctx, options.Client().ApplyURI(c.Config.MongoURI))
	if err != nil {
		log.Fatalf("bootstrap: failed to connect to mongodb: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("bootstrap: failed to ping mongodb: %v", err)
	}
	c.db = client.Database(c.Config.MongoDB)
	log.Printf("bootstrap: connected to mongodb database=%s", c.Config.MongoDB)
}

func (c *Container) setupRepositories() {
	c.CustomerRepo = mongodb.NewCustomerRepository(c.db)
	c.VehicleRepo = mongodb.NewVehicleRepository(c.db)
	c.ServiceOrderRepo = mongodb.NewServiceOrderRepository(c.db)
	c.ServiceRepo = mongodb.NewServiceRepository(c.db)
	c.PartRepo = mongodb.NewPartRepository(c.db)
}

func (c *Container) setupUseCases() {
	c.CustomerUseCase = usecase.NewCustomerUseCase(c.CustomerRepo)
	c.VehicleUseCase = usecase.NewVehicleUseCase(c.VehicleRepo, c.CustomerRepo)
	c.ServiceOrderUseCase = usecase.NewServiceOrderUseCase(
		c.ServiceOrderRepo,
		c.CustomerRepo,
		c.VehicleRepo,
		c.ServiceRepo,
		c.PartRepo,
	)
	c.ServiceUseCase = usecase.NewServiceUseCase(c.ServiceRepo)
	c.PartUseCase = usecase.NewPartUseCase(c.PartRepo)
}

func (c *Container) setupHandlers() {
	c.CustomerHandler = handler.NewCustomerHandler(c.CustomerUseCase)
	c.VehicleHandler = handler.NewVehicleHandler(c.VehicleUseCase)
	c.ServiceOrderHandler = handler.NewServiceOrderHandler(c.ServiceOrderUseCase)
	c.ServiceHandler = handler.NewServiceHandler(c.ServiceUseCase)
	c.PartHandler = handler.NewPartHandler(c.PartUseCase)
}
