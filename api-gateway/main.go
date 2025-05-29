package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	apipb "github.com/abaika-abay/live_sports_project/api-gateway/proto"
	matchpb "github.com/abaika-abay/live_sports_project/match-service/proto"
	userpb "github.com/abaika-abay/live_sports_project/user-service/proto"
)

type apiGatewayServer struct {
	apipb.UnimplementedApiGatewayServiceServer
	userClient  userpb.UserServiceClient
	matchClient matchpb.MatchServiceClient
}

func (s *apiGatewayServer) RegisterUser(ctx context.Context, req *apipb.RegisterUserRequest) (*apipb.RegisterUserResponse, error) {
	res, err := s.userClient.Register(ctx, &userpb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &apipb.RegisterUserResponse{
		UserId:  res.UserId,
		Message: res.Message,
		Success: res.Success,
	}, nil
}

func (s *apiGatewayServer) CreateMatch(ctx context.Context, req *apipb.CreateMatchRequest) (*apipb.CreateMatchResponse, error) {
	res, err := s.matchClient.CreateMatch(ctx, &matchpb.CreateMatchRequest{
		MatchId:  req.MatchId,
		HomeTeam: req.HomeTeam,
		AwayTeam: req.AwayTeam,
	})
	if err != nil {
		return nil, err
	}

	return &apipb.CreateMatchResponse{
		MatchId: res.MatchId,
		Status:  res.Status,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	userConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to UserService: %v", err)
	}
	userClient := userpb.NewUserServiceClient(userConn)

	matchConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to MatchService: %v", err)
	}
	matchClient := matchpb.NewMatchServiceClient(matchConn)

	apipb.RegisterApiGatewayServiceServer(grpcServer, &apiGatewayServer{
		userClient:  userClient,
		matchClient: matchClient,
	})

	log.Println("API Gateway running on port 5000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
