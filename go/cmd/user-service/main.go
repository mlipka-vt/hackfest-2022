package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"kafka-go-client/user"
	"net"
)

type Service struct {
	user.UnimplementedUserServiceServer
}

func (s Service) Login(ctx context.Context, request *user.LoginRequest) (*user.LoginResponse, error) {
	fmt.Println(fmt.Sprintf("INFO login request for %s", request.Username))
	return &user.LoginResponse{Status: "success"}, nil
}

func (s Service) Logout(ctx context.Context, request *user.LogoutRequest) (*user.LogoutResponse, error) {
	fmt.Println(fmt.Sprintf("INFO logout request for %s", request.Username))
	return &user.LogoutResponse{Status: "success"}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":50051"))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &Service{})

	fmt.Println("INFO user-service listening")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
