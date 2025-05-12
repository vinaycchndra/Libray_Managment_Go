package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var db *sql.DB

const dbTimeout = time.Second * 3

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		Author: &Author{},
		Book:   &Book{},
		User:   &User{},
	}
}

type Models struct {
	Author *Author
	Book   *Book
	User   *User
}

type Author struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	About     string    `json:"about"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Book struct {
	ID         int       `json:"int"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	Publisher  string    `json:"pubisher"`
	BookCount  int       `json:"book_count"`
	Price      float32   `json:"price"`
	FinePerDay float32   `json:"fine_per_day"`
	AuthorId   int       `json:"author_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Archive    bool      `json:"archive"`
}

type User struct {
	ID          int       `json:"int"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
	IsAdmin     bool      `json:"is_admin"`
}

// Get author with id
func (a *Author) GetAuthorWithId(id int) (*Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var author Author

	query := `select * from author where id = $1;`
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
func (a *Author) InsertAuthor(author Author) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	var newId int

	stmt := `insert into author (name, about, created_at, updated_at) values ($1, $2, $3, $4) returning id;`

	row := db.QueryRowContext(ctx, stmt, author.Name, author.About, time.Now(), time.Now())

	err := row.Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

// Create a book
func (b *Book) InsertBook(book Book) (*Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout*2)

	defer cancel()

	//Category check for the book
	var category_exists bool
	category_check_query := `select case when count(*) > 0 then True else False end from category where category_name = $1;`

	row := db.QueryRowContext(ctx, category_check_query, book.Category)
	err := row.Scan(&category_exists)

	if err != nil {
		return nil, err
	}

	if !category_exists {
		return nil, errors.New(fmt.Sprintf("%v this category does not exists.", book.Category))
	}

	//author id check for the book
	var author_id_exists bool
	author_id_check_query := `select case when count(*) > 0 then True else False end from author where id = $1;`

	row = db.QueryRowContext(ctx, author_id_check_query, book.AuthorId)
	err = row.Scan(&author_id_exists)

	if err != nil {
		return nil, err
	}

	if !author_id_exists {
		return nil, errors.New(fmt.Sprintf("%v this author_id does not exists.", book.AuthorId))
	}

	// Inserting book into the db.
	var inserted_book Book
	stmt := `insert into book (title, category, publisher, book_count, price, fine_per_day, created_at, updated_at, author_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id, title, category,  publisher, book_count, price, fine_per_day, created_at, updated_at, author_id;`

	row = db.QueryRowContext(ctx, stmt, book.Title, book.Category, book.Publisher, book.BookCount, book.Price,
		book.FinePerDay, book.AuthorId, time.Now(), time.Now())

	err = row.Scan(
		&inserted_book.ID,
		&inserted_book.Title,
		&inserted_book.Category,
		&inserted_book.BookCount,
		&inserted_book.Price,
		&inserted_book.FinePerDay,
		&inserted_book.CreatedAt,
		&inserted_book.UpdatedAt,
		&inserted_book.AuthorId,
	)

	if err != nil {
		return nil, err
	}
	return &inserted_book, nil
}

// Get a book with id
func (b *Book) GetBookWithId(id int) (*Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var book Book

	query := `select * from book where id = $1;`

	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Category,
		&book.Publisher,
		&book.BookCount,
		&book.Price,
		&book.FinePerDay,
		&book.AuthorId,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.Archive,
	)

	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (u *User) CreateUser(user User) (*User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout*2)

	defer cancel()

	// User exists check
	var user_exists bool

	user_exists_query := `select case when count(*) > 0 then True else False end from users where email = $1;`
	row := db.QueryRowContext(ctx, user_exists_query, user.Email)
	err := row.Scan(&user_exists)

	if err != nil {
		return nil, err
	}

	if user_exists {
		return nil, errors.New(fmt.Sprintf("User with email: %s already exists", user.Email))
	}

	// Inserting the user into db
	var inserted_user User
	insert_stmt := `insert into users (name, email, password, phone_number, created_at, updated_at, is_active, is_admin) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id, name, email, password, phone_number, created_at, updated_at, is_active, is_admin;`

	row = db.QueryRowContext(ctx, insert_stmt,
		user.Name,
		user.Email,
		user.Password,
		user.PhoneNumber,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsActive,
		user.IsAdmin)

	err = row.Scan(
		&inserted_user.ID,
		&inserted_user.Name,
		&inserted_user.Email,
		&inserted_user.Password,
		&inserted_user.PhoneNumber,
		&inserted_user.CreatedAt,
		&inserted_user.UpdatedAt,
		&inserted_user.IsActive,
		&inserted_user.IsAdmin,
	)

	if err != nil {
		return nil, err
	}
	return &inserted_user, nil
}
