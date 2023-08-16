package token

import (
	"testing"
	"time"

	"github.com/pawpaw2022/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasteoMaker(util.RandomString(32))
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

func TestExpiredPasetoMaker(t *testing.T) {
	maker, err := NewPasteoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
