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

func (l *LibraryService) InsertBook(title, category, publisher string, book_count int, price float32, fine_per_day float32, author_id int) (*data.Book, error) {
	book_to_insert := data.Book{
		Title:      title,
		Category:   category,
		Publisher:  publisher,
		BookCount:  book_count,
		Price:      price,
		FinePerDay: fine_per_day,
		AuthorId:   author_id,
	}
	book, err := l.model.Book.InsertBook(book_to_insert)

	if err != nil {
		return nil, err
	}
	return book, nil
}

func (l *LibraryService) UpdateBook(book_id int, input_json map[string]any) (*data.Book, error) {

	book, err := l.model.Book.UpdateBook(book_id, input_json)

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

func (l *LibraryService) InsertAuthor(name, about string) (*data.Author, error) {
	authorInput := data.Author{
		Name:  name,
		About: about,
	}

	// Add the author.
	addedAuthor, err := l.model.Author.InsertAuthor(authorInput)
	if err != nil {
		return nil, err
	}

	return addedAuthor, nil
}

func (l *LibraryService) GetAuthor(id int, name string) ([]data.Author, error) {
	var output_authors []data.Author

	// Get the author with id.
	if id != 0 {
		output_author, err := l.model.Author.GetAuthorWithId(id)
		if err != nil {
			return nil, err
		}
		output_authors = append(output_authors, output_author)
		return output_authors, nil

	} else {

		output_authors, err := l.model.Author.GetAuthorWithDetails(name)

		if err != nil {
			return nil, err
		}
		return output_authors, nil
	}
}

func (l *LibraryService) GetBooks(input_json map[string]any) ([]*data.Book_with_name, error) {
	book_list, err := l.model.Book.GetBook(input_json)

	if err != nil {
		return nil, err
	}
	return book_list, nil
}
