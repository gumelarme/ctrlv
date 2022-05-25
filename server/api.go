package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gumelarme/ctrlv/db"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// // ApiGetPost get a single post
func (s *server) ApiGetPost(c echo.Context) error {
	id := c.Param("id")
	post, err := s.database.GetPostById(c.Request().Context(), id)

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
	err := s.database.UpdatePost(c.Request().Context(), c.Param("id"), post)
	fmt.Println(err)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, data(post))
}

// // ApiDeletePost delete a post
// func (s *server) ApiDeletePost(c echo.Context) error {
// 	if err := db.Delete(c.Param("id")); err != nil {
// 		return err
// 	}

// 	invalidateCache()
// 	return c.NoContent(http.StatusNoContent)
// }
