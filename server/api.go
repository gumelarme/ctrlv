package server

import (
	"net/http"

	"github.com/gumendol/ctrlv/db"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (s *server) ApiGetPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"data": echo.Map{
			"post": db.GetPosts(nil),
		},
	})
}

func (s *server) ApiSavePost(c echo.Context) error {
	post, err := s.doSavePost(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, data(post))
}

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
