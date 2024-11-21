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
	"github.com/bioform/go-web-app-template/pkg/util/crypt"
)

type UserRepository interface {
	Create(user *model.User) (uint, error)
	FindByID(id uint) (*model.User, error)
	FindByEmailAndPassword(email, password string) (*model.User, error)
}

type userRepositoryImpl struct {
	// database connection or ORM
	db  *gorm.DB
	ctx context.Context
}

func NewUserRepository(ctx context.Context) *userRepositoryImpl {
	return &userRepositoryImpl{
		db:  database.Get(ctx),
		ctx: ctx,
	}
}

func (r *userRepositoryImpl) Create(user *model.User) (uint, error) {

	if len(user.Password) > 0 {
		// Password was updated, hash it
		hashedPassword, err := crypt.HashPassword(user.Password)
		if err != nil {
			return 0, err
		}
		user.PasswordHash = hashedPassword
		user.Password = "" // Clear the plain password after hashing
	}

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

func (r *userRepositoryImpl) FindByID(id uint) (*model.User, error) {
	// Logic to retrieve a user by ID from the database
	user := &model.User{}
	if err := r.db.First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found. %w", id, ErrRecordNotFound)
		}
		return nil, err // Other database error
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByEmailAndPassword(email, password string) (_ *model.User, err error) {
	defer util.WrapError(&err, "repository.FindByEmailAndPassword(%q)", email)

	// Logic to retrieve a user by email from the database
	user := &model.User{}
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
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
