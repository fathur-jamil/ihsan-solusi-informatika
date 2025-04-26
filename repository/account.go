package repository

import (
	"github.com/jinzhu/gorm"

	"account_service/model"
)

type AccountRepository interface {
	Create(account *model.Account) error
	FindByAccountNumber(accountNumber string) (*model.Account, error)
	UpdateBalance(accountID int64, newBalance float64) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db}
}

func (r *accountRepository) Create(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) FindByAccountNumber(accountNumber string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("account_number = ?", accountNumber).First(&account).Error
	return &account, err
}

func (r *accountRepository) UpdateBalance(accountID int64, newBalance float64) error {
	return r.db.Model(&model.Account{}).Where("id = ?", accountID).Update("balance", newBalance).Error
}
