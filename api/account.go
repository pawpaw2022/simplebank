package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/pawpaw2022/simplebank/db/postgresql"
)

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required,gt=2,lte=100"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountParams

	// Assign the request body to the req variable
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

type GetAccountParams struct {
	ID int64 `uri:"id" binding:"required,min=1"` // uri: path parameter
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountParams

	// Assign the request body to the req variable
	if err := ctx.ShouldBindUri(&req); err != nil {
		// Invalid User Input
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the account
	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Account not found
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// Database Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Insert success, return the account
	ctx.JSON(http.StatusOK, account)
}

type ListAccountParams struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"` // form: query parameter
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountParams

	// Assign the request body to the req variable
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// Invalid User Input
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the account
	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		// Database Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Insert success, return the account
	ctx.JSON(http.StatusOK, accounts)
}

type UpdateAccountBalanceUri struct {
	ID int64 `uri:"id" binding:"required,min=1"` // uri: path parameter
}

type UpdateAccountBalanceJSON struct {
	Amount int64 `json:"amount" binding:"required,ne=0"`
}

func (server *Server) updateAccountBalance(ctx *gin.Context) {
	var uri UpdateAccountBalanceUri
	var json UpdateAccountBalanceJSON

	// Assign the request body to the req variable
	if err := ctx.ShouldBindUri(&uri); err != nil {
		// Invalid User Input
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Assign the request body to the req variable
	if err := ctx.ShouldBindJSON(&json); err != nil {
		// Invalid User Input
		fmt.Println("err", err)

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get the account
	account, err := server.store.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		ID:     uri.ID,
		Amount: json.Amount,
	})

	if err != nil {
		// Database Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Insert success, return the account
	ctx.JSON(http.StatusOK, account)
}
