package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/pawpaw2022/simplebank/db/postgresql"
)

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountParams

	// Check if the request body is valid
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Invalid User Input
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Create the account
	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	})

	if err != nil {
		// Database Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Insert success, return the account
	ctx.JSON(http.StatusOK, account)
}
