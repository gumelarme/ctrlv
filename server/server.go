package server

import (
	"html/template"
	"io"

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

type server struct {
	database db.Database
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
	// e.POST("/p", s.SavePost)
	// e.POST("/p/delete", s.DeletePost)

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

func (s server) ApiIndex(c echo.Context) error {
	return c.JSON(200, data("Hello"))
}

func data(d interface{}) echo.Map {
	return map[string]interface{}{
		"data": d,
	}
}
