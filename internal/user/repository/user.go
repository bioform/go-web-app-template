package repository

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/bioform/go-web-app-template/pkg/util"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (uint, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) (uint, error) {
	db := r.db

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
	db := r.db
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

	db := r.db
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

// Check if the email exists in the database
func (r *userRepositoryImpl) IsEmailUnique(ctx context.Context, email string) bool {

	db := r.db

	var count int

	err := db.Exec("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)

	if err != nil {

		// handle error

		return false

	}

	return count == 0

}
