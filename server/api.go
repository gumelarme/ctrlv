package server

import (
	"log"
	"net/http"

	"github.com/gumelarme/ctrlv/db"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Index for api
func (s server) ApiIndex(c echo.Context) error {
	return c.JSON(200, data("Hello"))
}

// ApiGetPost get a single post
func (s *server) ApiGetPost(c echo.Context) error {
	id := c.Param("id")
	post, err := s.database.GetPostById(c.Request().Context(), id)

	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": err.Error(),
			"nice":  "Hello",
		})
	}

	return c.JSON(http.StatusOK, data(echo.Map{
		"post": post,
	}))
}

// ApiGetPosts get all posts
// TODO: paginate
func (s *server) ApiGetPosts(c echo.Context) error {
	posts, err := s.database.GetPosts(c.Request().Context())
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "error while fetching posts")
	}

	return c.JSON(http.StatusOK, data(echo.Map{
		"posts": posts,
	}))
}

// ApiSavePost take a whole Post object and upsert it
func (s *server) ApiSavePost(c echo.Context) error {
	var post *db.Post
	if err := c.Bind(&post); err != nil {
		return errors.Wrap(err, "bad request")
	}

	err := s.database.CreatePost(c.Request().Context(), post)
	if err != nil {
		return errors.Wrap(err, "something wrong while creating post")
	}

	return c.JSON(http.StatusCreated, data(post))
}

// ApiUpdatePost update a post
func (s *server) ApiUpdatePost(c echo.Context) error {
	var post *db.Post
	c.Bind(&post)
	post.Id = c.Param("id")

	err := s.database.UpdatePost(c.Request().Context(), post)
	if err != nil {
		return err
	}

	post, err = s.database.GetPostById(c.Request().Context(), post.Id)
	if err != nil {
		return c.JSON(http.StatusMultiStatus, echo.Map{
			"message": "update successul but failed to retrieve the changed post",
		})
	}

	return c.JSON(http.StatusCreated, data(post))
}

// ApiDeletePost delete a post
func (s *server) ApiDeletePost(c echo.Context) error {
	err := s.database.DeletePost(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
