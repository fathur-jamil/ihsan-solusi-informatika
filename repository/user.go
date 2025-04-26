package repository

import (
	"github.com/jinzhu/gorm"

	"account_service/model"
)

type UserRepository interface {
	CountByNIKOrPhoneNumber(nik, phone string) (int64, error)
	Create(user *model.User) error
	DeleteByID(id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CountByNIKOrPhoneNumber(nik, phone string) (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("nik = ? OR phone_number = ?", nik, phone).Count(&count).Error
	return count, err
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) DeleteByID(id int64) error {
	return r.db.Delete(&model.User{}, id).Error
}
