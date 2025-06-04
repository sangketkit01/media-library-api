package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	SessionID uuid.UUID `json:"session_id"`
	ID uuid.UUID `json:"id"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(id uuid.UUID, duration time.Duration) (*Payload, error){
	sessionId, err := uuid.NewUUID()
	if err != nil{
		return nil, err
	}

	payload := &Payload{
		SessionID: sessionId,
		ID: id,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error{
	if time.Now().After(p.ExpiredAt) {
		return errors.New("token is expired.")
	}

	if p.SessionID == uuid.Nil{
		return errors.New("invalid session ID")
	}

	if p.ID == uuid.Nil{
		return errors.New("invalid ID")
	}

	return nil
}