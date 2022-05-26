package db

import (
	"context"
)

type Database interface {
	CreatePost(context.Context, *Post) error
	GetPosts(context.Context) ([]*Post, error)
	GetPostById(ctx context.Context, id string) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, id string) error
	// SearchPost(keyword string) []Post
}
