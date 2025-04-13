package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"time"

	"links-service/internal/infra/repository"
	pb "links-service/proto"
	"links-service/utils"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedLinksServiceServer
	repo *repository.Queries
}

func uuidFromString(s string) (pgtype.UUID, error) {
	// Remove hyphens from UUID string
	cleanUUID := strings.ReplaceAll(s, "-", "")
	bytes, err := hex.DecodeString(cleanUUID)
	if err != nil {
		return pgtype.UUID{}, err
	}
	if len(bytes) != 16 {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID length")
	}
	var uuid [16]byte
	copy(uuid[:], bytes)
	return pgtype.UUID{Bytes: uuid, Valid: true}, nil
}

func generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return utils.ConfigInstance.FrontendSource + "/" + string(b)
}

func (s *server) CreateLink(ctx context.Context, req *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	if req.OriginalUrl == "" {
		return nil, errors.New("original_url is required")
	}

	// Validate URL format
	if _, err := url.Parse(req.OriginalUrl); err != nil {
		return nil, errors.New("invalid URL format")
	}

	var shortURL string
	var customSlug pgtype.Text
	var expiresAt pgtype.Timestamptz

	if req.ExpirationDate != nil {
		expirationTime, err := time.Parse(time.RFC3339, *req.ExpirationDate)
		if err != nil {
			return nil, errors.New("invalid expiration date format. Use RFC3339 format (e.g., 2024-12-31T23:59:59Z)")
		}
		if expirationTime.Before(time.Now()) {
			return nil, errors.New("expiration date must be in the future")
		}
		expiresAt = pgtype.Timestamptz{Time: expirationTime, Valid: true}
	}

	if req.CustomSlug != "" {
		exists, err := s.repo.CheckCustomSlugExists(ctx, pgtype.Text{String: req.CustomSlug, Valid: true})
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("custom slug already exists")
		}
		customSlug = pgtype.Text{String: req.CustomSlug, Valid: true}
		shortURL = utils.ConfigInstance.FrontendSource + "/" + req.CustomSlug
	} else {
		for {
			shortURL = generateShortURL()
			exists, err := s.repo.CheckShortURLExists(ctx, shortURL)
			if err != nil {
				return nil, err
			}
			if !exists {
				break
			}
		}
	}

	customerID, err := uuidFromString(req.CustomerId)
	if err != nil {
		return nil, err
	}

	link, err := s.repo.CreateLink(ctx, repository.CreateLinkParams{
		OriginalUrl: req.OriginalUrl,
		ShortUrl:    shortURL,
		CustomSlug:  customSlug,
		CustomerID:  customerID,
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		return nil, err
	}

	response := &pb.CreateLinkResponse{
		Id:         hex.EncodeToString(link.ID.Bytes[:]),
		ShortUrl:   link.ShortUrl,
		CustomSlug: link.CustomSlug.String,
		Clicks:     int32(link.Clicks),
		CreatedAt:  link.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:  link.UpdatedAt.Time.Format(time.RFC3339),
		CustomerId: hex.EncodeToString(link.CustomerID.Bytes[:]),
	}

	if link.ExpiresAt.Valid {
		expirationDate := link.ExpiresAt.Time.Format(time.RFC3339)
		response.ExpirationDate = &expirationDate
	}

	return response, nil
}

func (s *server) GetLink(ctx context.Context, req *pb.GetLinkRequest) (*pb.GetLinkResponse, error) {
	link, err := s.repo.GetLinkByShortURL(ctx, utils.ConfigInstance.FrontendSource+"/"+req.ShortUrl)
	if err != nil {
		return nil, err
	}

	// Check if link is expired
	if link.ExpiresAt.Valid && link.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("link has expired")
	}

	response := &pb.GetLinkResponse{
		Id:          link.ID.String(),
		OriginalUrl: link.OriginalUrl,
		ShortUrl:    link.ShortUrl,
		CustomSlug:  link.CustomSlug.String,
		Clicks:      int32(link.Clicks),
		CreatedAt:   link.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   link.UpdatedAt.Time.Format(time.RFC3339),
	}

	if link.ExpiresAt.Valid {
		expirationDate := link.ExpiresAt.Time.Format(time.RFC3339)
		response.ExpirationDate = &expirationDate
	}

	return response, nil
}

