// internal/infra/grpc/events/client.go
package events

import (
	"context"

	eventspb "auth-service/internal/infra/grpc/events/pb/proto"
	"auth-service/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	events eventspb.EventsClient
}

func NewClient() (*Client, error) {
	addr := utils.ConfigInstance.EventsServiceURL

	cc, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   cc,
		events: eventspb.NewEventsClient(cc),
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) CreateEvent(ctx context.Context, req *eventspb.CreateEventRequest) (*eventspb.Event, error) {
	return c.events.CreateEvent(ctx, req)
}

func (c *Client) GetEvent(ctx context.Context, req *eventspb.GetEventRequest) (*eventspb.Event, error) {
	return c.events.GetEvent(ctx, req)
}

func (c *Client) ListEvents(ctx context.Context, req *eventspb.ListEventsRequest) (*eventspb.ListEventsResponse, error) {
	return c.events.ListEvents(ctx, req)
}

func (c *Client) UpdateEvent(ctx context.Context, req *eventspb.UpdateEventRequest) (*eventspb.Event, error) {
	return c.events.UpdateEvent(ctx, req)
}

func (c *Client) DeleteEvent(ctx context.Context, req *eventspb.DeleteEventRequest) (*eventspb.DeleteEventResponse, error) {
	return c.events.DeleteEvent(ctx, req)
}

func (c *Client) CutEventFrom(ctx context.Context, req *eventspb.CutEventFromRequest) (*eventspb.Event, error) {
	return c.events.CutEventFrom(ctx, req)
}

func (c *Client) ListOccurrences(ctx context.Context, req *eventspb.ListOccurrencesRequest) (*eventspb.ListOccurrencesResponse, error) {
	return c.events.ListOccurrences(ctx, req)
}
