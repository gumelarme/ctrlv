package main

import (
	"fmt"
	"net/http"

	"github.com/gumendol/ctrlv/db"
	"github.com/gumendol/ctrlv/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} [${status}] ${uri} (${latency_human}) from ${remote_ip} ${user_agent}\nError: ${error}\n",
	}))

	e.Renderer = server.NewRenderer()
	server.InitServer(e)

	fs := http.FileServer(http.Dir("./public/template/static"))

	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", fs)))
	e.Logger.Fatal(e.Start(":1234"))
}

func _main() {
	fmt.Println(db.NewULID())
}
