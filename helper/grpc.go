package helper

import (
	"context"
	grpcserver "github.com/gemm123/vkrf-ticket/internal/grpc"
	"google.golang.org/grpc"
	"log"
)

func GetUserByEmailGrpc(conn *grpc.ClientConn, email string) (*grpcserver.GetUserByEmailResponse, error) {
	c := grpcserver.NewUserServiceClient(conn)
	userRequest := grpcserver.GetUserByEmailRequest{
		Email: email,
	}
	resp, err := c.GetUserByEmail(context.Background(), &userRequest)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	return resp, nil
}

func GetUserByUserIdGrpc(conn *grpc.ClientConn, userId string) (*grpcserver.GetUserByUserIdResponse, error) {
	c := grpcserver.NewUserServiceClient(conn)
	userRequest := grpcserver.GetUserByUserIdRequest{
		UserId: userId,
	}
	resp, err := c.GetUserByUserId(context.Background(), &userRequest)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	return resp, nil
}
