package server

import (
	"github.com/gumendol/ctrlv/db"
	"github.com/labstack/echo/v4"
)

type postWithTimestamp struct {
	db.Post
	Timestamp string
}

func (s *server) Index(c echo.Context) error {
	var posts []postWithTimestamp
	for _, p := range db.GetPosts(nil) {
		t := db.GetTimeFromId(p.Id)
		posts = append(posts, postWithTimestamp{
			Post:      p,
			Timestamp: t.Format("Mon, 02 Jan 2006 15:04:05"),
		})
	}

	return c.Render(200, "index.html", map[string]interface{}{
		"Name":  "world!",
		"Items": posts,
		"Post":  posts[0],
	})
}
