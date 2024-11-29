package repository

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/util"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (uint, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error)
}

type userRepositoryImpl struct{}

func NewUserRepository() *userRepositoryImpl {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) (uint, error) {
	db := database.Get(ctx)

	err := db.Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, &model.EmailDuplicateError{Email: user.Email}
		}
		return 0, err
	}
	return user.ID, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.User, error) {
	db := database.Get(ctx)
	// Logic to retrieve a user by ID from the database
	user := &model.User{}
	if err := db.First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found. %w", id, ErrRecordNotFound)
		}
		return nil, err // Other database error
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByEmailAndPassword(ctx context.Context, email, password string) (_ *model.User, err error) {
	defer util.WrapError(&err, "repository.FindByEmailAndPassword(%q)", email)

	db := database.Get(ctx)
	// Logic to retrieve a user by email from the database
	user := &model.User{}
	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err // Other database error
	}

	// Compare the password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	return user, nil
}
