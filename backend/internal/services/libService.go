package services

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/data"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/utils"
)

type LibraryService struct {
	model data.Models
}

func NewLibraryService(db *sql.DB) *LibraryService {
	return &LibraryService{
		model: data.New(db),
	}
}

func (l *LibraryService) GetBook(id int) (*data.Book, error) {

	book, err := l.model.Book.GetBookWithId(id)

	if err != nil {
		return nil, err
	}
	return book, nil
}

func (l *LibraryService) RegisterUser(name, email, password, phone_number string, is_active, is_admin bool) (*data.User, error) {
	err := utils.ValidatePassword(password)

	if err != nil {
		return nil, err
	}

	hashed_password, err := utils.HashPassword(password)

	if err != nil {
		return nil, err
	}

	userInput := data.User{
		Name:        name,
		Email:       email,
		Password:    hashed_password,
		PhoneNumber: phone_number,
		IsActive:    is_active,
		IsAdmin:     is_admin,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	user, err := l.model.User.CreateUser(userInput)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (l *LibraryService) LoginUser(email, password string) (string, error) {
	userInput := data.User{
		Email: email,
	}

	user, err := l.model.User.GetUserWithEmail(userInput)

	if err != nil {
		return "", err
	}

	if is_same := utils.CheckPasswordHash(password, user.Password); !is_same {
		return "", errors.New("Invalid password")
	}

	if err != nil {
		return "", err
	}

	token, err := utils.CreateToken(strconv.Itoa(user.ID), user.Email, user.IsAdmin)

	if err != nil {
		return "", err
	}

	return token, nil
}
