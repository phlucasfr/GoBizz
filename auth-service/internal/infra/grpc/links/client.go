package links

import (
	"context"

	"auth-service/internal/infra/grpc/links/pb/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client proto.LinksServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: proto.NewLinksServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateLink(ctx context.Context, request *proto.CreateLinkRequest) (*proto.CreateLinkResponse, error) {
	return c.client.CreateLink(ctx, request)
}

func (c *Client) GetLink(ctx context.Context, request *proto.GetLinkRequest) (*proto.GetLinkResponse, error) {
	return c.client.GetLink(ctx, request)
}

func (c *Client) GetCustomerLinks(ctx context.Context, request *proto.GetCustomerLinksRequest) (*proto.GetCustomerLinksResponse, error) {
	return c.client.GetCustomerLinks(ctx, request)
}

func (c *Client) DeleteLink(ctx context.Context, request *proto.DeleteLinkRequest) (*proto.DeleteLinkResponse, error) {
	return c.client.DeleteLink(ctx, request)
}

func (c *Client) UpdateLink(ctx context.Context, request *proto.UpdateLinkRequest) (*proto.UpdateLinkResponse, error) {
	return c.client.UpdateLink(ctx, request)
}

func (c *Client) UpdateLinkClicks(ctx context.Context, request *proto.UpdateLinkClicksRequest) (*proto.UpdateLinkClicksResponse, error) {
	return c.client.UpdateLinkClicks(ctx, request)
}
