package server

import (
	"context"
	"fmt"
	"links-service-read/internal/infra/repository"
	"links-service-read/internal/logger"
	pb "links-service-read/proto"
	"links-service-read/utils"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedLinksServiceReadServer
	repo *repository.LinksRepository
}

// NewGRPCServer creates a new instance of GRPCServer with the provided LinksRepository.
// It initializes the GRPCServer with the given repository to handle gRPC requests.
//
// Parameters:
//   - repo: A pointer to a LinksRepository instance that provides access to the data layer.
//
// Returns:
//
//	A pointer to a newly created GRPCServer instance.
func NewGRPCServer(repo *repository.LinksRepository) *GRPCServer {
	return &GRPCServer{repo: repo}
}

// GetLink handles the retrieval of a link based on its short URL.
// It validates the input, processes the short URL, and fetches the corresponding link
// from the repository. If the link is found, it checks for expiration and returns
// the link details in the response.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation and deadlines.
//   - req: The request containing the short URL to retrieve the link.
//
// Returns:
//   - *pb.GetLinkResponse: The response containing the link details, including
//     original URL, short URL, custom slug, click count, creation and update timestamps,
//     and expiration date (if applicable).
//   - error: An error if the request is invalid, the link is not found, or an internal
//     error occurs during processing.
//
// Possible Errors:
//   - codes.InvalidArgument: Returned if the short URL is missing in the request.
//   - codes.NotFound: Returned if the link corresponding to the short URL is not found.
//   - codes.FailedPrecondition: Returned if the link has expired.
//   - codes.Internal: Returned if an internal error occurs while fetching the link.
//
// Notes:
//   - If the short URL includes the domain, it is stripped before processing.
//   - Expiration is checked against the current time, and an error is returned if the link has expired.
func (s *GRPCServer) GetLink(ctx context.Context, req *pb.GetLinkRequest) (*pb.GetLinkResponse, error) {
	if req.ShortUrl == "" {
		logger.Log.Error("short_url is required")
		return nil, status.Error(codes.InvalidArgument, "short_url is required")
	}

	shortURL := req.ShortUrl
	baseURL := utils.ConfigInstance.FrontendSource
	if strings.HasPrefix(shortURL, baseURL+"/") {
		shortURL = strings.TrimPrefix(shortURL, baseURL+"/")
	}

	link, err := s.repo.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			logger.Log.Error("link not found")
			return nil, status.Error(codes.NotFound, "link not found")
		}

		logger.Log.Error("failed to get link", zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get link: %v", err))
	}

	if link.ExpirationDate != nil && *link.ExpirationDate != "" {
		expirationTime, err := time.Parse(time.RFC3339, *link.ExpirationDate)
		if err == nil && expirationTime.Before(time.Now()) {
			logger.Log.Error("link has expired", zap.String("expiration_date", *link.ExpirationDate))
			return nil, status.Error(codes.FailedPrecondition, "link has expired")
		}
	}

	logger.Log.Info("link retrieved successfully", zap.String("short_url", shortURL))

	return &pb.GetLinkResponse{
		Id:             link.ID,
		OriginalUrl:    link.OriginalURL,
		ShortUrl:       baseURL + "/" + link.ShortURL,
		CustomSlug:     link.CustomSlug,
		Clicks:         link.Clicks,
		CreatedAt:      link.CreatedAt,
		UpdatedAt:      link.UpdatedAt,
		ExpirationDate: link.ExpirationDate,
	}, nil
}

// GetCustomerLinks retrieves a list of links associated with a specific customer ID.
// It validates the input request to ensure the customer ID is provided, fetches the links
// from the repository, and constructs a response containing the link details.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellation signals.
//   - req: A pointer to a GetCustomerLinksRequest containing the customer ID.
//
// Returns:
//   - A pointer to a GetCustomerLinksResponse containing the list of links associated with the customer.
//   - An error if the customer ID is missing, or if there is an issue retrieving the links from the repository.
//
// Errors:
//   - codes.InvalidArgument: Returned if the customer ID is not provided in the request.
//   - codes.Internal: Returned if there is an internal error while fetching the links.
//
// The response includes details such as the link ID, original URL, short URL, custom slug,
// click count, creation and update timestamps, and expiration date.
func (s *GRPCServer) GetCustomerLinks(ctx context.Context, req *pb.GetCustomerLinksRequest) (*pb.GetCustomerLinksResponse, error) {
	if req.CustomerId == "" {
		logger.Log.Error("customer_id is required")
		return nil, status.Error(codes.InvalidArgument, "customer_id is required")
	}

	links, err := s.repo.GetCustomerLinks(ctx, req.CustomerId)
	if err != nil {
		logger.Log.Error("failed to get customer links", zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get customer links: %v", err))
	}

	baseURL := utils.ConfigInstance.FrontendSource
	response := &pb.GetCustomerLinksResponse{
		Links: make([]*pb.GetLinkResponse, 0, len(links)),
	}

	for _, link := range links {
		linkResponse := &pb.GetLinkResponse{
			Id:             link.ID,
			OriginalUrl:    link.OriginalURL,
			ShortUrl:       baseURL + "/" + link.ShortURL,
			CustomSlug:     link.CustomSlug,
			Clicks:         link.Clicks,
			CreatedAt:      link.CreatedAt,
			UpdatedAt:      link.UpdatedAt,
			ExpirationDate: link.ExpirationDate,
		}

		response.Links = append(response.Links, linkResponse)
	}

	logger.Log.Info("customer links retrieved successfully", zap.String("customer_id", req.CustomerId))
	return response, nil
}

// StartGRPCServer starts a gRPC server on the specified port and registers the LinksServiceReadServer.
// It also enables server reflection for tools like grpcurl.
//
// Parameters:
//   - port: The port on which the gRPC server will listen.
//   - repo: A pointer to the LinksRepository, which provides the necessary data access layer.
//
// Returns:
//   - error: An error if the server fails to start or listen on the specified port.
//
// Example usage:
//
//	err := StartGRPCServer("50051", repo)
//	if err != nil {
//	    log.Fatalf("Failed to start gRPC server: %v", err)
//	}
func StartGRPCServer(port string, repo *repository.LinksRepository) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterLinksServiceReadServer(server, NewGRPCServer(repo))

	// Habilitar reflection para ferramentas como grpcurl
	reflection.Register(server)

	logger.Log.Info("gRPC server listening", zap.String("port", port))
	return server.Serve(lis)
}
