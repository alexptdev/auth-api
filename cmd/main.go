package main

import (
	"context"
	"fmt"
	desc "github.com/alexptdev/auth-api/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

const grpcPort = 4000

type userServer struct {
	desc.UnimplementedUserV1Server
}

func (s *userServer) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	fmt.Println(req.Name)
	fmt.Println(req.Email)
	fmt.Println(req.Password)
	fmt.Println(req.PasswordConfirm)
	fmt.Println(req.Role)

	return &desc.CreateResponse{Id: 1}, nil
}

func (s *userServer) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	fmt.Printf("Id := %d", req.Id)

	return &desc.GetResponse{
		Id:        1,
		Name:      "Test",
		Email:     "email",
		Role:      0,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func (s *userServer) Update(_ context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {

	fmt.Printf("Id := %d \n", req.Id)

	if req.Name != nil {
		fmt.Println(req.Name.Value)
	}

	if req.Email != nil {
		fmt.Println(req.Email.Value)
	}

	fmt.Println(req.Role)

	return &emptypb.Empty{}, nil
}

func (s *userServer) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {

	fmt.Printf("Id := %d", req.Id)

	return &emptypb.Empty{}, nil
}

func main() {

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen port: %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)
	desc.RegisterUserV1Server(server, &userServer{})

	log.Printf("Server listenenig at %v", listen.Addr())
	if err = server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
