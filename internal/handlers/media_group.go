package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sangketkit01/media-library-api/internal/db/sqlc"
	"github.com/sangketkit01/media-library-api/internal/token"
	"github.com/sangketkit01/media-library-api/internal/util"
)

type CreateGroupRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateGroupResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"create_at"`
}

func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	var req CreateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	}

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	}

	p := c.Locals("payload")
	payload, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
	}

	user, err := h.Store.GetUserByID(c.Context(), pgtype.UUID{
		Bytes: payload.ID,
		Valid: true,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		util.RouteCustomError(err, c.Path())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retreive user data.")
	}

	arg := db.CreateMediaGroupParams{
		UserID: pgtype.UUID{
			Bytes: user.ID.Bytes,
			Valid: true,
		},

		Name: req.Name,
	}

	group, err := h.Store.CreateMediaGroup(c.Context(), arg)

	if err != nil {

		util.RouteCustomError(err, c.Path())

		return fiber.NewError(fiber.StatusInternalServerError, "failed to create group.")
	}

	response := CreateGroupResponse{
		ID:        group.ID.Bytes,
		UserID:    group.UserID.Bytes,
		Name:      group.Name,
		CreatedAt: group.CreatedAt.Time,
	}

	return c.JSON(response)
}
