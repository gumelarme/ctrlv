package mongo

import (
	"context"
	"fmt"

	"github.com/gumelarme/ctrlv/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoAPI) GetPostById(ctx context.Context, id string) (*db.Post, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("the provided id is invalid")
	}

	var post Post
	err = m.withMongo(ctx, func(d *mongo.Database) error {
		coll := d.Collection("posts")
		res := coll.FindOne(ctx, bson.M{"_id": _id})

		if err := res.Err(); err != nil {
			return fmt.Errorf("post `%s` doesn't exist", id)
		}

		res.Decode(&post)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return post.ToDBPost(), nil
}

func (m *MongoAPI) GetPosts(ctx context.Context) ([]*db.Post, error) {
	var posts []Post

	err := m.withMongo(ctx, func(db *mongo.Database) error {
		coll := db.Collection("posts")
		cursor, err := coll.Find(ctx, bson.D{})
		if err != nil {
			return err
		}

		return cursor.All(ctx, &posts)
	})

	if err != nil {
		return nil, err
	}

	var dbPosts []*db.Post
	for _, p := range posts {
		dbPosts = append(dbPosts, p.ToDBPost())
	}

	return dbPosts, nil
}

func (m *MongoAPI) CreatePost(ctx context.Context, post *db.Post) error {
	p, _ := NewFromPost(*post, false)
	err := m.withMongo(ctx, func(db *mongo.Database) error {
		posts := db.Collection("posts")
		res, err := posts.InsertOne(ctx, p)
		post.Id = res.InsertedID.(primitive.ObjectID).Hex()
		return err
	})

	if err != nil {
		return errors.Wrap(err, "error while creating new post")
	}

	return nil
}
