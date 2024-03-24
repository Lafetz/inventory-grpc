package inventory

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	productCollection *mongo.Collection
}

func (s *store) AddProduct(ctx context.Context, title string, desc string) (*product, error) {
	product := NewProduct(title, desc)
	res, err := s.productCollection.InsertOne(ctx, product)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	fmt.Printf(`%s`, res)
	return product, nil
}
func (s *store) GetProduct(ctx context.Context, id uuid.UUID) (*product, error) {
	var result product
	err := s.productCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrProductNotFound
		}
		fmt.Println(err)
		return nil, err
	}

	return &product{Id: id, Title: result.Title, Description: result.Description}, nil
}

func NewDb() *store {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://admin:admin11@localhost/inventory?authSource=admin"))
	if err != nil {
		log.Fatal("db connection falled", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("db connection falled", err)
	}
	productCollection := client.Database("inventory").Collection("product")
	return &store{
		productCollection: productCollection,
	}
}
