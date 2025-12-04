package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/gabezy/go-grpc/internal/database"
	"github.com/gabezy/go-grpc/internal/pb"
	"github.com/gabezy/go-grpc/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDB := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDB)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	// Open TCP port to listen the gRPC default port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	fmt.Println("Running on port 50051")

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
