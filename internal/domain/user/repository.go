package user

import (
	"errors"

	"gorm.io/gorm"
)

var ErrorAlreadyExists = errors.New("User with this email already exists")

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(user *User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrorAlreadyExists
		}
	}
	return result.Error
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User
	result := r.db.Where("email = ?", email).First(&user)
	// result := r.db.Where(&User{Email: email}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
