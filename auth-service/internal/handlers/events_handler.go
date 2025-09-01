// internal/handlers/events_handler.go
package handlers

import (
	"context"
	"errors"
	"strings"

	eventsgw "auth-service/internal/infra/grpc/events"
	eventspb "auth-service/internal/infra/grpc/events/pb/proto"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type OccurrenceHTTP struct {
	EventId string `json:"eventId"`
	Name    string `json:"name"`
	Date    string `json:"date"`
}
type ListOccurrencesHTTPResp struct {
	Start       string           `json:"start"`
	End         string           `json:"end"`
	Occurrences []OccurrenceHTTP `json:"occurrences"`
}

type EventsHandler struct {
	eventsClient *eventsgw.Client
}

func NewEventsHandler(eventsClient *eventsgw.Client) *EventsHandler {
	return &EventsHandler{eventsClient: eventsClient}
}

func (h *EventsHandler) Close() error { return h.eventsClient.Close() }

func (h *EventsHandler) CreateEvent(ctx context.Context, req *eventspb.CreateEventRequest) (*eventspb.Event, error) {
	if req.CustomerId == "" || req.Name == "" || req.StartDate == "" || req.IntervalDays <= 0 {
		return nil, errors.New("customer_id, name, start_date and interval_days are required")
	}
	return h.eventsClient.CreateEvent(ctx, req)
}
func (h *EventsHandler) GetEvent(ctx context.Context, req *eventspb.GetEventRequest) (*eventspb.Event, error) {
	if req.CustomerId == "" || req.Id == "" {
		return nil, errors.New("customer_id and id are required")
	}
	return h.eventsClient.GetEvent(ctx, req)
}
func (h *EventsHandler) ListEvents(ctx context.Context, req *eventspb.ListEventsRequest) (*eventspb.ListEventsResponse, error) {
	if req.CustomerId == "" {
		return nil, errors.New("customer_id is required")
	}
	return h.eventsClient.ListEvents(ctx, req)
}
func (h *EventsHandler) UpdateEvent(ctx context.Context, req *eventspb.UpdateEventRequest) (*eventspb.Event, error) {
	if req.CustomerId == "" || req.Id == "" {
		return nil, errors.New("customer_id and id are required")
	}
	return h.eventsClient.UpdateEvent(ctx, req)
}
func (h *EventsHandler) DeleteEvent(ctx context.Context, req *eventspb.DeleteEventRequest) (*eventspb.DeleteEventResponse, error) {
	if req.CustomerId == "" || req.Id == "" {
		return nil, errors.New("customer_id and id are required")
	}
	return h.eventsClient.DeleteEvent(ctx, req)
}
func (h *EventsHandler) CutEventFrom(ctx context.Context, req *eventspb.CutEventFromRequest) (*eventspb.Event, error) {
	if req.CustomerId == "" || req.Id == "" || req.From == "" {
		return nil, errors.New("customer_id, id and from are required")
	}
	return h.eventsClient.CutEventFrom(ctx, req)
}
func (h *EventsHandler) ListOccurrences(ctx context.Context, req *eventspb.ListOccurrencesRequest) (*eventspb.ListOccurrencesResponse, error) {
	if req.CustomerId == "" || req.Start == "" || req.End == "" {
		return nil, errors.New("customer_id, start and end are required")
	}
	return h.eventsClient.ListOccurrences(ctx, req)
}

// ---------------- HTTP (Fiber) ----------------

func (h *EventsHandler) CreateEventHTTP(c *fiber.Ctx) error {
	var body struct {
		Name         string  `json:"name"`
		StartDate    string  `json:"startDate"`
		IntervalDays int32   `json:"intervalDays"`
		StopAt       *string `json:"stopAt"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}
	customerId := c.Params("customerId")
	if customerId == "" {
		if v := c.Locals("user_id"); v != nil {
			if s, ok := v.(string); ok {
				customerId = s
			}
		}
	}

	if customerId == "" || body.Name == "" || body.StartDate == "" || body.IntervalDays <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer_id, name, start_date and interval_days are required"})
	}

	req := &eventspb.CreateEventRequest{
		CustomerId:   customerId,
		Name:         body.Name,
		StartDate:    body.StartDate,
		IntervalDays: body.IntervalDays,
	}
	if body.StopAt != nil && *body.StopAt != "" {
		req.StopAt = body.StopAt
	}

	resp, err := h.CreateEvent(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *EventsHandler) GetEventHTTP(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}
	customerId := c.Params("customerId")
	if customerId == "" {
		if v := c.Locals("user_id"); v != nil {
			if s, ok := v.(string); ok {
				customerId = s
			}
		}
	}
	if customerId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer_id is required"})
	}

	req := &eventspb.GetEventRequest{CustomerId: customerId, Id: id}
	resp, err := h.GetEvent(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *EventsHandler) ListEventsHTTP(c *fiber.Ctx) error {
	customerId := c.Params("customerId")
	if customerId == "" {
		if v := c.Locals("user_id"); v != nil {
			if s, ok := v.(string); ok {
				customerId = s
			}
		}
	}
	if customerId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer_id is required"})
	}
	req := &eventspb.ListEventsRequest{CustomerId: customerId}
	if v := c.QueryInt("limit"); v > 0 {
		req.Limit = int64(v)
	}
	if v := c.QueryInt("offset"); v > 0 {
		req.Offset = int64(v)
	}

	resp, err := h.ListEvents(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *EventsHandler) ListOccurrencesHTTP(c *fiber.Ctx) error {
	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "start and end are required (YYYY-MM-DD)"})
	}

	customerId := c.Params("customerId")
	if customerId == "" {
		if v := c.Locals("user_id"); v != nil {
			if s, ok := v.(string); ok {
				customerId = s
			}
		}
	}
	if customerId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer_id is required"})
	}

	req := &eventspb.ListOccurrencesRequest{
		CustomerId: customerId,
		Start:      start,
		End:        end,
	}

	raw := strings.TrimSpace(c.Query("name"))
	if raw != "" && raw != "undefined" && raw != "null" {
		req.Name = wrapperspb.String(raw)
	}

	grpcResp, err := h.ListOccurrences(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	out := ListOccurrencesHTTPResp{
		Start:       grpcResp.Start,
		End:         grpcResp.End,
		Occurrences: make([]OccurrenceHTTP, len(grpcResp.Occurrences)),
	}
	for i, o := range grpcResp.Occurrences {
		out.Occurrences[i] = OccurrenceHTTP{
			EventId: o.EventId,
			Name:    o.Name,
			Date:    o.Date,
		}
	}

	return c.Status(fiber.StatusOK).JSON(out)
}

func (h *EventsHandler) UpdateEventHTTP(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}
	customerId := c.Params("customerId")
	if customerId == "" {
		if v := c.Locals("user_id"); v != nil {
			if s, ok := v.(string); ok {
				customerId = s
			}
		}
	}
	if customerId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer_id is required"})
	}

	var body struct {
		Name         *string `json:"name"`
		StartDate    *string `json:"startDate"`
		IntervalDays *int32  `json:"intervalDays"`
		StopAt       *string `json:"stopAt"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}

	req := &eventspb.UpdateEventRequest{
		CustomerId: customerId,
		Id:         id,
	}

	if body.Name != nil {
		req.Name = wrapperspb.String(*body.Name)
	}
	if body.StartDate != nil {
		req.StartDate = wrapperspb.String(*body.StartDate)
	}
	if body.IntervalDays != nil {
		req.IntervalDays = wrapperspb.Int32(*body.IntervalDays)
	}
	if body.StopAt != nil && *body.StopAt != "" {
		req.StopAtChange = &eventspb.UpdateEventRequest_StopAt{StopAt: *body.StopAt}
	}

	resp, err := h.UpdateEvent(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *EventsHandler) DeleteOrCutEventHTTP(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}
	customerId := c.Params("customerId")
	if customerId == "" {
		if v := c.Locals("user_id"); v != nil {
			if s, ok := v.(string); ok {
				customerId = s
			}
		}
	}
	if customerId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer_id is required"})
	}

	if from := c.Query("from"); from != "" {
		req := &eventspb.CutEventFromRequest{CustomerId: customerId, Id: id, From: from}
		resp, err := h.CutEventFrom(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}

	req := &eventspb.DeleteEventRequest{CustomerId: customerId, Id: id}
	resp, err := h.DeleteEvent(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
