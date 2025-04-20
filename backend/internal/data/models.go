package data

import (
	"context"
	"database/sql"
	"time"
)

var db *sql.DB

const dbTimeout = time.Second * 3

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		Author: Author{},
	}
}

type Models struct {
	Author Author
}

type Author struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	About     string    `json:"about"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Get author with id
func (a *Author) getAuthorWithId(id int) (*Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var author Author

	query := `select * from author where id = $1`
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&author.ID,
		&author.Name,
		&author.About,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &author, nil
}

// Create author
func (a *Author) insertAuthor(author Author) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	var newId int

	stmt := `insert into author (name, about, created_at, updated_at) values ($1, $2, $3, $4) returning id`

	row := db.QueryRowContext(ctx, stmt, author.Name, author.About, time.Now(), time.Now())

	err := row.Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}
