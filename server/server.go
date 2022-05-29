package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/gumelarme/ctrlv/config"
	"github.com/gumelarme/ctrlv/db"
	"github.com/gumelarme/ctrlv/db/mongo"
	"github.com/labstack/echo/v4"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer() *Renderer {
	return &Renderer{
		templates: template.Must(template.ParseGlob("public/template/*.html")),
	}
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return r.templates.ExecuteTemplate(w, name, data)
}

func InitServer(e *echo.Echo) {
	s := server{
		database: &mongo.MongoAPI{
			Host:     config.Conf.MongoDB.Host,
			Port:     config.Conf.MongoDB.Port,
			Username: config.Conf.MongoDB.Username,
			Password: config.Conf.MongoDB.Password,
			Database: config.Conf.MongoDB.Database,
		},
	}

	e.GET("/", s.Index)
	e.GET("/p/:id", s.GetPost)
	e.POST("/p", s.SavePost)
	e.POST("/p/delete", s.DeletePost)
	echo.NotFoundHandler = s.NotFoundHandler

	api := e.Group("/api")
	{
		api.GET("", s.ApiIndex)

		api.GET("/p", s.ApiGetPosts)
		api.GET("/p/:id", s.ApiGetPost)
		api.POST("/p", s.ApiSavePost)
		api.PUT("/p/:id", s.ApiUpdatePost)
		api.DELETE("/p/:id", s.ApiDeletePost)
	}
}

func data(d interface{}) echo.Map {
	return map[string]interface{}{
		"data": d,
	}
}

type server struct {
	database db.Database
}

func (s server) NotFoundHandler(c echo.Context) error {
	err := fmt.Errorf("resource not found")
	if strings.HasPrefix(c.Request().URL.Path, "/api") {
		return c.JSON(404, echo.Map{
			"error": err.Error(),
		})
	}

	return RenderError(c, 404, err, "Not found", "This page doesn't exist anymore")
}

func Render500(c echo.Context, err error, message string) error {
	return RenderError(c, http.StatusInternalServerError, err, "Internal Server Error", message)
}

func RenderError(c echo.Context, code int, err error, title, message string) error {
	fmt.Println(err)
	return c.Render(code, "error.html", echo.Map{
		"Code":    code,
		"Title":   title,
		"Message": message,
	})
}
