package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"

	"account_service/repository"

	errInternal "account_service/errors"
)

type TransactionService interface {
	Deposit(accountNumber string, amount float64) (*float64, error)
	Withdraw(accountNumber string, amount float64) (*float64, error)
	GetBalance(accountNumber string) (*float64, error)
}

type transactionService struct {
	accountRepository repository.AccountRepository
	userRepository    repository.UserRepository
}

func NewTransactionService(accountRepository repository.AccountRepository, userRepository repository.UserRepository) TransactionService {
	return &transactionService{
		accountRepository: accountRepository,
		userRepository:    userRepository,
	}
}

func (s *transactionService) Deposit(accountNumber string, amount float64) (*float64, error) {
	var newBalance float64

	account, err := s.accountRepository.FindByAccountNumber(accountNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &newBalance, errInternal.ErrAccountNotFound
		} else {
			log.Errorf("Error when finding account: %s", err)
			return &newBalance, err
		}
	}

	newBalance = account.Balance + amount
	err = s.accountRepository.UpdateBalance(account.ID, newBalance)
	if err != nil {
		log.Errorf("Error when updating balance: %s", err)
		return &newBalance, err
	}
	return &newBalance, nil
}

func (s *transactionService) Withdraw(accountNumber string, amount float64) (*float64, error) {
	var newBalance float64

	account, err := s.accountRepository.FindByAccountNumber(accountNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &newBalance, errInternal.ErrAccountNotFound
		} else {
			log.Errorf("Error when finding account: %s", err)
			return &newBalance, err
		}
	}

	newBalance = account.Balance - amount
	if newBalance < 0 {
		return &newBalance, errInternal.ErrInsufficientBalance
	}

	err = s.accountRepository.UpdateBalance(account.ID, newBalance)
	if err != nil {
		log.Errorf("Error when updating balance: %s", err)
		return &newBalance, err
	}
	return &newBalance, nil
}

func (s *transactionService) GetBalance(accountNumber string) (*float64, error) {
	account, err := s.accountRepository.FindByAccountNumber(accountNumber)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &account.Balance, errInternal.ErrAccountNotFound
	}
	return &account.Balance, err
}
