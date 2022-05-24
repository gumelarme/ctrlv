package mongo

import (
	"context"

	"github.com/gumelarme/ctrlv/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoAPI) GetPosts(ctx context.Context) ([]db.Post, error) {
	var posts []db.Post

	err := m.withMongo(ctx, func(db *mongo.Database) error {
		coll := db.Collection("posts")
		cursor, err := coll.Find(ctx, nil)
		if err != nil {
			return err
		}

		return cursor.All(ctx, posts)
	})

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *MongoAPI) CreatePost(ctx context.Context, post *db.Post) error {
	err := m.withMongo(ctx, func(db *mongo.Database) error {
		posts := db.Collection("posts")
		res, err := posts.InsertOne(ctx, post)
		post.Id = res.InsertedID.(primitive.ObjectID).Hex()
		return err
	})

	if err != nil {
		return errors.Wrap(err, "error while creating new post")
	}

	return nil
}
