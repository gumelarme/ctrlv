package server

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
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
}

func InitServer(e *echo.Echo) {
	s := server{}

	e.GET("/", s.Index)

	api := e.Group("/api")
	{
		api.GET("", s.ApiIndex)

		api.GET("/post", s.ApiGetPosts)
		api.POST("/post", s.SavePost)
		api.PUT("/post/:id", s.SavePost)
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
