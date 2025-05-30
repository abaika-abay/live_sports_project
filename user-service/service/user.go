package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	pb "github.com/abaika-abay/live_sports_project/user-service/proto"
	"github.com/abaika-abay/live_sports_project/user-service/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		repo: repository.NewUserRepository(db),
	}
}
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "all fields are required")
	}

	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	if existingUser != nil {
		return nil, status.Errorf(codes.AlreadyExists, "email already registered")
	}

	userID := primitive.NewObjectID().Hex()
	user := &repository.User{
		UserID:   userID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.RegisterResponse{
		UserId:  userID,
		Message: "User registered successfully",
		Success: true,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	if user == nil || user.Password != req.Password {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	// Simple token (in production, use JWT)
	token := primitive.NewObjectID().Hex()

	return &pb.LoginResponse{
		UserId:  user.UserID,
		Token:   token,
		Success: true,
		Message: "Login successful",
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	user, err := s.repo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.GetProfileResponse{
		UserId:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Success:   true,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.GetProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	update := bson.M{}
	if req.Username != "" {
		update["username"] = req.Username
	}
	if req.Email != "" {
		update["email"] = req.Email
	}

	if len(update) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no fields to update")
	}

	if err := s.repo.UpdateUser(ctx, req.UserId, update); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	user, err := s.repo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return &pb.GetProfileResponse{
		UserId:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Success:   true,
	}, nil
}
