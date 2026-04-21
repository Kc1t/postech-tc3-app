package bootstrap

import (
	"log"
	"sync"

	"github.com/fiap/postech-tc1/config"
	"github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql"
	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	"github.com/fiap/postech-tc1/internal/ports"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	customerhandler "github.com/fiap/postech-tc1/internal/adapters/inbound/http/customer"
	parthandler "github.com/fiap/postech-tc1/internal/adapters/inbound/http/part"
	servicehandler "github.com/fiap/postech-tc1/internal/adapters/inbound/http/service"
	serviceorderhandler "github.com/fiap/postech-tc1/internal/adapters/inbound/http/service_order"
	vehiclehandler "github.com/fiap/postech-tc1/internal/adapters/inbound/http/vehicle"

	customeruc "github.com/fiap/postech-tc1/internal/application/usecase/customer"
	partuc "github.com/fiap/postech-tc1/internal/application/usecase/part"
	serviceuc "github.com/fiap/postech-tc1/internal/application/usecase/service"
	serviceorderuc "github.com/fiap/postech-tc1/internal/application/usecase/service_order"
	vehicleuc "github.com/fiap/postech-tc1/internal/application/usecase/vehicle"
)

// Container centraliza o acesso as dependencias da aplicacao.
// A tag `container` documenta a camada de cada dependencia.
type Container struct {
	Config *config.Config
	db     *gorm.DB

	// Repositories — outbound adapters (persistencia)
	CustomerRepo     ports.CustomerRepository     `container:"repository"`
	VehicleRepo      ports.VehicleRepository      `container:"repository"`
	ServiceOrderRepo ports.ServiceOrderRepository `container:"repository"`
	ServiceRepo      ports.ServiceRepository      `container:"repository"`
	PartRepo         ports.PartRepository         `container:"repository"`

	// Use Cases — customer
	CreateCustomer        ports.CreateCustomerUseCase        `container:"usecase"`
	GetCustomer           ports.GetCustomerUseCase           `container:"usecase"`
	GetCustomerByDocument ports.GetCustomerByDocumentUseCase `container:"usecase"`
	ListCustomers         ports.ListCustomersUseCase         `container:"usecase"`
	UpdateCustomer        ports.UpdateCustomerUseCase        `container:"usecase"`
	DeleteCustomer        ports.DeleteCustomerUseCase        `container:"usecase"`

	// Use Cases — vehicle
	CreateVehicle          ports.CreateVehicleUseCase          `container:"usecase"`
	GetVehicle             ports.GetVehicleUseCase             `container:"usecase"`
	ListVehicles           ports.ListVehiclesUseCase           `container:"usecase"`
	ListVehiclesByCustomer ports.ListVehiclesByCustomerUseCase `container:"usecase"`
	UpdateVehicle          ports.UpdateVehicleUseCase          `container:"usecase"`
	DeleteVehicle          ports.DeleteVehicleUseCase          `container:"usecase"`

	// Use Cases — service order
	CreateServiceOrder          ports.CreateServiceOrderUseCase          `container:"usecase"`
	GetServiceOrder             ports.GetServiceOrderUseCase             `container:"usecase"`
	GetServiceOrderByCode       ports.GetServiceOrderByCodeUseCase       `container:"usecase"`
	ListServiceOrders           ports.ListServiceOrdersUseCase           `container:"usecase"`
	ListServiceOrdersByCustomer ports.ListServiceOrdersByCustomerUseCase `container:"usecase"`
	ListServiceOrdersByDocument ports.ListServiceOrdersByDocumentUseCase `container:"usecase"`
	UpdateServiceOrderStatus    ports.UpdateServiceOrderStatusUseCase    `container:"usecase"`
	UpdateServiceOrder          ports.UpdateServiceOrderUseCase          `container:"usecase"`
	DeleteServiceOrder          ports.DeleteServiceOrderUseCase          `container:"usecase"`
	UpdateServiceOrderStatusByCode ports.UpdateServiceOrderStatusByCodeUseCase `container:"usecase"`
	GetAverageExecutionTime        ports.GetAverageExecutionTimeUseCase        `container:"usecase"`

	// Use Cases — service
	CreateService ports.CreateServiceUseCase `container:"usecase"`
	GetService    ports.GetServiceUseCase    `container:"usecase"`
	ListServices  ports.ListServicesUseCase  `container:"usecase"`
	UpdateService ports.UpdateServiceUseCase `container:"usecase"`
	DeleteService ports.DeleteServiceUseCase `container:"usecase"`

	// Use Cases — part
	CreatePart      ports.CreatePartUseCase      `container:"usecase"`
	GetPart         ports.GetPartUseCase         `container:"usecase"`
	ListParts       ports.ListPartsUseCase       `container:"usecase"`
	UpdatePart      ports.UpdatePartUseCase      `container:"usecase"`
	DeletePart      ports.DeletePartUseCase      `container:"usecase"`
	AdjustPartStock ports.AdjustPartStockUseCase `container:"usecase"`

	// Handlers — inbound adapters (HTTP)
	CustomerHandler     *customerhandler.CustomerHandler         `container:"handler"`
	VehicleHandler      *vehiclehandler.VehicleHandler           `container:"handler"`
	ServiceOrderHandler *serviceorderhandler.ServiceOrderHandler `container:"handler"`
	ServiceHandler      *servicehandler.ServiceHandler           `container:"handler"`
	PartHandler         *parthandler.PartHandler                 `container:"handler"`
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
	_ = godotenv.Load()
	c.Config = config.Load()
	c.setupDatabase()
	c.setupRepositories()
	c.setupUseCases()
	c.setupHandlers()
}

