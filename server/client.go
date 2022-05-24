package server

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gumelarme/ctrlv/db"
// 	"github.com/labstack/echo/v4"
// )

// // Index is the main page
// func (s *server) Index(c echo.Context) error {
// 	posts := getAllPost()

// 	var post db.Post
// 	if len(posts) > 0 {
// 		post = posts[0]
// 	}

// 	return c.Render(200, "index.html", echo.Map{
// 		"Items": posts,
// 		"Post":  post,
// 	})
// }

// // GetPost get a post by id
// func (s *server) GetPost(c echo.Context) error {
// 	id := c.Param("id")
// 	post, err := s.doGetPostById(id)
// 	if err != nil {
// 		return c.Render(404, "404.html", nil)
// 	}

// 	return c.Render(200, "index.html", echo.Map{
// 		"Items": getAllPost(),
// 		"Post":  post,
// 	})
// }

// // SavePost upsert a post, it is the "action" of the index page's form
// func (s *server) SavePost(c echo.Context) error {
// 	post, err := s.doSavePost(c)
// 	if err != nil {
// 		return err
// 	}

// 	invalidateCache()
// 	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/p/%s", post.Id))
// }

// func (s *server) DeletePost(c echo.Context) error {
// 	var post db.Post
// 	c.Bind(&post)

// 	if err := db.Delete(post.Id); err != nil {
// 		return err
// 	}

// 	invalidateCache()
// 	return c.Redirect(http.StatusSeeOther, "/")
// }
