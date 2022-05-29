package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/gumelarme/ctrlv/config"
	"github.com/gumelarme/ctrlv/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	db.PostData `bson:",inline"`
}

func (p *Post) ToDBPost() *db.Post {
	data := p.PostData
	data.CreatedAt = p.Id.Timestamp()

	return &db.Post{
		Id:       p.Id.Hex(),
		PostData: data,
	}
}

type MongoAPI struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func NewMongoAPI() *MongoAPI {
	return &MongoAPI{
		Host:     config.MongoDB.Host,
		Port:     config.MongoDB.Port,
		Username: config.MongoDB.Username,
		Password: config.MongoDB.Password,
		Database: config.MongoDB.Database,
	}
}

func (m *MongoAPI) withMongo(ctx context.Context, f func(*mongo.Database) error) error {
	timeout := time.Second * 3
	client, err := mongo.Connect(
		ctx,
		&options.ClientOptions{
			Hosts:          []string{fmt.Sprintf("%s:%d", m.Host, m.Port)},
			ConnectTimeout: &timeout,
			Auth: &options.Credential{
				Username:      m.Username,
				Password:      m.Password,
				AuthMechanism: "SCRAM-SHA-1",
				AuthSource:    "admin",
			},
		},
	)

	if err != nil {
		return err
	}

	var closingError error
	defer func() {
		closingError = client.Disconnect(ctx)
	}()

	if err := f(client.Database(m.Database)); err != nil {
		return err
	}

	return closingError
}
