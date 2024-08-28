package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/yinnohs/simple-bank/db/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

type GetAccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) createAccout(ctx *gin.Context) {

	var request CreateAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateAccountParams{
		Owner:    request.Owner,
		Currency: request.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request GetAccountByIDRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccountById(ctx, request.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountsRequest struct {
	Limit  int32 `form:"page_size"`
	Offset int32 `form:"page_id"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var request ListAccountsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.FindAllAccountsParams{
		Limit:  request.Limit,
		Offset: (request.Offset - 1) * request.Limit,
	}

	fmt.Println(args)

	accounts, err := server.store.FindAllAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(accounts) <= 0 {
		ctx.JSON(http.StatusNotFound, "No records avalaible")
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type AddAccountBalanceRequest struct {
	AccountId int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required,min=1"`
}

type AddBalanceResponse struct {
	Entry   db.Entry   `json:"entry"`
	Account db.Account `json:"account"`
}

func (server *Server) addAccountBalance(ctx *gin.Context) {
	var request AddAccountBalanceRequest
	if err := ctx.ShouldBindBodyWithJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	accountArgs := db.AddToAccountBalanceParams{
		ID:     request.AccountId,
		Amount: request.Amount,
	}
	entryArgs := db.CreateEntryParams{
		AccountID: request.AccountId,
		Amount:    request.Amount,
	}

	updatedAccount, err := server.store.AddToAccountBalance(ctx, accountArgs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	entry, err := server.store.CreateEntry(ctx, entryArgs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, AddBalanceResponse{
		Account: updatedAccount,
		Entry:   entry,
	})

}
