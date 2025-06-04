package main

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const(
	authorizationHeader = "authorization"
	bearer = "bearer"
	payloadHeader = "payload"
)

func (app *App) authMiddleware() fiber.Handler{
	return func(c *fiber.Ctx) error{
		authorization := c.Get(authorizationHeader)
		if strings.TrimSpace(authorization) == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing authorization header.")
		}

		parts := strings.Split(authorization, " ")
		if len(parts) < 2 || strings.ToLower(parts[0]) != bearer{
			return fiber.NewError(fiber.StatusBadRequest, "invalid authorization header format.")
		}

		token := parts[1]
		payload, err := app.tokenMaker.VerifyToken(token)
		if err != nil{
			return fiber.NewError(fiber.StatusBadRequest, "invalid token.")
		}

		c.Locals(payloadHeader, payload)

		return c.Next()
	}
}

func (app *App) loggerMiddleware() fiber.Handler{
	return func(c *fiber.Ctx) error{
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		log.Printf("[%s] %s - %v\n", c.Method(), c.Path(), duration)

		return err
	}
}