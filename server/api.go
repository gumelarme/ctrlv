package server

import (
	"net/http"

	"github.com/gumendol/ctrlv/db"
	"github.com/labstack/echo/v4"
)

func (s *server) ApiGetPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"data": echo.Map{
			"post": db.GetPosts(nil),
		},
	})
}

func (s *server) SavePost(c echo.Context) error {
	var post db.Post
	if err := c.Bind(&post); err != nil {
		return err
	}

	if c.Request().Method == "PUT" {
		post.Id = c.Param("id")
	}

	id, err := post.Save()
	if err != nil {
		//TODO wrap error
		return err
	}

	return c.JSON(http.StatusCreated, data(echo.Map{
		"Id": id,
	}))
}
