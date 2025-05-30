package main

import (
	"context"
	"log"
	"net"

	"github.com/abaika-abay/live_sports_project/common/pkg/config"
	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"github.com/abaika-abay/live_sports_project/user-service/service"
	pb "github.com/abaika-abay/live_sports_project/user-service/proto"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to MongoDB
	mongoDB, err := db.InitMongoDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Disconnect(context.Background())

	// Initialize user service
	userService := service.NewUserService(mongoDB.Database)

	// Start gRPC server
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	log.Println("User Service starting on " + cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}