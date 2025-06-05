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
	router.Post("/login-user", handler.LoginUser)
	router.Post("/refresh-token", handler.RefreshToken)
	

	authRouter := router.Use(middleware.AuthMiddleware())
	authRouter.Get("/user", handler.GetCurrentUser)
	authRouter.Get("/logout", handler.LogoutUser)

	authRouter.Post("/media/upload", handler.UploadFile)
	authRouter.Post("/groups", handler.CreateGroup)
	authRouter.Patch("/media/:id/group/:group_id", handler.AssignMediaToGroup)

	authRouter.Get("/media", handler.GetCurrentUserMedia)
	authRouter.Get("/media/:id/download", handler.DownloadMedia)

	return &Route{
		App: router,
		Middleware: middleware,
		Handler: handler,
	}
}
