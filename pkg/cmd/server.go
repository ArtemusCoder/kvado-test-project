package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/ArtemusCoder/kvado-test-project/pkg/service"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/ArtemusCoder/kvado-test-project/pkg/protocol/grpc"
)

// Config - конфигурация для Server
type Config struct {
	// Раздел параметров запуска сервера gRPC
	// gRPCPort — это порт TCP для прослушивания сервером gRPC.
	GRPCPort string

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

// RunServer - запускает gRPC сервер и HTTP-шлюз
func RunServer() error {
	ctx := context.Background()
	// Получаем конфигурацию
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "8080", "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "library-user", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "password", "Database password")
	flag.StringVar(&cfg.DatastoreDBName, "db-name", "librarydb", "Database name")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// Составление url к подключению базы данных
	dsn := fmt.Sprintf("%s:%s@tcp(%s:9000)/%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := service.NewLibraryServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
