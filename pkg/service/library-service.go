package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ArtemusCoder/kvado-test-project/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SqlQuery Общая часть SQL запроса
const (
	SqlQuery = "SELECT %s FROM book join author on book.author_id = author.id WHERE %s = ?"
)

// LibraryServiceServer реализация интерфейса
type LibraryServiceServer struct {
	api.UnimplementedLibraryServiceServer
	db *sql.DB
}

func NewLibraryServiceServer(db *sql.DB) api.LibraryServiceServer {
	return &LibraryServiceServer{db: db}
}

func (s *LibraryServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// GetBooksByAuthor реализация метода
func (s *LibraryServiceServer) GetBooksByAuthor(ctx context.Context, req *api.GetBooksByAuthorRequest) (*api.GetBooksByAuthorResponse, error) {
	// Подключение к базе данных
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	//
	query := fmt.Sprintf(SqlQuery, "book.title", "author.name")

	rows, err := c.QueryContext(ctx, query, req.Name)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Library-> "+err.Error())
	}
	defer rows.Close()

	var books []*api.Book
	for rows.Next() {
		book := new(api.Book)
		if err := rows.Scan(&book.Title); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from Library row-> "+err.Error())
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Library-> "+err.Error())
	}
	return &api.GetBooksByAuthorResponse{
		Book: books,
	}, nil
}

// GetAuthorsByBook реализация метода
func (s *LibraryServiceServer) GetAuthorsByBook(ctx context.Context, req *api.GetAuthorsByBookRequest) (*api.GetAuthorByBookResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	query := fmt.Sprintf(SqlQuery, "author.name", "book.title")
	rows, err := c.QueryContext(ctx, query, req.Title)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Library-> "+err.Error())
	}
	defer rows.Close()

	var authors []*api.Author
	for rows.Next() {
		author := new(api.Author)
		if err := rows.Scan(&author.Name); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from Library row-> "+err.Error())
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Library-> "+err.Error())
	}
	return &api.GetAuthorByBookResponse{
		Author: authors,
	}, nil
}
