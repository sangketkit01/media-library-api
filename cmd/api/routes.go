package main

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func (app *App) routes() *fiber.App{
	router := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	router.Get("/", func(c *fiber.Ctx) error{
		return c.JSON(fiber.Map{"message" : "Hello world"})
	})

	return router
}