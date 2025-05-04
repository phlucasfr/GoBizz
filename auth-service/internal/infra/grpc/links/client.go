package links

import (
	"context"

	"auth-service/internal/infra/grpc/links/pb/proto"
	"auth-service/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	connRead   *grpc.ClientConn
	connWrite  *grpc.ClientConn
	linksRead  proto.LinksServiceReadClient
	linksWrite proto.LinksServiceWriteClient
}

// NewClient initializes and returns a new instance of Client, which provides
// gRPC connections for both reading and writing to the Links service.
// It establishes two separate gRPC clients using the URLs specified in the
// configuration: one for read operations and another for write operations.
//
// Returns:
//   - *Client: A pointer to the initialized Client instance.
//   - error: An error if the gRPC client connections could not be established.
func NewClient() (*Client, error) {
	read, err := grpc.NewClient(utils.ConfigInstance.LinksServiceReadUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	write, err := grpc.NewClient(utils.ConfigInstance.LinksServiceWriteUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		linksRead:  proto.NewLinksServiceReadClient(read),
		linksWrite: proto.NewLinksServiceWriteClient(write),
		connRead:   read,
		connWrite:  write,
	}, nil
}

func (c *Client) CloseRead() error {
	return c.connRead.Close()
}

func (c *Client) CloseWrite() error {
	return c.connWrite.Close()
}

func (c *Client) CreateLink(ctx context.Context, request *proto.CreateLinkRequest) (*proto.CreateLinkResponse, error) {
	return c.linksWrite.CreateLink(ctx, request)
}

func (c *Client) GetLink(ctx context.Context, request *proto.GetLinkRequest) (*proto.GetLinkResponse, error) {
	return c.linksRead.GetLink(ctx, request)
}

func (c *Client) GetCustomerLinks(ctx context.Context, request *proto.GetCustomerLinksRequest) (*proto.GetCustomerLinksResponse, error) {
	return c.linksRead.GetCustomerLinks(ctx, request)
}

func (c *Client) DeleteLink(ctx context.Context, request *proto.DeleteLinkRequest) (*proto.DeleteLinkResponse, error) {
	return c.linksWrite.DeleteLink(ctx, request)
}

func (c *Client) UpdateLink(ctx context.Context, request *proto.UpdateLinkRequest) (*proto.UpdateLinkResponse, error) {
	return c.linksWrite.UpdateLink(ctx, request)
}

func (c *Client) UpdateLinkClicks(ctx context.Context, request *proto.UpdateLinkClicksRequest) (*proto.UpdateLinkClicksResponse, error) {
	return c.linksWrite.UpdateLinkClicks(ctx, request)
}
