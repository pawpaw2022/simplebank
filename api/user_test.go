package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	mockdb "github.com/pawpaw2022/simplebank/db/mock"
	db "github.com/pawpaw2022/simplebank/db/postgresql"
	"github.com/pawpaw2022/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.ComparePassword(arg.HashedPassword, e.password)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checker    func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				// Build stubs
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(args, password)).
					Times(1).
					Return(user, nil)
			},
			checker: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "Invalid email",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "invalid_email",
			},
			buildStubs: func(store *mockdb.MockStore) {

				// Build stubs
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checker: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Password too short",
			body: gin.H{
				"username":  user.Username,
				"password":  "123",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Build stubs
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checker: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Username",
			body: gin.H{
				"username":  "invalid-username",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {

				// Build stubs
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checker: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				// Build stubs
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(args, password)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checker: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Duplicate Username",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				// Build stubs
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(args, password)).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)
			},
			checker: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the response
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// Build stubs
			tc.buildStubs(store)

			// Create a test server
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// Create a request
			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			require.NoError(t, err)

			// Send the request
			server.router.ServeHTTP(recorder, request)

			// Check the response
			tc.checker(t, recorder)
		})
	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
