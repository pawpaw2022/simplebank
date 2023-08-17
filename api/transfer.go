package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/pawpaw2022/simplebank/db/postgresql"
	"github.com/pawpaw2022/simplebank/token"
)

// TransferTxParams contains the input parameters of the transfer transaction
type TransferInputParams struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// Authorization: a logged-in user can only send money from his own account.
func (s *Server) createTransfer(ctx *gin.Context) {
	var req TransferInputParams

	// Assign the request body to the req variable
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Invalid User Input
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := s.validateUser(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	// Get the owner from the token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := fmt.Errorf("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = s.validateUser(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// Execute the transfer transaction
	result, err := s.store.TransferTx(ctx, arg)

	if err != nil {
		// Database Error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Insert success, return the account
	ctx.JSON(http.StatusOK, result)
}

// validateUser validates the currency of the account
func (s *Server) validateUser(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	// Get the account
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return account, false
	}

	// Validate the currency
	if account.Currency != currency {

		err := fmt.Errorf("accountID [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
