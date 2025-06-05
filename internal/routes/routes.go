package routes

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/sangketkit01/media-library-api/internal/handlers"
	"github.com/sangketkit01/media-library-api/internal/middleware"
)

type Route struct{
	App *fiber.App
	Middleware *middleware.Middleware 
	Handler *handlers.Handler
}

func NewRoute(middleware *middleware.Middleware, handler *handlers.Handler) *Route {
	router := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})


	router.Use(middleware.LoggerMiddleware())
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello world"})
	})

	router.Post("/create-user", handler.CreateUser)

	_ = router.Use(middleware.AuthMiddleware())

	return &Route{
		App: router,
		Middleware: middleware,
		Handler: handler,
	}
}
