package server

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gumendol/ctrlv/db"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var cachedPosts []db.Post

// invalidateCache clear all cache
func invalidateCache() {
	cachedPosts = []db.Post{}
}

// getAllPost retrieve all posts and read cached if any
func getAllPost() []db.Post {
	if len(cachedPosts) > 0 {
		return cachedPosts
	}

	log.Println("no cache available, retreiving fresh data")

	posts := db.GetPosts(nil)
	// ULID feature, sorted lexicographically ==  chronographically
	sort.Slice(posts, func(i, j int) bool {
		return strings.Compare(posts[i].Id, posts[j].Id) < 0
	})

	// TODO: limit cache size, FIFO
	cachedPosts = posts
	return posts
}

// doSavePost upsert post
func (s *server) doSavePost(c echo.Context) (*db.Post, error) {
	var post db.Post
	if err := c.Bind(&post); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	_, err := post.Save()
	if err != nil {
		return nil, errors.Wrap(err, "failed to save post")
	}
	return &post, nil
}

// doGetPostById retreive post by id by first checking inside cache
func (s *server) doGetPostById(id string) (*db.Post, error) {
	posts := getAllPost()
	index := -1
	for i, p := range posts {
		if p.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		// TODO: directly query the db
		return nil, fmt.Errorf("post with id `%s` is not found", id)
	}

	return &posts[index], nil
}