func (c *Container) setupDatabase() {
	db, err := gorm.Open(postgres.Open(c.Config.PostgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("bootstrap: failed to connect to postgresql: %v", err)
	}

	if err := db.AutoMigrate(
		&pgmodel.Customer{},
		&pgmodel.Vehicle{},
		&pgmodel.Service{},
		&pgmodel.Part{},
		&pgmodel.ServiceOrder{},
	); err != nil {
		log.Fatalf("bootstrap: failed to run migrations: %v", err)
	}

	c.db = db
	log.Printf("bootstrap: connected to postgresql and migrations applied")
}

func (c *Container) setupRepositories() {
	c.CustomerRepo = postgresql.NewCustomerRepository(c.db)
	c.VehicleRepo = postgresql.NewVehicleRepository(c.db)
	c.ServiceOrderRepo = postgresql.NewServiceOrderRepository(c.db)
	c.ServiceRepo = postgresql.NewServiceRepository(c.db)
	c.PartRepo = postgresql.NewPartRepository(c.db)
}

func (c *Container) setupUseCases() {
	// Customer
	c.CreateCustomer = customeruc.NewCreateCustomer(c.CustomerRepo)
	c.GetCustomer = customeruc.NewGetCustomer(c.CustomerRepo)
	c.GetCustomerByDocument = customeruc.NewGetCustomerByDocument(c.CustomerRepo)
	c.ListCustomers = customeruc.NewListCustomers(c.CustomerRepo)
	c.UpdateCustomer = customeruc.NewUpdateCustomer(c.CustomerRepo)
	c.DeleteCustomer = customeruc.NewDeleteCustomer(c.CustomerRepo)

	// Vehicle
	c.CreateVehicle = vehicleuc.NewCreateVehicle(c.VehicleRepo, c.CustomerRepo)
	c.GetVehicle = vehicleuc.NewGetVehicle(c.VehicleRepo)
	c.ListVehicles = vehicleuc.NewListVehicles(c.VehicleRepo)
	c.ListVehiclesByCustomer = vehicleuc.NewListVehiclesByCustomer(c.VehicleRepo)
	c.UpdateVehicle = vehicleuc.NewUpdateVehicle(c.VehicleRepo)
	c.DeleteVehicle = vehicleuc.NewDeleteVehicle(c.VehicleRepo)

	// ServiceOrder
	c.CreateServiceOrder = serviceorderuc.NewCreateServiceOrder(
		c.ServiceOrderRepo, c.CustomerRepo, c.VehicleRepo,
	)
	c.GetServiceOrder = serviceorderuc.NewGetServiceOrder(c.ServiceOrderRepo)
	c.GetServiceOrderByCode = serviceorderuc.NewGetServiceOrderByCode(c.ServiceOrderRepo, c.CustomerRepo)
	c.ListServiceOrders = serviceorderuc.NewListServiceOrders(c.ServiceOrderRepo)
	c.ListServiceOrdersByCustomer = serviceorderuc.NewListServiceOrdersByCustomer(c.ServiceOrderRepo)
	c.ListServiceOrdersByDocument = serviceorderuc.NewListServiceOrdersByDocument(c.ServiceOrderRepo, c.CustomerRepo)
	c.UpdateServiceOrderStatus = serviceorderuc.NewUpdateServiceOrderStatus(
		c.ServiceOrderRepo, c.ServiceRepo, c.PartRepo,
	)
	c.UpdateServiceOrder = serviceorderuc.NewUpdateServiceOrder(c.ServiceOrderRepo)
	c.DeleteServiceOrder = serviceorderuc.NewDeleteServiceOrder(c.ServiceOrderRepo)
	c.UpdateServiceOrderStatusByCode = serviceorderuc.NewUpdateServiceOrderStatusByCode(c.ServiceOrderRepo, c.CustomerRepo, c.PartRepo)
	c.GetAverageExecutionTime = serviceorderuc.NewGetAverageExecutionTime(c.ServiceOrderRepo)

	// Service
	c.CreateService = serviceuc.NewCreateService(c.ServiceRepo)
	c.GetService = serviceuc.NewGetService(c.ServiceRepo)
	c.ListServices = serviceuc.NewListServices(c.ServiceRepo)
	c.UpdateService = serviceuc.NewUpdateService(c.ServiceRepo)
	c.DeleteService = serviceuc.NewDeleteService(c.ServiceRepo)

	// Part
	c.CreatePart = partuc.NewCreatePart(c.PartRepo)
	c.GetPart = partuc.NewGetPart(c.PartRepo)
	c.ListParts = partuc.NewListParts(c.PartRepo)
	c.UpdatePart = partuc.NewUpdatePart(c.PartRepo)
	c.DeletePart = partuc.NewDeletePart(c.PartRepo)
	c.AdjustPartStock = partuc.NewAdjustPartStock(c.PartRepo)
}

func (c *Container) setupHandlers() {
	c.CustomerHandler = customerhandler.NewCustomerHandler(
		c.CreateCustomer, c.GetCustomer, c.GetCustomerByDocument,
		c.ListCustomers, c.UpdateCustomer, c.DeleteCustomer,
		c.ListVehiclesByCustomer,
	)
	c.VehicleHandler = vehiclehandler.NewVehicleHandler(
		c.CreateVehicle, c.GetVehicle, c.ListVehicles,
		c.UpdateVehicle, c.DeleteVehicle,
	)
	c.ServiceOrderHandler = serviceorderhandler.NewServiceOrderHandler(
		c.CreateServiceOrder, c.GetServiceOrder, c.GetServiceOrderByCode,
		c.ListServiceOrders, c.ListServiceOrdersByDocument,
		c.UpdateServiceOrderStatus, c.UpdateServiceOrder, c.DeleteServiceOrder,
		c.UpdateServiceOrderStatusByCode,
		c.GetAverageExecutionTime,
	)
	c.ServiceHandler = servicehandler.NewServiceHandler(
		c.CreateService, c.GetService, c.ListServices,
		c.UpdateService, c.DeleteService,
	)
	c.PartHandler = parthandler.NewPartHandler(
		c.CreatePart, c.GetPart, c.ListParts,
		c.UpdatePart, c.DeletePart, c.AdjustPartStock,
	)
}
