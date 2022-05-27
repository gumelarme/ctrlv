package mongo

import (
	"context"
	"fmt"
	"log"

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

func (m *MongoAPI) withMongo(ctx context.Context, f func(*mongo.Database) error) error {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%d/admin?connectTimeoutMS=10000&authSource=admin&authMechanism=SCRAM-SHA-1",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
	)

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(connStr),
	)

	if err != nil {
		log.Fatal(err)
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
