package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/sangketkit01/media-library-api/internal/config"
	db "github.com/sangketkit01/media-library-api/internal/db/sqlc"
	"github.com/sangketkit01/media-library-api/internal/token"
)

const (
	webPort = "8099"
)

type App struct {
	store      db.Store
	router     *fiber.App
	config     *config.Config
	tokenMaker token.Maker
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

	fmt.Println(env)

	config, err := config.NewConfig("../../", env)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("environment:", config.Environment)
	fmt.Println(config)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	database, err := pgx.Connect(ctx, config.DatabaseUrl)
	if err != nil {
		log.Fatalln(err)
	}

	if err := database.Ping(ctx); err != nil {
		log.Fatalln(err)
	}

	log.Println("Sucess connected to database")

	store := db.NewStore(database)

	tokenMaker, err := token.NewPasetoMaker(config.Secretkey)
	if err != nil {
		log.Fatalln(tokenMaker)
	}

	app := App{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	app.router = app.routes()

	app.router.Listen(fmt.Sprintf(":%s", webPort))
}
