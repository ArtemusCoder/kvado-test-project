package test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ArtemusCoder/kvado-test-project/pkg/service"
	"log"
	"net"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ArtemusCoder/kvado-test-project/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type Config struct {
	// Раздел параметров Базы Данных
	// DatastoreDBHost - хост базы данных
	DatastoreDBHost string
	// DatastoreDBUser - пользователь для подключения к базе данных
	DatastoreDBUser string
	// DatastoreDBPassword - пароль для подключения к базе данных
	DatastoreDBPassword string
	// DatastoreDBName - название базы данных
	DatastoreDBName string
}

func server(ctx context.Context) (api.LibraryServiceClient, func()) {
	buffer := 101024 * 101024
	lis := bufconn.Listen(buffer)

	var cfg Config
	cfg.DatastoreDBHost = "localhost"
	cfg.DatastoreDBUser = "library-user"
	cfg.DatastoreDBPassword = "password"
	cfg.DatastoreDBName = "librarydb"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:9000)/%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("failed to open database: %v", err)
	}

	baseServer := grpc.NewServer()
	api.RegisterLibraryServiceServer(baseServer, service.NewLibraryServiceServer(db))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("Error with serving server: %s", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
		db.Close()
	}

	client := api.NewLibraryServiceClient(conn)

	return client, closer
}

// Проверка GetBooksByAuthor endpoint
func TestLibraryServer_GetBooksByAuthor(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *api.GetBooksByAuthorResponse
		err error
	}

	tests := map[string]struct {
		in       *api.GetBooksByAuthorRequest
		expected expectation
	}{
		"Get Books": {
			in: &api.GetBooksByAuthorRequest{
				Name: "Eric Freeman",
			},
			expected: expectation{
				out: &api.GetBooksByAuthorResponse{
					Book: []*api.Book{
						{
							Title: "Head First Design Patterns",
						},
						{
							Title: "Head First JavaScript Programming: A Brain-Friendly Guide",
						},
					},
				},
				err: nil,
			},
		},
		"No such author": {
			in: &api.GetBooksByAuthorRequest{
				Name: "Eric",
			},
			expected: expectation{
				out: &api.GetBooksByAuthorResponse{
					Book: []*api.Book{},
				},
				err: nil,
			},
		},
		"One book from author": {
			in: &api.GetBooksByAuthorRequest{
				Name: "Robert C. Martin",
			},
			expected: expectation{
				out: &api.GetBooksByAuthorResponse{
					Book: []*api.Book{
						{
							Title: "Clean Code",
						},
					},
				},
				err: nil,
			},
		},
		"Blank value": {
			in: &api.GetBooksByAuthorRequest{
				Name: "",
			},
			expected: expectation{
				out: &api.GetBooksByAuthorResponse{
					Book: []*api.Book{},
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetBooksByAuthor(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if len(tt.expected.out.Book) != len(out.Book) {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
				if len(tt.expected.out.Book) != 0 {
					for i := range tt.expected.out.Book {
						if tt.expected.out.Book[i].Title != out.Book[i].Title {
							t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
						}
					}
				}
			}
		})
	}
}

// Проверка GetAuthorsByBook endpoint
func TestLibraryServer_GetAuthorsByBook(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *api.GetAuthorByBookResponse
		err error
	}

	tests := map[string]struct {
		in       *api.GetAuthorsByBookRequest
		expected expectation
	}{
		"Get few Authors": {
			in: &api.GetAuthorsByBookRequest{
				Title: "Head First Design Patterns",
			},
			expected: expectation{
				out: &api.GetAuthorByBookResponse{
					Author: []*api.Author{
						{
							Name: "Eric Freeman",
						},
						{
							Name: "Elizabeth Freeman",
						},
					},
				},
				err: nil,
			},
		},
		"No such book": {
			in: &api.GetAuthorsByBookRequest{
				Title: "TestBook",
			},
			expected: expectation{
				out: &api.GetAuthorByBookResponse{
					Author: []*api.Author{},
				},
				err: nil,
			},
		},
		"One author from book": {
			in: &api.GetAuthorsByBookRequest{
				Title: "Clean Code",
			},
			expected: expectation{
				out: &api.GetAuthorByBookResponse{
					Author: []*api.Author{
						{
							Name: "Robert C. Martin",
						},
					},
				},
				err: nil,
			},
		},
		"Blank value": {
			in: &api.GetAuthorsByBookRequest{
				Title: "",
			},
			expected: expectation{
				out: &api.GetAuthorByBookResponse{
					Author: []*api.Author{},
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetAuthorsByBook(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if len(tt.expected.out.Author) != len(out.Author) {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
				if len(tt.expected.out.Author) != 0 {
					for i := range tt.expected.out.Author {
						if tt.expected.out.Author[i].Name != out.Author[i].Name {
							t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
						}
					}
				}
			}
		})
	}
}
