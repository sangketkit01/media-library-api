package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sangketkit01/media-library-api/internal/config"
)

const(
	webPort = "8099"
)

type App struct {
	router *fiber.App
	config *config.Config
}

func init(){
	if err := godotenv.Load("../../.env.local") ; err != nil{

		log.Println("failed to load .env.local. Try to load .env.production...")
		if err := godotenv.Load("../../.env.production") ;  err != nil{
			log.Println("failed to load .env.production")
		}else{
			log.Println("Success Loading .env.production")
		}

	}else{
		log.Println("Succeed Loading .env.local")
	}
}

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local"
	}

	config, err := config.NewConfig("../../", env)
	if err != nil{
		log.Fatalln(err)
	}

	app := App{
		config: config,
	}
	
	app.router = app.routes()

	app.router.Listen(fmt.Sprintf(":%s", webPort))
}