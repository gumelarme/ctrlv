package server

import (
	"net/http"

	"github.com/gumendol/ctrlv/db"
	"github.com/labstack/echo/v4"
)

// ApiGetPost get a single post
func (s *server) ApiGetPost(c echo.Context) error {
	id := c.Param("id")
	post, err := s.doGetPostById(id)

	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, data(echo.Map{
		"post": post,
	}))
}

// ApiGetPosts get all posts
// TODO: paginate
func (s *server) ApiGetPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, data(echo.Map{
		"posts": getAllPost(),
	}))
}

// ApiSavePost take a whole Post object and upsert it
func (s *server) ApiSavePost(c echo.Context) error {
	post, err := s.doSavePost(c)
	if err != nil {
		return err
	}

	invalidateCache()
	return c.JSON(http.StatusCreated, data(post))
}

// ApiUpdatePost update a post
func (s *server) ApiUpdatePost(c echo.Context) error {
	body := make(map[string]string)
	c.Bind(&body)
	post, err := db.UpdatePostByMap(c.Param("id"), body)

	if err != nil {
		return err
	}

	invalidateCache()
	return c.JSON(http.StatusCreated, data(post))
}

// ApiDeletePost delete a post
func (s *server) ApiDeletePost(c echo.Context) error {
	if err := db.Delete(c.Param("id")); err != nil {
		return err
	}

	invalidateCache()
	return c.NoContent(http.StatusNoContent)
}
