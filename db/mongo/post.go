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

// GetPostById return post with given string id
func (m *MongoAPI) GetPostById(ctx context.Context, id string) (*db.Post, error) {
	_id, err := isIdValid(id)
	if err != nil {
		return nil, err
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

// GetPosts get multiple posts
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

// CreatePost insert the given post to the databse
func (m *MongoAPI) CreatePost(ctx context.Context, post *db.Post) error {
	p := Post{PostData: post.PostData}
	err := m.withMongo(ctx, func(db *mongo.Database) error {
		res, err := db.Collection("posts").InsertOne(ctx, p)
		post.Id = res.InsertedID.(primitive.ObjectID).Hex()
		return err
	})

	if err != nil {
		return errors.Wrap(err, "error while creating new post")
	}

	return nil
}

// UpdatePost update post by its id, return error if id is empty
func (m *MongoAPI) UpdatePost(ctx context.Context, post *db.Post) error {
	_id, err := isIdValid(post.Id)
	if err != nil {
		return err
	}

	p := Post{PostData: post.PostData}
	err = m.withMongo(ctx, func(db *mongo.Database) error {
		_, err := db.Collection("posts").UpdateByID(ctx, _id, bson.M{"$set": p})
		return err
	})

	if err != nil {
		return errors.Wrapf(err, "error while updating post")
	}

	return nil
}

// DeletePost delete post that match the given id
func (m *MongoAPI) DeletePost(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = m.withMongo(ctx, func(db *mongo.Database) error {
		_, err := db.Collection("posts").DeleteOne(ctx, bson.M{"_id": _id})
		return err
	})

	if err != nil {
		return errors.Wrap(err, "error while deleting post")
	}

	return nil
}

// isIdValid check if the given id is valid
func isIdValid(id string) (primitive.ObjectID, error) {
	if len(id) == 0 {
		return primitive.NilObjectID, fmt.Errorf("invalid empty id")
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("the provided id is invalid")
	}
	return _id, nil
}
