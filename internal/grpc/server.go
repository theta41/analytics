package grpc

import (
	"context"
	"fmt"
	"net"

	"gitlab.com/g6834/team41/analytics/internal/ports"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "gitlab.com/g6834/team41/analytics/api/events"
)

func StartServer(host string, events ports.Events) error {

	lis, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("grpc failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterAnalyticsServiceServer(s, &server{events: events})

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("grpc failed to serve: %w", err)
	}

	return nil
}

type server struct {
	events ports.Events
	pb.UnimplementedAnalyticsServiceServer
}

func (s *server) CreateTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskResponse, error) {
	_, err := s.events.CreateTask(in.ObjectId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method CreateTask %v", err)
	}
	return &pb.TaskResponse{
		Success: true,
	}, nil
}

func (s *server) FinishTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskResponse, error) {
	err := s.events.FinishTask(in.ObjectId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method FinishTask %v", err)
	}
	return &pb.TaskResponse{
		Success: true,
	}, nil
}

func (s *server) CreateLetter(ctx context.Context, in *pb.LetterRequest) (*pb.LetterResponse, error) {
	_, err := s.events.CreateLetter(in.ObjectId, in.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method CreateLetter %v", err)
	}
	return &pb.LetterResponse{
		Success: true,
	}, nil
}

func (s *server) AcceptedLetter(ctx context.Context, in *pb.LetterRequest) (*pb.LetterResponse, error) {
	err := s.events.AcceptedLetter(in.ObjectId, in.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method AcceptedLetter %v", err)
	}
	return &pb.LetterResponse{
		Success: true,
	}, nil
}

func (s *server) DeclinedLetter(ctx context.Context, in *pb.LetterRequest) (*pb.LetterResponse, error) {
	err := s.events.DeclinedLetter(in.ObjectId, in.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method DeclinedLetter %v", err)
	}
	return &pb.LetterResponse{
		Success: true,
	}, nil
}
