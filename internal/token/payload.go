package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidSessionID = errors.New("invalid session ID")
	ErrInvalidUserID = errors.New("invalid ID")
)

type Payload struct {
	SessionID uuid.UUID `json:"session_id"`
	ID        uuid.UUID `json:"id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(id uuid.UUID, sessionID uuid.UUID, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		SessionID: sessionID,
		ID:        id,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	if p.SessionID == uuid.Nil {
		return ErrInvalidSessionID
	}

	if p.ID == uuid.Nil {
		return ErrInvalidUserID
	}

	return nil
}