func (s *server) GetCustomerLinks(ctx context.Context, req *pb.GetCustomerLinksRequest) (*pb.GetCustomerLinksResponse, error) {
	if req.CustomerId == "" {
		return nil, status.Error(codes.InvalidArgument, "customer_id is required")
	}

	customerID, err := uuidFromString(req.CustomerId)
	if err != nil {
		return nil, err
	}

	var limit int32
	if req.Limit != nil {
		limit = *req.Limit
	}

	var offset int32
	if req.Offset != nil {
		offset = *req.Offset
	}

	searchTerm := ""
	if req.Search != nil && *req.Search != "" {
		searchTerm = strings.TrimSpace(*req.Search)
	}

	sortBy := getStringValue(req.SortBy)
	sortDirection := getStringValue(req.SortDirection)

	links, err := s.repo.GetLinksByCustomer(ctx, repository.GetLinksByCustomerParams{
		CustomerID: customerID,
		Column2:    limit,
		Column3:    offset,
		Column4:    searchTerm,
		Column5:    getStringValue(req.Status),
		Column6:    getStringValue(req.SlugType),
		Column7:    sortBy,
	})
	if err != nil {
		log.Printf("Error retrieving links: %v", err)
		return nil, status.Error(codes.Internal, "failed to get customer links")
	}

	log.Printf("Retrieved %d links for current page", len(links))

	if sortDirection == "asc" {
		for i, j := 0, len(links)-1; i < j; i, j = i+1, j-1 {
			links[i], links[j] = links[j], links[i]
		}
	}

	response := &pb.GetCustomerLinksResponse{
		Links: make([]*pb.GetLinkResponse, 0, len(links)),
	}

	for _, link := range links {
		response.Links = append(response.Links, &pb.GetLinkResponse{
			Id:          hex.EncodeToString(link.ID.Bytes[:]),
			OriginalUrl: link.OriginalUrl,
			ShortUrl:    link.ShortUrl,
			CustomSlug:  link.CustomSlug.String,
			Clicks:      int32(link.Clicks),
			CreatedAt:   link.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:   link.UpdatedAt.Time.Format(time.RFC3339),
		})

		if link.ExpiresAt.Valid {
			expirationDate := link.ExpiresAt.Time.Format(time.RFC3339)
			response.Links[len(response.Links)-1].ExpirationDate = &expirationDate
		}
	}

	if len(links) > 0 {
		response.Total = int32(links[0].TotalCount)
	}

	return response, nil
}

// Helper function to safely get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *server) DeleteLink(ctx context.Context, req *pb.DeleteLinkRequest) (*pb.DeleteLinkResponse, error) {
	id, err := uuidFromString(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.repo.DeleteLink(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteLinkResponse{
		Success: true,
	}, nil
}

func (s *server) UpdateLink(ctx context.Context, req *pb.UpdateLinkRequest) (*pb.UpdateLinkResponse, error) {
	id, err := uuidFromString(req.Id)
	if err != nil {
		return nil, err
	}

	currentLink, err := s.repo.GetLinkByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var expiresAt pgtype.Timestamptz

	if req.ExpirationDate != nil {
		if *req.ExpirationDate == "" {
			expiresAt = pgtype.Timestamptz{Valid: false}
		} else {
			expirationTime, err := time.Parse(time.RFC3339, *req.ExpirationDate)
			if err != nil {
				return nil, errors.New("invalid expiration date format. Use RFC3339 format (e.g., 2024-12-31T23:59:59Z)")
			}
			if expirationTime.Before(time.Now()) {
				return nil, errors.New("expiration date must be in the future")
			}
			expiresAt = pgtype.Timestamptz{Time: expirationTime, Valid: true}
		}
	} else {
		expiresAt = pgtype.Timestamptz{Valid: false}
	}

	var shortURL string
	if req.CustomSlug != "" && req.CustomSlug != currentLink.CustomSlug.String {
		exists, err := s.repo.CheckCustomSlugExists(ctx, pgtype.Text{String: req.CustomSlug, Valid: true})
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("custom slug already exists")
		}
		shortURL = utils.ConfigInstance.FrontendSource + "/" + req.CustomSlug
	} else {
		shortURL = currentLink.ShortUrl
	}

	link, err := s.repo.UpdateLink(ctx, repository.UpdateLinkParams{
		ID:          id,
		OriginalUrl: req.OriginalUrl,
		ShortUrl:    shortURL,
		CustomSlug:  pgtype.Text{String: req.CustomSlug, Valid: req.CustomSlug != ""},
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		return nil, err
	}

	response := &pb.UpdateLinkResponse{
		Id:          hex.EncodeToString(link.ID.Bytes[:]),
		OriginalUrl: link.OriginalUrl,
		ShortUrl:    link.ShortUrl,
		CustomSlug:  link.CustomSlug.String,
		Clicks:      int32(link.Clicks),
		CreatedAt:   link.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   link.UpdatedAt.Time.Format(time.RFC3339),
		CustomerId:  hex.EncodeToString(link.CustomerID.Bytes[:]),
	}

	if link.ExpiresAt.Valid {
		expirationDate := link.ExpiresAt.Time.Format(time.RFC3339)
		response.ExpirationDate = &expirationDate
	}

	return response, nil
}

func (s *server) UpdateLinkClicks(ctx context.Context, req *pb.UpdateLinkClicksRequest) (*pb.UpdateLinkClicksResponse, error) {
	id, err := uuidFromString(req.Id)
	if err != nil {
		return nil, err
	}

	link, err := s.repo.UpdateLinkClicks(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateLinkClicksResponse{
		Id:          hex.EncodeToString(link.ID.Bytes[:]),
		OriginalUrl: link.OriginalUrl,
		ShortUrl:    link.ShortUrl,
		CustomSlug:  link.CustomSlug.String,
		Clicks:      int32(link.Clicks),
		CreatedAt:   link.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   link.UpdatedAt.Time.Format(time.RFC3339),
		CustomerId:  hex.EncodeToString(link.CustomerID.Bytes[:]),
	}, nil
}

func StartGRPCServer(port string, repo *repository.Queries) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterLinksServiceServer(s, &server{repo: repo})

	log.Printf("gRPC server listening at %v", lis.Addr())
	return s.Serve(lis)
}
