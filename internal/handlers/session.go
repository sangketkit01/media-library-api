package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad request.")
	}

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "refresh token is not provided.")
	}

	payload, err := h.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token.")
	}

	arg := pgtype.UUID{
		Bytes: payload.SessionID,
		Valid: true,
	}

	session, err := h.Store.GetSession(c.Context(), arg)
	if err != nil || session.IsBlocked.Bool || session.UserID.Bytes != payload.ID ||
		session.RefreshToken != req.RefreshToken || time.Now().After(session.ExpiresAt.Time) {

		return fiber.NewError(fiber.StatusInternalServerError, "invalid session")
	}

	accessToken, _, err := h.tokenMaker.CreateToken(payload.ID, payload.SessionID, h.Config.AccessTokenDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "cannot create new access token")
	}

	return c.JSON(fiber.Map{"access_token": accessToken})
}
