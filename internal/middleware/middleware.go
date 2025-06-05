package middleware

import "github.com/sangketkit01/media-library-api/internal/token"

type Middleware struct {
	tokenMaker token.Maker
}

func NewMiddleware(tokenMaker token.Maker) *Middleware{
	return &Middleware{
		tokenMaker: tokenMaker,
	}
}