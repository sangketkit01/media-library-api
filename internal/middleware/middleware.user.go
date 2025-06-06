package middleware

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sangketkit01/media-library-api/internal/token"
)

const(
	authorizationHeader = "authorization"
	bearer = "bearer"
	payloadHeader = "payload"
)

func (m *Middleware)AuthMiddleware() fiber.Handler{
	return func(c *fiber.Ctx) error{
		authorization := c.Get(authorizationHeader)
		if strings.TrimSpace(authorization) == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing authorization header.")
		}

		parts := strings.Split(authorization, " ")
		if len(parts) < 2 || strings.ToLower(parts[0]) != bearer{
			return fiber.NewError(fiber.StatusBadRequest, "invalid authorization header format.")
		}

		accessToken := parts[1]
		payload, err := m.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) {
				return fiber.NewError(fiber.StatusUnauthorized, "token expired")
			}
		
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		

		c.Locals(payloadHeader, payload)

		return c.Next()
	}
}

func (m *Middleware) LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		log.Printf("%s [%d] %s %s - %v\n", c.IP(), status, c.Method(), c.Path(), duration)

		return err
	}
}
