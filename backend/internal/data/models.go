package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sanggonlee/gosq"
)

var db *sql.DB

const dbTimeout = time.Second * 3

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		Author:         &Author{},
		Book:           &Book{},
		User:           &User{},
		BookBorrowList: &BookBorrowList{},
	}
}

type Models struct {
	Author         *Author
	Book           *Book
	User           *User
	BookBorrowList *BookBorrowList
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

type Book_with_name struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	Publisher  string    `json:"publisher"`
	Price      float32   `json:"price"`
	FinePerDay float32   `json:"fine_per_day"`
	BookCount  int       `json:"book_count"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BookBorrowList struct {
	ID        int            `json:"id"`
	DueDate   time.Time      `json:"due_date"`
	UserId    int            `json:"user_id"`
	Closed    bool           `json:"closed"`
	FinePaid  float32        `json:"fine_paid"`
	CreateAt  time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	BookList  []*BookBorrorw `json:"lended_books"`
}

type BookBorrorw struct {
	ID       int  `json:"id"`
	BookId   int  `json:"book_id"`
	ListId   int  `json:"list_id"`
	Returned bool `json:"returned"`
	Extended bool `json:"extended"`
}

// Get author with id
func (a *Author) GetAuthorWithId(id int) (Author, error) {
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
		return author, err
	}
	return author, nil
}

// Get authors details with name and about details
func (a *Author) GetAuthorWithDetails(name string) ([]Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*dbTimeout)

	defer cancel()

	var authors []Author

	if name != "" {
		stmt := `select * from author where name ilike $1 order by created_at desc;`
		name_param := "%" + strings.ToLower(name) + "%"
		rows, err := db.QueryContext(ctx, stmt, name_param)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var author Author
			err = rows.Scan(&author.ID, &author.Name, &author.About, &author.CreatedAt, &author.UpdatedAt)
			if err != nil {
				return nil, err
			}
			authors = append(authors, author)
		}
	} else {
		stmt := `select * from author order by created_at desc;`
		rows, err := db.QueryContext(ctx, stmt)
		for rows.Next() {
			var author Author
			err = rows.Scan(&author.ID, &author.Name, &author.About, &author.CreatedAt, &author.UpdatedAt)
			if err != nil {
				return nil, err
			}
			authors = append(authors, author)
		}
	}
	return authors, nil
}

// Create author
func (a *Author) InsertAuthor(author Author) (*Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	var addedAuthor Author

	stmt := `insert into author (name, about, created_at, updated_at) values ($1, $2, $3, $4) returning id, name, about, created_at, updated_at;`

	row := db.QueryRowContext(ctx, stmt, author.Name, author.About, time.Now(), time.Now())

	err := row.Scan(
		&addedAuthor.ID,
		&addedAuthor.Name,
		&addedAuthor.About,
		&addedAuthor.CreatedAt,
		&addedAuthor.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &addedAuthor, nil
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
	stmt := `insert into book (title, category, publisher, book_count, price, fine_per_day, created_at, updated_at, author_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id, title, category,  publisher, book_count, price, fine_per_day, created_at, updated_at, author_id;`

	row = db.QueryRowContext(ctx, stmt, book.Title, book.Category, book.Publisher, book.BookCount, book.Price,
		book.FinePerDay, time.Now(), time.Now(), book.AuthorId)
	// id, title, category,  publisher, book_count, price, fine_per_day, created_at, updated_at, author_id;
	err = row.Scan(
		&inserted_book.ID,
		&inserted_book.Title,
		&inserted_book.Category,
		&inserted_book.Publisher,
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

// Update book
func (b *Book) UpdateBook(book_id int, input_json map[string]any) (*Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout*3)
	defer cancel()

	type fields struct {
		Title        bool
		Category     bool
		Publisher    bool
		Book_count   bool
		Price        bool
		Fine_per_day bool
		Author_id    bool
	}

	field := fields{}

	query_count := 0

	query_args := make([]any, 0, 7)

	if title_value, ok := input_json["title"]; ok {
		field.Title = true
		query_count++
		query_args = append(query_args, title_value)
	}

	if category_value, ok := input_json["category"]; ok {
		field.Category = true
		query_count++
		query_args = append(query_args, category_value)

		//Category check for the book
		var category_exists bool
		category_check_query := `select case when count(*) > 0 then True else False end from category where category_name = $1;`

		row := db.QueryRowContext(ctx, category_check_query, category_value.(string))
		err := row.Scan(&category_exists)

		if err != nil {
			return nil, err
		}

		if !category_exists {
			return nil, errors.New(fmt.Sprintf("%v this category does not exists.", category_value))
		}

	}

	if publisher_value, ok := input_json["publisher"]; ok {
		field.Publisher = true
		query_count++
		query_args = append(query_args, publisher_value)
	}

	if book_count_value, ok := input_json["book_count"]; ok {
		field.Book_count = true
		query_count++
		query_args = append(query_args, int(book_count_value.(float64)))
	}

	if price_value, ok := input_json["price"]; ok {
		field.Price = true
		query_count++
		query_args = append(query_args, float32(price_value.(float64)))
	}

	if fine_per_day_value, ok := input_json["fine_per_day"]; ok {
		field.Fine_per_day = true
		query_count++
		query_args = append(query_args, float32(fine_per_day_value.(float64)))
	}

	if author_id_value, ok := input_json["author_id"]; ok {
		field.Author_id = true
		query_count++
		query_args = append(query_args, int(author_id_value.(float64)))
	}

	if !(field.Title || field.Category || field.Publisher || field.Book_count || field.Price || field.Fine_per_day || field.Author_id) {
		return nil, errors.New(fmt.Sprintf("Nonthing to update for book with book_id :: %v", book_id))
	}

	query_args = append(query_args, time.Now())
	query_args = append(query_args, book_id)
	query_count = query_count + 2

	annotation_list := make([]any, 0, query_count)

	for i := 1; i <= query_count; i++ {
		annotation_list = append(annotation_list, i)
	}

	query, err := gosq.Compile(`
				update book set
	 			{{ [if] .Title [then]  title = $%d, }}
				{{ [if] .Category [then] category = $%d, }}
				{{ [if] .Publisher [then] publisher = $%d, }}
				{{ [if] .Book_count [then] book_count = $%d, }}
				{{ [if] .Price [then] price = $%d, }}
				{{ [if] .Fine_per_day [then] fine_per_day = $%d, }} 
				{{ [if] .Author_id [then] author_id = $%d, }}  
				updated_at = $%d where id = $%d 
				returning id, title, category,  publisher, book_count, 
				price, fine_per_day, created_at, updated_at, author_id;
				`, field)

	query = fmt.Sprintf(query, annotation_list...)
	row := db.QueryRowContext(ctx, query, query_args...)

	var inserted_book Book

	err = row.Scan(
		&inserted_book.ID,
		&inserted_book.Title,
		&inserted_book.Category,
		&inserted_book.Publisher,
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

func (b *Book) GetBook(input_json map[string]any) ([]*Book_with_name, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout*3)
	defer cancel()

	type fields struct {
		Title      bool
		Category   bool
		Publisher  bool
		AuthorName bool
	}
	fmt.Println("here is 1")
	var query string
	var err error

	field := fields{}
	query_count := 0
	query_args := make([]any, 0, 4)
	results := make([]*Book_with_name, 0)

	if title_value, ok := input_json["title"]; ok {
		field.Title = true
		query_count++
		query_args = append(query_args, "%"+title_value.(string)+"%")
	}

	if category_value, ok := input_json["category"]; ok {
		field.Category = true
		query_count++
		query_args = append(query_args, "%"+category_value.(string)+"%")

	}

	if publisher_value, ok := input_json["publisher"]; ok {
		field.Publisher = true
		query_count++
		query_args = append(query_args, "%"+publisher_value.(string)+"%")
	}

	if author_id_value, ok := input_json["author_name"]; ok {
		field.AuthorName = true
		query_count++
		query_args = append(query_args, "%"+author_id_value.(string)+"%")
	}
	fmt.Println("here is 2")
	if !(field.Title || field.Category || field.Publisher || field.AuthorName) {
		query = `select t1.id, t1.title, t1.category, t1.publisher, t1.price, t1.fine_per_day,
					t1.book_count, t2.name, t1.created_at, t1.updated_at 
					from book as t1 inner join author as t2 on t1.author_id = t2.id order by t1.created_at desc;`
		fmt.Println(query, "inside the block")
	} else {
		annotation_list := make([]any, 0, query_count)

		for i := 1; i <= query_count; i++ {
			annotation_list = append(annotation_list, i)
		}

		query, err = gosq.Compile(`
					select t1.id, t1.title, t1.category, t1.publisher, t1.price, t1.fine_per_day,
					t1.book_count, t2.name, t1.created_at, t1.updated_at 
					from book as t1 inner join author as t2 on t1.author_id = t2.id where 1=1    
					{{ [if] .Title [then]     and t1.title ilike $%d }}
					{{ [if] .Category [then]  and t1.category ilike $%d }}
					{{ [if] .Publisher [then] and t1.publisher ilike $%d }}
					{{ [if] .AuthorName [then] and t2.name ilike $%d }} order by t1.created_at desc;`,
			field)
		if err != nil {
			return nil, err
		}
		query = fmt.Sprintf(query, annotation_list...)
	}
	fmt.Println(query, "outside the block")
	rows, err := db.QueryContext(ctx, query, query_args...)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var output_book Book_with_name

		err = rows.Scan(
			&output_book.ID,
			&output_book.Title,
			&output_book.Category,
			&output_book.Publisher,
			&output_book.Price,
			&output_book.FinePerDay,
			&output_book.BookCount,
			&output_book.AuthorName,
			&output_book.CreatedAt,
			&output_book.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		results = append(results, &output_book)
	}

	return results, nil
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

func (u *User) GetUserWithEmail(user User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	// User exists check
	var existing_user User

	get_user_query := `select * from users where email = $1;`

	row := db.QueryRowContext(ctx, get_user_query, user.Email)

	err := row.Scan(
		&existing_user.ID,
		&existing_user.Name,
		&existing_user.Email,
		&existing_user.Password,
		&existing_user.PhoneNumber,
		&existing_user.CreatedAt,
		&existing_user.UpdatedAt,
		&existing_user.IsActive,
		&existing_user.IsAdmin,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(fmt.Sprintf("User with email %s does not exist.", user.Email))
	}

	return &existing_user, nil
}

func (b *BookBorrowList) CreateBookBorrowList(input_json map[string]any) (*BookBorrowList, error) {
	return nil, nil
}
