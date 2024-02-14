package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	"github.com/alexptdev/auth-api/internal/config"
	"github.com/alexptdev/auth-api/internal/config/env"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	desc "github.com/alexptdev/auth-api/pkg/user_v1"
)

var configPath string

type userServer struct {
	desc.UnimplementedUserV1Server
	conPool *pgxpool.Pool
}

func init() {
	flag.StringVar(
		&configPath,
		"config-path",
		".env",
		"path to config file",
	)
}

func (s *userServer) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	insertBuilder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("user_name", "user_email", "user_password", "user_role").
		Values(req.Name, req.Email, req.Password, req.Role).
		Suffix("RETURNING user_id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v \n", err)
		return nil, err
	}

	var userId int64
	err = s.conPool.QueryRow(ctx, query, args...).Scan(&userId)
	if err != nil {
		log.Printf("failed to create user: %v \n", err)
		return nil, err
	}

	return &desc.CreateResponse{Id: userId}, nil
}

func (s *userServer) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	selectBuilder := sq.Select(
		"user_id",
		"user_name",
		"user_email",
		"user_role",
		"user_created_at",
		"user_updated_at",
	).
		From("users").
		Where(sq.Eq{"user_id": req.GetId()}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v \n", err)
		return nil, err
	}

	var id int64
	var name, email string
	var role int32
	var createdAt time.Time
	var updatedAt sql.NullTime

	row := s.conPool.QueryRow(ctx, query, args...)
	err = row.Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("failed to select user: %v \n", err)
		return nil, err
	}

	var updatedAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      desc.UserRole(role),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: updatedAtTime,
	}, nil
}

func (s *userServer) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {

	updateBuilder := sq.Update("users")

	if req.Name != nil {
		updateBuilder = updateBuilder.Set("user_name", req.Name.Value)
	}

	if req.Email != nil {
		updateBuilder = updateBuilder.Set("user_email", req.Email.Value)
	}

	updateBuilder = updateBuilder.Set("user_updated_at", sq.Expr("now()"))

	updateBuilder = updateBuilder.Where(sq.Eq{"user_id": req.Id})

	query, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Printf("failed to build query: %v \n", err)
		return nil, err
	}

	_, err = s.conPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update user: %v \n", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *userServer) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {

	deleteBuilder := sq.Delete("users").
		Where(sq.Eq{"user_id": req.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v \n", err)
		return nil, err
	}

	_, err = s.conPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete user: %v \n", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func main() {

	ctx := context.Background()

	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGrpcConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPgConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	listen, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("Failed to listen port: %v", err)
	}

	conPool, err := pgxpool.Connect(ctx, pgConfig.Dsn())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conPool.Close()

	server := grpc.NewServer()
	reflection.Register(server)
	desc.RegisterUserV1Server(server, &userServer{
		conPool: conPool,
	})

	log.Printf("Server listenenig at %v", listen.Addr())
	if err = server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
