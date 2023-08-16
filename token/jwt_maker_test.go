package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pawpaw2022/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issueAt := time.Now()
	expireAt := issueAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// Verify the payload data.
	require.Equal(t, username, payload.Username)
	require.NotZero(t, payload.ID)
	require.WithinDuration(t, issueAt, payload.IssueAt, time.Second)
	require.WithinDuration(t, expireAt, payload.ExpireAt, time.Second)
}

func TestExpiredJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestInvalidJWTMaker(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"id":        payload.ID,
		"username":  payload.Username,
		"issue_at":  payload.IssueAt.Unix(),
		"expire_at": payload.ExpireAt.Unix(),
	})
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, "invalid token signing method: none")
	require.Nil(t, payload)
}
