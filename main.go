package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/gumelarme/ctrlv/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	e          *echo.Echo
	echoLambda echoadapter.EchoLambdaV2
)

func init() {
	e = echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} [${status}] ${uri} (${latency_human}) from ${remote_ip} ${user_agent}",
	}))

	e.Renderer = server.NewRenderer()
	server.InitServer(e)

	fs := http.FileServer(http.Dir("./public/assets/"))

	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", fs)))
	echoLambda = *echoadapter.NewV2(e)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	if name := os.Getenv("AWS_LAMBDA_FUNCTION_NAME"); len(name) > 0 {
		log.Print("AWS echo cold start")
		lambda.Start(Handler)
	} else {
		fmt.Println("Starting local server...")
		if err := e.Start(":1234"); err != nil {
			panic(err)
		}
	}
}
