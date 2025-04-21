package services

import (
	"database/sql"

	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/data"
)

type LibraryService struct {
	model data.Models
}

func NewLibraryService(db *sql.DB) *LibraryService {
	return &LibraryService{
		model: data.New(db),
	}
}

func (l *LibraryService) getBook(id int) (*data.Book, error) {
	book, err := l.model.Book.GetBookWithId(id)

	if err != nil {
		return nil, err
	}
	return book, nil
}
