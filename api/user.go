package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/pawpaw2022/simplebank/db/postgresql"
	"github.com/pawpaw2022/simplebank/util"
)

type CreateUserParams struct {
	Username string `json:"username" binding:"required,min=6,alphanum"` // alphanum: only allow alphanumeric characters
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserParams

	// Assign the request body to the req variable
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Invalid User Input
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		// Hashing failed
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create the user
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})

	if err != nil {

		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}

		// Database Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Fill the response
	res := UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	// Insert success, return the account
	ctx.JSON(http.StatusOK, res)
}
