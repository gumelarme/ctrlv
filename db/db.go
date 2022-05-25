package db

import (
	"context"
)

type Database interface {
	CreatePost(context.Context, *Post) error
	GetPosts(context.Context) ([]*Post, error)
	GetPostById(ctx context.Context, id string) (*Post, error)
	// SearchPost(keyword string) []Post
	// UpdatePost(context.Context, id string, *Post) error
	// DeletePost(context.Context, id string) error
}
