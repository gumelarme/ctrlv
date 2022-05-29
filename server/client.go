package server

import (
	"fmt"
	"net/http"

	"github.com/gumelarme/ctrlv/db"
	"github.com/labstack/echo/v4"
)

// Index is the main page
func (s *server) Index(c echo.Context) error {
	posts, err := s.database.GetPosts(c.Request().Context())
	if err != nil {
		return Render500(c, err, "Error while retrieving posts")
	}

	var post db.Post
	if len(posts) > 0 {
		post = *posts[0]
	}

	var realPosts []db.Post
	for _, p := range posts {
		realPosts = append(realPosts, *p)

	}

	err = c.Render(200, "index.html", echo.Map{
		"Items": posts,
		"Post":  post,
	})

	fmt.Println(err)
	return err
}

// GetPost get a post by id
func (s *server) GetPost(c echo.Context) error {
	id := c.Param("id")
	post, err := s.database.GetPostById(c.Request().Context(), id)
	if err != nil {
		return c.Render(404, "404.html", nil)
	}

	posts, _ := s.database.GetPosts(c.Request().Context())
	return c.Render(200, "index.html", echo.Map{
		"Items": posts,
		"Post":  post,
	})
}

// SavePost upsert a post, it is the "action" of the index page's form
func (s *server) SavePost(c echo.Context) error {
	var post db.Post
	err := c.Bind(&post)
	if err != nil {
		// TODO: Bad requets error
		return err
	}

	saveAction := s.database.CreatePost
	if len(post.Id) > 0 {
		saveAction = s.database.UpdatePost
	}

	err = saveAction(c.Request().Context(), &post)
	if err != nil {
		return Render500(c, err, "Error while saving a post")
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/p/%s", post.Id))
}

func (s *server) DeletePost(c echo.Context) error {
	var post db.Post
	err := c.Bind(&post)
	if err != nil {
		return err
	}

	err = s.database.DeletePost(c.Request().Context(), post.Id)
	if err != nil {
		return Render500(c, err, "Error while deleting post")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
