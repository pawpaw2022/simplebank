package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/google/uuid"
)

var minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker.
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker.
func NewJWTMaker(secretKey string) (TokenMaker, error) {

	// Check if the secret key size is at least 32 bytes.
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

// implements TokenMaker interface

// CreateToken creates a new token for a specific username and duration.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        payload.ID,
		"username":  payload.Username,
		"issue_at":  payload.IssueAt.Unix(),
		"expire_at": payload.ExpireAt.Unix(),
	})

	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Check if the signing method is HMAC
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract the payload
	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Extract the payload data
	id, ok := payload["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid id claim")
	}

	username, ok := payload["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username claim")
	}

	issueAt, ok := payload["issue_at"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid issue_at claim")
	}

	expireAt, ok := payload["expire_at"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expire_at claim")
	}

	// check if token is expired
	if time.Unix(int64(expireAt), 0).Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return &Payload{
		ID:       uuid.MustParse(id),
		Username: username,
		IssueAt:  time.Unix(int64(issueAt), 0),
		ExpireAt: time.Unix(int64(expireAt), 0),
	}, nil

}
