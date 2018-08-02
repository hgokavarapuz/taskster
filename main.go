package main

import (
	"log"
	"os"
	"github.com/kataras/iris"
	"github.com/joho/godotenv"
	A "github.com/zendesk/taskster/api"
	R "github.com/zendesk/taskster/routes"
	"github.com/kataras/iris/middleware/recover"
)

func newApp() *iris.Application {
	app := iris.New()
	//app.Logger().SetLevel("debug")
	app.Use(recover.New())
	//app.Use(logger.New())
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))
	app.StaticWeb("/", "./assets")

	// Routes
	app.Get("/", R.IndexHandler)

	// group all api routes here
	apiGroup := app.Party("/api")
	{
		apiGroup.Get("/hello", A.HelloHandler)
	}

	// group all non-api routes here
	nonApiGroup := app.Party("/hello")
	{
		nonApiGroup.Get("/:id", R.IdHandler)
	}

	return app
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := newApp()
	app.Run(iris.Addr(":" + getPort()), iris.WithoutServerError(iris.ErrServerClosed))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	if port[0] == ':' {
		port = port[1:]
	}
    return port
}
