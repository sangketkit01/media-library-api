package handlers

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sangketkit01/media-library-api/internal/db/sqlc"
	"github.com/sangketkit01/media-library-api/internal/token"
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
		ID:        user.ID.Bytes,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}

	return c.JSON(response)
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserResponse struct {
	Token                 string    `json:"token"`
	RefreshToken          string    `json:"refresh_token"`
	SessionID             uuid.UUID `json:"session_id"`
	TokenIssuedAt         time.Time `json:"token_issued_at"`
	TokenExpiredAt        time.Time `json:"token_expired_at"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired"`
	ID                    uuid.UUID `json:"id"`
	Email                 string    `json:"email"`
	CreatedAt             time.Time `json:"created_at"`
}

func (h *Handler) LoginUser(c *fiber.Ctx) error {
	var req LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	}

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	}

	user, err := h.Store.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "invalid email credential.")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := util.CheckPassword(user.Password, req.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid password credential")
	}

	sessionId, _ := uuid.NewUUID()

	accessToken, accessPayload, err := h.tokenMaker.CreateToken(user.ID.Bytes, sessionId,h.Config.AccessTokenDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, refreshPayloay, err := h.tokenMaker.CreateToken(user.ID.Bytes, sessionId,h.Config.RefreshTokenDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	arg := db.CreateSessionParams{
		RefreshToken: refreshToken,

		ID: pgtype.UUID{
			Bytes: sessionId,
			Valid: true,
		},

		UserID: pgtype.UUID{
			Bytes: accessPayload.ID,
			Valid: true,
		},

		UserAgent: pgtype.Text{
			String: c.Get("User-Agent"),
			Valid: true,
		},

		ClientIp: pgtype.Text{
			String: c.IP(),
			Valid: true,
		},

		IsBlocked: pgtype.Bool{
			Bool: false,
			Valid: true,
		},

		ExpiresAt: pgtype.Timestamptz{
			Time: refreshPayloay.ExpiredAt,
			Valid: true,
		},
	}

	_, err = h.Store.CreateSession(c.Context(), arg)

	response := LoginUserResponse{
		Token:                 accessToken,
		RefreshToken:          refreshToken,
		SessionID:             accessPayload.SessionID,
		TokenIssuedAt:         accessPayload.IssuedAt,
		TokenExpiredAt:        accessPayload.ExpiredAt,
		RefreshTokenExpiredAt: refreshPayloay.ExpiredAt,
		ID:                    user.ID.Bytes,
		Email:                 user.Email,
		CreatedAt:             user.CreatedAt.Time,
	}

	return c.JSON(response)
}

func (h *Handler) GetCurrentUser(c *fiber.Ctx) error {
	p := c.Locals("payload")
	payload, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid payload")
	}

	arg := pgtype.UUID{
		Bytes: payload.ID,
		Valid: true,
	}

	user, err := h.Store.GetUserByID(c.Context(), arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(fiber.StatusForbidden, "user not found")
		}

		log.Println(util.RouteCustomError(err, c.Path()))
		return fiber.NewError(fiber.StatusInternalServerError, "cannot get server")
	}

	return c.JSON(user)
}

func (h *Handler) LogoutUser(c *fiber.Ctx) error {
    p := c.Locals("payload")
    payload, ok := p.(*token.Payload)
    if !ok {
        return fiber.NewError(fiber.StatusUnauthorized, "invalid payload")
    }

    err := h.Store.BlockSessionByID(c.Context(), pgtype.UUID{
		Bytes: payload.SessionID,
		Valid: true,
	})
    if err != nil {
		log.Println(err, payload.SessionID)
        return fiber.NewError(fiber.StatusInternalServerError, "cannot block session")
    }

    return c.JSON(fiber.Map{
        "message": "logged out successfully",
    })
}

