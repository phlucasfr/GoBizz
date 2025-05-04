package server

import (
	"context"
	"fmt"
	"links-service-write/internal/infra/repository"
	"links-service-write/internal/logger"
	pb "links-service-write/proto"
	"links-service-write/utils"
	"net"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedLinksServiceWriteServer
	repo *repository.LinksRepository
}

// NewGRPCServer creates a new instance of GRPCServer with the provided LinksRepository.
// It initializes the server with the given repository to handle gRPC requests.
//
// Parameters:
//   - repo: A pointer to a LinksRepository instance that provides access to the data layer.
//
// Returns:
//
//	A pointer to a GRPCServer instance configured with the provided repository.
func NewGRPCServer(repo *repository.LinksRepository) *GRPCServer {
	return &GRPCServer{repo: repo}
}

// CreateLink handles the creation of a new shortened link.
// It validates the input request, generates a unique ID and slug, and stores the link in the repository.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation and deadlines.
//   - req: A pointer to a CreateLinkRequest containing the details of the link to be created.
//
// Returns:
//   - A pointer to a CreateLinkResponse containing the details of the created link, including the short URL.
//   - An error if the creation process fails due to invalid input, repository issues, or other internal errors.
//
// Validation:
//   - Ensures the OriginalUrl field is not empty.
//   - Validates the format of the OriginalUrl.
//   - If an ExpirationDate is provided, ensures it is in RFC3339 format and is a future date.
//   - If a CustomSlug is provided, checks for its uniqueness in the repository.
//
// Behavior:
//   - Generates a unique ID for the link.
//   - If no CustomSlug is provided, generates a random slug and ensures its uniqueness.
//   - Creates a new link record in the repository with the provided and generated details.
//   - Constructs the short URL using the base frontend source URL.
//
// Possible Errors:
//   - InvalidArgument: If required fields are missing or invalid (e.g., empty OriginalUrl, invalid URL format).
//   - AlreadyExists: If the provided CustomSlug or generated slug already exists.
//   - Internal: If there are issues generating the ID/slug or interacting with the repository.
func (s *GRPCServer) CreateLink(ctx context.Context, req *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	if req.OriginalUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "original_url is required")
	}

	if _, err := url.Parse(req.OriginalUrl); err != nil {
		logger.Log.Error("invalid URL format", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid URL format")
	}

	var expirationDate *string
	if req.ExpirationDate != nil && *req.ExpirationDate != "" {
		expirationTime, err := time.Parse(time.RFC3339, *req.ExpirationDate)
		if err != nil {
			logger.Log.Error("invalid expiration date format", zap.Error(err))
			return nil, status.Error(codes.InvalidArgument,
				"invalid expiration date format. Use RFC3339 format (e.g., 2024-12-31T23:59:59Z)")
		}
		if expirationTime.Before(time.Now()) {
			logger.Log.Error("expiration date must be in the future", zap.Error(err))
			return nil, status.Error(codes.InvalidArgument, "expiration date must be in the future")
		}
		expirationDate = req.ExpirationDate
	}

	id, err := utils.GenerateRandomSlug(10)
	if err != nil {
		logger.Log.Error("failed to generate ID", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to generate ID")
	}

	var shortSlug string
	var customSlug string = req.CustomSlug

	if customSlug != "" {
		existingLink, err := s.repo.GetLinkByCustomSlug(ctx, customSlug)
		if err == nil && existingLink != nil {
			logger.Log.Error("custom slug already exists", zap.String("custom_slug", customSlug))
			return nil, status.Error(codes.AlreadyExists, "custom slug already exists")
		} else if err != nil && !strings.Contains(err.Error(), "not found") {
			logger.Log.Error("error checking custom slug", zap.Error(err))
			return nil, status.Error(codes.Internal, fmt.Sprintf("error checking custom slug: %v", err))
		}

		shortSlug = customSlug
	} else {
		for {
			generatedSlug, err := utils.GenerateRandomSlug(6)
			if err != nil {
				logger.Log.Error("failed to generate slug", zap.Error(err))
				return nil, status.Error(codes.Internal, "failed to generate slug")
			}
			existingLink, err := s.repo.GetLinkByShortURL(ctx, generatedSlug)
			if err != nil && strings.Contains(err.Error(), "not found") {
				shortSlug = generatedSlug
				break
			} else if err == nil && existingLink != nil && existingLink.CustomSlug != customSlug {
				logger.Log.Error("generated slug already exists", zap.String("generated_slug", generatedSlug))
				return nil, status.Error(codes.AlreadyExists, "generated slug already exists")
			} else if err == nil && existingLink != nil {
				logger.Log.Error("generated slug already exists", zap.String("generated_slug", generatedSlug))
				return nil, status.Error(codes.AlreadyExists, "generated slug already exists")
			} else if err != nil {
				logger.Log.Error("error checking slug", zap.Error(err))
				return nil, status.Error(codes.Internal, fmt.Sprintf("error checking slug: %v", err))
			}
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)

	link := repository.Link{
		ID:             id,
		ShortURL:       shortSlug,
		OriginalURL:    req.OriginalUrl,
		CustomSlug:     customSlug,
		CustomerID:     req.CustomerId,
		Clicks:         0,
		CreatedAt:      now,
		UpdatedAt:      now,
		ExpirationDate: expirationDate,
	}

	createdLink, err := s.repo.CreateLink(ctx, link)
	if err != nil {
		logger.Log.Error("failed to create link", zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create link: %v", err))
	}

	baseURL := utils.ConfigInstance.FrontendSource

	logger.Log.Info("link created successfully", zap.String("short_url", createdLink.ShortURL))
	return &pb.CreateLinkResponse{
		Id:             createdLink.ID,
		ShortUrl:       baseURL + "/" + createdLink.ShortURL,
		CustomSlug:     createdLink.CustomSlug,
		Clicks:         createdLink.Clicks,
		CreatedAt:      createdLink.CreatedAt,
		UpdatedAt:      createdLink.UpdatedAt,
		CustomerId:     createdLink.CustomerID,
		ExpirationDate: createdLink.ExpirationDate,
	}, nil
}

// DeleteLink handles the deletion of a link based on the provided request.
// It validates the input parameters and interacts with the repository to perform the deletion.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation and deadlines.
//   - req: A pointer to a DeleteLinkRequest containing the ID of the link to delete
//     and the customer ID associated with the link.
//
// Returns:
//   - A pointer to a DeleteLinkResponse indicating the success of the operation.
//   - An error if the operation fails, which could be one of the following:
//   - codes.InvalidArgument: If the link ID or customer ID is missing.
//   - codes.NotFound: If the link does not exist.
//   - codes.PermissionDenied: If the link does not belong to the specified customer.
//   - codes.Internal: If an internal error occurs during the deletion process.
func (s *GRPCServer) DeleteLink(ctx context.Context, req *pb.DeleteLinkRequest) (*pb.DeleteLinkResponse, error) {
	if req.Id == "" {
		logger.Log.Error("link ID is required")
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.CustomerId == "" {
		logger.Log.Error("customer ID is required")
		return nil, status.Error(codes.InvalidArgument, "customer_id is required")
	}

	err := s.repo.DeleteLink(ctx, req.Id, req.CustomerId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			logger.Log.Error("link not found", zap.String("link_id", req.Id))
			return nil, status.Error(codes.NotFound, "link not found")
		}
		if strings.Contains(err.Error(), "does not belong") {
			logger.Log.Error("link does not belong to this customer", zap.String("customer_id", req.CustomerId))
			return nil, status.Error(codes.PermissionDenied, "link does not belong to this customer")
		}
		logger.Log.Error("failed to delete link", zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete link: %v", err))
	}

	logger.Log.Info("link deleted successfully", zap.String("link_id", req.Id))
	return &pb.DeleteLinkResponse{
		Success: true,
	}, nil
}

// UpdateLink handles the update of an existing link in the system.
// It validates the input request, checks for ownership, and ensures
// that the provided data adheres to the required constraints before
// updating the link in the repository.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - req: A pointer to pb.UpdateLinkRequest containing the details of the link to be updated.
//
// Returns:
//   - A pointer to pb.UpdateLinkResponse containing the updated link details.
//   - An error if the update operation fails or the input is invalid.
//
// Validation:
//   - The `id` field in the request must not be empty.
//   - The `original_url` field must not be empty and must be a valid URL format.
//   - If `custom_slug` is provided, it must not conflict with an existing slug.
//   - If `expiration_date` is provided, it must be in RFC3339 format and set to a future date.
//   - The `customer_id` field must not be empty and cannot be changed from the original value.
//
// Errors:
//   - codes.InvalidArgument: If required fields are missing or invalid.
//   - codes.NotFound: If the link with the specified ID does not exist.
//   - codes.AlreadyExists: If the custom slug is already in use by another link.
//   - codes.PermissionDenied: If the `customer_id` is modified.
//   - codes.Internal: If there is an internal error during the update process.
func (s *GRPCServer) UpdateLink(ctx context.Context, req *pb.UpdateLinkRequest) (*pb.UpdateLinkResponse, error) {
	if req.Id == "" {
		logger.Log.Error("link ID is required")
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.OriginalUrl == "" {
		logger.Log.Error("original_url is required")
		return nil, status.Error(codes.InvalidArgument, "original_url is required")
	}

	if _, err := url.Parse(req.OriginalUrl); err != nil {
		logger.Log.Error("invalid URL format", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid URL format")
	}

	existingLink, err := s.repo.GetLinkByID(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			logger.Log.Error("link not found", zap.String("link_id", req.Id))
			return nil, status.Error(codes.NotFound, "link not found")
		}
		logger.Log.Error("failed to get link", zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get link: %v", err))
	}

	if req.CustomSlug != "" && req.CustomSlug != existingLink.CustomSlug {
		existingSlugLink, err := s.repo.GetLinkByCustomSlug(ctx, req.CustomSlug)
		if err == nil && existingSlugLink != nil && existingSlugLink.ID != req.Id {
			logger.Log.Error("custom slug already in use", zap.String("custom_slug", req.CustomSlug))
			return nil, status.Error(codes.AlreadyExists, "custom slug already in use")
		}
	}

	var expirationDate *string
	if req.ExpirationDate != nil {
		if *req.ExpirationDate != "" {
			expirationTime, err := time.Parse(time.RFC3339, *req.ExpirationDate)
			if err != nil {
				logger.Log.Error("invalid expiration date format", zap.Error(err))
				return nil, status.Error(codes.InvalidArgument,
					"invalid expiration date format. Use RFC3339 format (e.g., 2024-12-31T23:59:59Z)")
			}
			if expirationTime.Before(time.Now()) {
				logger.Log.Error("expiration date must be in the future", zap.Error(err))
				return nil, status.Error(codes.InvalidArgument, "expiration date must be in the future")
			}
		}
		expirationDate = req.ExpirationDate
	}

	if req.CustomerId == "" {
		logger.Log.Error("customer_id is required")
		return nil, status.Error(codes.InvalidArgument, "customer_id is required")
	}

	if req.CustomerId != "" && req.CustomerId != existingLink.CustomerID {
		logger.Log.Error("customer_id cannot be changed", zap.String("customer_id", req.CustomerId))
		return nil, status.Error(codes.PermissionDenied, "customer_id cannot be changed")
	}

	updatedLink := repository.Link{
		ID:             req.Id,
		ShortURL:       existingLink.ShortURL,
		OriginalURL:    req.OriginalUrl,
		CustomSlug:     req.CustomSlug,
		CustomerID:     req.CustomerId,
		Clicks:         existingLink.Clicks,
		CreatedAt:      existingLink.CreatedAt,
		ExpirationDate: expirationDate,
	}

	result, err := s.repo.UpdateLink(ctx, updatedLink)
	if err != nil {
		logger.Log.Error("failed to update link", zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update link: %v", err))
	}

	logger.Log.Info("link updated successfully", zap.String("link_id", result.ID))
	baseURL := utils.ConfigInstance.FrontendSource
	return &pb.UpdateLinkResponse{
		Id:             result.ID,
		OriginalUrl:    result.OriginalURL,
		ShortUrl:       baseURL + "/" + result.ShortURL,
		CustomSlug:     result.CustomSlug,
		Clicks:         result.Clicks,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
		CustomerId:     result.CustomerID,
		ExpirationDate: result.ExpirationDate,
	}, nil
}

// UpdateLinkClicks updates the click count for a specific link identified by its ID.
// It validates the input request to ensure the ID is provided, and interacts with the repository
// to update the click count. If the link is not found, it returns a NotFound error. For other
// errors, it returns an Internal error with details.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - req: A pointer to pb.UpdateLinkClicksRequest containing the ID of the link to update.
//
// Returns:
//   - A pointer to pb.UpdateLinkClicksResponse containing the updated link details, including
//     the ID, original URL, short URL, custom slug, click count, creation and update timestamps,
//     customer ID, and expiration date.
//   - An error if the operation fails, with appropriate gRPC status codes.
func (s *GRPCServer) UpdateLinkClicks(ctx context.Context, req *pb.UpdateLinkClicksRequest) (*pb.UpdateLinkClicksResponse, error) {
	if req.Id == "" {
		logger.Log.Error("link ID is required")
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	updatedLink, err := s.repo.UpdateLinkClicks(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			logger.Log.Error("link not found", zap.String("link_id", req.Id))
			return nil, status.Error(codes.NotFound, "link not found")
		}
		logger.Log.Error("failed to update link clicks", zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update link clicks: %v", err))
	}

	logger.Log.Info("link clicks updated successfully", zap.String("link_id", updatedLink.ID))
	baseURL := utils.ConfigInstance.FrontendSource
	return &pb.UpdateLinkClicksResponse{
		Id:             updatedLink.ID,
		OriginalUrl:    updatedLink.OriginalURL,
		ShortUrl:       baseURL + "/" + updatedLink.ShortURL,
		CustomSlug:     updatedLink.CustomSlug,
		Clicks:         updatedLink.Clicks,
		CreatedAt:      updatedLink.CreatedAt,
		UpdatedAt:      updatedLink.UpdatedAt,
		CustomerId:     updatedLink.CustomerID,
		ExpirationDate: updatedLink.ExpirationDate,
	}, nil
}

// StartGRPCServer starts a gRPC server on the specified port and registers the LinksServiceWriteServer.
// It also enables server reflection for tools like grpcurl.
//
// Parameters:
//   - port: The port on which the gRPC server will listen.
//   - repo: A pointer to the LinksRepository, which provides the necessary data operations.
//
// Returns:
//   - error: An error if the server fails to start or encounters an issue.
//
// This function sets up a TCP listener, initializes a gRPC server, registers the LinksServiceWriteServer
// implementation, and enables reflection for debugging and testing purposes.
func StartGRPCServer(port string, repo *repository.LinksRepository) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Log.Error("failed to listen", zap.Error(err))
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterLinksServiceWriteServer(server, NewGRPCServer(repo))

	// Habilitar reflection para ferramentas como grpcurl
	reflection.Register(server)

	logger.Log.Info("gRPC server listening on port", zap.String("port", port))
	return server.Serve(lis)
}
