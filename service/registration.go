package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/labstack/gommon/log"

	"account_service/model"
	"account_service/repository"
)

type RegistrationService interface {
	IsNIKOrPhoneNumberRegistered(nik string, phoneNumber string) (bool, error)
	Register(name, nik, phoneNumber string) (string, error)
}

type registrationService struct {
	accountRepository repository.AccountRepository
	userRepository    repository.UserRepository
}

func NewRegistrationService(accountRepository repository.AccountRepository, userRepository repository.UserRepository) RegistrationService {
	return &registrationService{
		accountRepository: accountRepository,
		userRepository:    userRepository,
	}
}

func (s *registrationService) IsNIKOrPhoneNumberRegistered(nik, phoneNumber string) (bool, error) {
	count, err := s.userRepository.CountByNIKOrPhoneNumber(nik, phoneNumber)
	if err != nil {
		log.Errorf("Error when checking nik or phone number: %s", err)
		return false, err
	}
	return count > 0, nil
}

func (s *registrationService) Register(name, nik, phoneNumber string) (string, error) {
	var (
		accountNumber string
		timeNow       = time.Now()
	)

	user := &model.User{
		Name:        name,
		NIK:         nik,
		PhoneNumber: phoneNumber,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}
	err := s.userRepository.Create(user)
	if err != nil {
		log.Errorf("Error when creating user: %s", err)
		return accountNumber, err
	}

	accountNumber = generateAccountNumber()
	err = s.accountRepository.Create(&model.Account{
		UserID:        user.ID,
		AccountNumber: accountNumber,
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
	})
	if err != nil {
		log.Errorf("Error when creating account: %s", err)
		errDeleteByID := s.userRepository.DeleteByID(user.ID)
		if errDeleteByID != nil {
			log.Errorf("Error when deleting user: %s", errDeleteByID)
		}
		return accountNumber, err
	}

	return accountNumber, nil
}

func generateAccountNumber() string {
	return time.Now().Format("200601021504") + fmt.Sprintf("%03d", rand.Intn(999)+1)
}
