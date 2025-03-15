// db/mongo/mongo.go
package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI        string
	Database   string
	Collection string
}

func (c *MongoConfig) Validate() error {
	if c.URI == "" {
		return errors.New("empty URI")
	}
	return nil
}

func (c *MongoConfig) DriverName() string {
	return "mongodb"
}

type MongoDriver struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (m *MongoDriver) Connect(config interface{}) error {
	cfg, ok := config.(*MongoConfig)
	if !ok {
		return errors.New("invalid config type")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	m.client = client
	m.collection = client.Database(cfg.Database).Collection(cfg.Collection)
	return err
}

// 实现其他接口方法...
