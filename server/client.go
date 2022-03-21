package server

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gumendol/ctrlv/db"
	"github.com/labstack/echo/v4"
)

type postWithTimestamp struct {
	db.Post
	Timestamp string
}

var cachedPosts []postWithTimestamp

func invalidateCache() {
	cachedPosts = []postWithTimestamp{}
}

func getAllPost() []postWithTimestamp {
	if len(cachedPosts) > 0 {
		return cachedPosts
	}

	log.Println("no cache available, retreiving fresh data")
	var posts []postWithTimestamp

	for _, p := range db.GetPosts(nil) {
		t := db.GetTimeFromId(p.Id)
		posts = append(posts, postWithTimestamp{
			Post:      p,
			Timestamp: t.Format("Mon, 02 Jan 2006 15:04:05"),
		})
	}

	// ULID feature, sorted lexicographically ==  chronographically
	sort.Slice(posts, func(i, j int) bool {
		return strings.Compare(posts[i].Id, posts[j].Id) < 0
	})
	cachedPosts = posts
	return posts
}

func (s *server) Index(c echo.Context) error {
	posts := getAllPost()
	return c.Render(200, "index.html", echo.Map{
		"Items": posts,
		"Post":  posts[0],
	})
}

func (s *server) GetPost(c echo.Context) error {
	id := c.Param("id")
	posts := getAllPost()
	index := -1
	for i, p := range posts {
		if p.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		return c.Render(404, "404.html", nil)
	}

	return c.Render(200, "index.html", echo.Map{
		"Items": posts,
		"Post":  posts[index],
	})
}

func (s *server) SavePost(c echo.Context) error {
	post, err := s.doSavePost(c)
	if err != nil {
		return err
	}

	invalidateCache()
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/p/%s", post.Id))
}
