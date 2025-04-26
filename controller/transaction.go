package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"account_service/dto"
	"account_service/service"

	errInternal "account_service/errors"
)

type transactionController struct {
	transactionService service.TransactionService
}

func NewTransactionController(transactionService service.TransactionService) *transactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

func (c *transactionController) Deposit(ctx echo.Context) error {
	var transactionRequest dto.TransactionRequest
	if err := ctx.Bind(&transactionRequest); err != nil {
		log.Errorf("Error when binding request: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorRequestIsNotValid,
			},
		)
	}

	if err := ctx.Validate(transactionRequest); err != nil {
		log.Errorf("Error when validating transaction request: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorRequestIsNotValid,
			},
		)
	}

	if transactionRequest.Amount <= 0 {
		log.Warnf("Invalid amount %f in transaction request", transactionRequest.Amount)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorInvalidAmount,
			},
		)
	}

	newBalance, err := c.transactionService.Deposit(
		transactionRequest.AccountNumber,
		transactionRequest.Amount,
	)
	if err == nil {
		log.Infof("Successfully deposit balance of %s", transactionRequest.AccountNumber)
		return ctx.JSON(
			http.StatusOK,
			dto.TransactionResponse{
				Balance: newBalance,
			},
		)
	}

	if errors.Is(err, errInternal.ErrAccountNotFound) {
		log.Warnf("Account %s not found in transaction service", transactionRequest.AccountNumber)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorAccountNotFound,
			},
		)
	}

	return ctx.JSON(
		http.StatusInternalServerError,
		dto.TransactionResponse{
			Remark: errInternal.TextErrorServer,
		},
	)
}

func (c *transactionController) Withdraw(ctx echo.Context) error {
	var transactionRequest dto.TransactionRequest
	if err := ctx.Bind(&transactionRequest); err != nil {
		log.Errorf("Error when binding request: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorRequestIsNotValid,
			},
		)
	}

	err := ctx.Validate(transactionRequest)
	if err != nil {
		log.Errorf("Error when validating transaction request: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorRequestIsNotValid,
			},
		)
	}

	if transactionRequest.Amount <= 0 {
		log.Warnf("Invalid amount %f in transaction request", transactionRequest.Amount)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorInvalidAmount,
			},
		)
	}

	newBalance, err := c.transactionService.Withdraw(
		transactionRequest.AccountNumber,
		transactionRequest.Amount,
	)
	if err == nil {
		log.Infof("Successfully withdraw balance of %s", transactionRequest.AccountNumber)
		return ctx.JSON(
			http.StatusOK,
			dto.TransactionResponse{
				Balance: newBalance,
			},
		)
	}

	if errors.Is(err, errInternal.ErrAccountNotFound) {
		log.Warnf("Account %s not found in transaction service", transactionRequest.AccountNumber)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorAccountNotFound,
			},
		)
	}

	if errors.Is(err, errInternal.ErrInsufficientBalance) {
		log.Warnf("Insufficient balance of %s", transactionRequest.AccountNumber)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorInsufficientBalance,
			},
		)
	}

	return ctx.JSON(
		http.StatusInternalServerError,
		dto.TransactionResponse{
			Remark: errInternal.TextErrorServer,
		},
	)
}

func (c *transactionController) GetBalance(ctx echo.Context) error {
	accountNumber := ctx.Param("no_rekening")
	balance, err := c.transactionService.GetBalance(accountNumber)
	if err == nil {
		log.Infof("Successfully get balance of %s", accountNumber)
		return ctx.JSON(
			http.StatusOK,
			dto.TransactionResponse{
				Balance: balance,
			},
		)
	}

	if errors.Is(err, errInternal.ErrAccountNotFound) {
		log.Warnf("Account %s not found in transaction service", accountNumber)
		return ctx.JSON(
			http.StatusBadRequest,
			dto.TransactionResponse{
				Remark: errInternal.TextErrorAccountNotFound,
			},
		)
	}

	log.Errorf("Error when getting balance: %s", err)
	return ctx.JSON(
		http.StatusInternalServerError,
		dto.TransactionResponse{
			Remark: errInternal.TextErrorServer,
		},
	)
}
