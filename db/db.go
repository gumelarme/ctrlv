package db

import (
	"context"
)

type Database interface {
	CreatePost(context.Context, *Post) error
	GetPosts(context.Context) ([]*Post, error)
	GetPostById(ctx context.Context, id string) (*Post, error)
	UpdatePost(ctx context.Context, id string, post *Post) error
	// SearchPost(keyword string) []Post
	// DeletePost(context.Context, id string) error
}
