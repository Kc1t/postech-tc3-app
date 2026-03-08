package mongodb

import (
	"context"

	mongomodel "github.com/fiap/postech-tc1/internal/adapters/outbound/mongodb/model"
	"github.com/fiap/postech-tc1/internal/domain/service"
	"github.com/fiap/postech-tc1/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const serviceCollection = "services"

type serviceRepository struct {
	collection *mongodriver.Collection
}

func NewServiceRepository(db *mongodriver.Database) ports.ServiceRepository {
	return &serviceRepository{collection: db.Collection(serviceCollection)}
}

func (r *serviceRepository) Create(ctx context.Context, s *service.Service) error {
	doc := mongomodel.FromService(s)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	s.SetID(doc.ID.Hex())
	return nil
}

func (r *serviceRepository) FindByID(ctx context.Context, id string) (*service.Service, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var doc mongomodel.Service
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, err
	}
	return doc.ToDomain(), nil
}

func (r *serviceRepository) FindAll(ctx context.Context) ([]*service.Service, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []mongomodel.Service
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	services := make([]*service.Service, 0, len(docs))
	for i := range docs {
		services = append(services, docs[i].ToDomain())
	}
	return services, nil
}

func (r *serviceRepository) Update(ctx context.Context, s *service.Service) error {
	doc := mongomodel.FromService(s)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": doc.ID}, doc)
	return err
}

func (r *serviceRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
