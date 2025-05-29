package main

import (
	"context"
	"log"
	"net"

	"github.com/abaika-abay/live_sports_project/common/pkg/config"
	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"github.com/abaika-abay/live_sports_project/common/pkg/logger"
	pb "github.com/abaika-abay/live_sports_project/user-service/proto"
	"github.com/abaika-abay/live_sports_project/user-service/service"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.InitLogger(cfg.Log.Level)

	// Connect to MongoDB
	mongoDB, err := db.NewMongoDB(cfg.Mongo.URI, cfg.Mongo.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Disconnect(context.Background())

	// Initialize user service
	userService := service.NewUserService(mongoDB.Database) // Use Database field directly

	// Start gRPC server
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	logger.InfoLogger.Println("User Service starting on " + cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
