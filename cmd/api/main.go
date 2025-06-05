package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sangketkit01/media-library-api/internal/config"
	"github.com/sangketkit01/media-library-api/internal/handlers"
	"github.com/sangketkit01/media-library-api/internal/middleware"
	"github.com/sangketkit01/media-library-api/internal/routes"
	"github.com/sangketkit01/media-library-api/internal/token"
)

const (
	webPort = "8099"
)

type App struct {
	handler    *handlers.Handler
	config     *config.Config
	tokenMaker token.Maker
	router     *routes.Route
}

func init() {
	if err := godotenv.Load("../../.env.local"); err != nil {

		log.Println("failed to load .env.local. Try to load .env.production...")

		if err := godotenv.Load("../../.env.production"); err != nil {
			log.Println("failed to load .env.production")
		} else {
			log.Println("Success Loading .env.production")
		}

	} else {
		log.Println("Succeed Loading .env.local")
	}
}

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local"
	}


	config, err := config.NewConfig("../../", env)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("environment:", config.Environment)

	tokenMaker, err := token.NewPasetoMaker(config.Secretkey)
	if err != nil {
		log.Panic(tokenMaker)
	}

	middleware := middleware.NewMiddleware(tokenMaker)
	
	handler, err := handlers.NewHandler(config)
	if err != nil {
		log.Panic(err)
	}

	router := routes.NewRoute(middleware, handler)


	app := App{
		config:     config,
		tokenMaker: tokenMaker,
		handler:    handler,
		router:     router,
	}

	err = app.router.App.Listen(fmt.Sprintf(":%s", webPort))
	if err != nil{
		log.Panic(err)
	}
}
