package grpc

import (
	"context"
	"github.com/ArtemusCoder/kvado-test-project/pkg/api"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

// RunServer запускает gRPC сервер с LibraryService сервисом
func RunServer(ctx context.Context, v1API api.LibraryServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// регистрируем сервис
	server := grpc.NewServer()
	api.RegisterLibraryServiceServer(server, v1API)

	// Добавляем возможность Server Reflection
	reflection.Register(server)

	// Остановка graceful stop сервера комбинацией Ctrl C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// запуск gRPC сервера
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
