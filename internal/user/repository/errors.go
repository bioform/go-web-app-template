package repository

import (
	"fmt"

	"github.com/bioform/go-web-app-template/internal/errors"
)

var (
	ErrRecordNotFound         = errors.ErrRecordNotFound
	ErrInvalidEmailOrPassword = fmt.Errorf("invalid email or password: %w", ErrRecordNotFound)
)
