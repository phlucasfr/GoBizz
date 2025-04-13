package handlers

import (
	"context"
	"errors"

	"auth-service/internal/infra/grpc/links"
	"auth-service/internal/infra/grpc/links/pb/proto"

	"github.com/gofiber/fiber/v2"
)

type LinksHandler struct {
	linksClient *links.Client
}

func NewLinksHandler(linksClient *links.Client) *LinksHandler {
	return &LinksHandler{
		linksClient: linksClient,
	}
}

// HTTP Handlers
func (h *LinksHandler) CreateLinkHTTP(c *fiber.Ctx) error {
	var req proto.CreateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	resp, err := h.CreateLink(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *LinksHandler) GetLinkHTTP(c *fiber.Ctx) error {
	shortUrl := c.Params("shortUrl")
	if shortUrl == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "short_url is required",
		})
	}

	req := &proto.GetLinkRequest{
		ShortUrl: shortUrl,
	}

	resp, err := h.GetLink(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *LinksHandler) GetCustomerLinksHTTP(c *fiber.Ctx) error {
	customerId := c.Params("customerId")
	if customerId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "customer_id is required",
		})
	}

	req := &proto.GetCustomerLinksRequest{
		CustomerId: customerId,
	}

	if limit := c.QueryInt("limit"); limit > 0 {
		limit32 := int32(limit)
		req.Limit = &limit32
	}
	if offset := c.QueryInt("offset"); offset > 0 {
		offset32 := int32(offset)
		req.Offset = &offset32
	}

	if search := c.Query("search"); search != "" {
		req.Search = &search
	}
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if slugType := c.Query("slug_type"); slugType != "" {
		req.SlugType = &slugType
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		req.SortBy = &sortBy
	}
	if sortDirection := c.Query("sort_direction"); sortDirection != "" {
		req.SortDirection = &sortDirection
	}

	resp, err := h.GetCustomerLinks(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *LinksHandler) DeleteLinkHTTP(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	customerId := c.Locals("user_id")
	if customerId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	req := &proto.DeleteLinkRequest{
		Id:         id,
		CustomerId: customerId.(string),
	}

	resp, err := h.DeleteLink(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *LinksHandler) UpdateLinkHTTP(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	var req proto.UpdateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	req.Id = id

	resp, err := h.UpdateLink(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *LinksHandler) UpdateLinkClicksHTTP(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	req := &proto.UpdateLinkClicksRequest{
		Id: id,
	}

	resp, err := h.UpdateLinkClicks(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// gRPC Handlers
func (h *LinksHandler) CreateLink(ctx context.Context, req *proto.CreateLinkRequest) (*proto.CreateLinkResponse, error) {
	if req.OriginalUrl == "" {
		return nil, errors.New("original_url is required")
	}

	resp, err := h.linksClient.CreateLink(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *LinksHandler) GetLink(ctx context.Context, req *proto.GetLinkRequest) (*proto.GetLinkResponse, error) {
	if req.ShortUrl == "" {
		return nil, errors.New("short_url is required")
	}

	resp, err := h.linksClient.GetLink(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *LinksHandler) GetCustomerLinks(ctx context.Context, req *proto.GetCustomerLinksRequest) (*proto.GetCustomerLinksResponse, error) {
	if req.CustomerId == "" {
		return nil, errors.New("customer_id is required")
	}

	resp, err := h.linksClient.GetCustomerLinks(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *LinksHandler) DeleteLink(ctx context.Context, req *proto.DeleteLinkRequest) (*proto.DeleteLinkResponse, error) {
	if req.Id == "" {
		return nil, errors.New("id is required")
	}
	if req.CustomerId == "" {
		return nil, errors.New("customer_id is required")
	}

	resp, err := h.linksClient.DeleteLink(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *LinksHandler) UpdateLink(ctx context.Context, req *proto.UpdateLinkRequest) (*proto.UpdateLinkResponse, error) {
	if req.Id == "" {
		return nil, errors.New("id is required")
	}

	if req.OriginalUrl == "" {
		return nil, errors.New("original_url is required")
	}

	resp, err := h.linksClient.UpdateLink(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *LinksHandler) UpdateLinkClicks(ctx context.Context, req *proto.UpdateLinkClicksRequest) (*proto.UpdateLinkClicksResponse, error) {
	if req.Id == "" {
		return nil, errors.New("id is required")
	}

	resp, err := h.linksClient.UpdateLinkClicks(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
