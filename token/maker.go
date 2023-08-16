package token

import "time"

// TokenMaker is an interface that creates and verifies tokens.
type TokenMaker interface {
	// CreateToken creates a new token for a specific username and duration.
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid or not.
	VerifyToken(token string) (*Payload, error)
}
