package grpcclient

import (
	"context"
	"fmt"

	pb "github.com/commitshark/notification-svc/gen"
	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/domain/ports"

	"google.golang.org/grpc"
)

type userDataGRPCClient struct {
	client pb.GrpcUserServiceClient
}

func NewUserDataGRPCClient(conn *grpc.ClientConn) ports.UserDataAdapter {
	return &userDataGRPCClient{
		client: pb.NewGrpcUserServiceClient(conn),
	}
}

func (c *userDataGRPCClient) GetContactInfo(ctx context.Context, userID string) (*domain.UserContactInfo, error) {
	req := &pb.GetUserContactInfoRequest{
		UserAuthId: userID,
	}

	resp, err := c.client.GetUserContactInfo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user contact info: %w", err)
	}

	return &domain.UserContactInfo{
		Email:    resp.Email,
		Phone:    resp.Phone,
		DeviceID: resp.Device,
	}, nil
}
