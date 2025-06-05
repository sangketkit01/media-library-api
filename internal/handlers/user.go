package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	db "github.com/sangketkit01/media-library-api/internal/db/sqlc"
	"github.com/sangketkit01/media-library-api/internal/util"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Passowrd string `json:"password" validate:"required,min=8"`
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	}

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	}

	hashedPassword, err := util.HashPassword(req.Passowrd)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	arg := db.CreateUserParams{
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := h.Store.CreateUser(c.Context(), arg)
	if err != nil {
		if pgError, ok := err.(*pgx.PgError); ok && pgError.Code == util.UniqueViolationErrCode {
			return fiber.NewError(fiber.StatusConflict, "Email is already exists.")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := CreateUserResponse{
		ID: user.ID.Bytes,
		Email: user.Email,
		CreatedAt: user.CreatedAt.Time,
	}

	return c.JSON(response)
}
