package repository

import (
	"context"
	"fmt"
	"links-service-read/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type Link struct {
	ID             string  `dynamodbav:"id"`
	ShortURL       string  `dynamodbav:"short_url"`
	OriginalURL    string  `dynamodbav:"original_url"`
	CustomSlug     string  `dynamodbav:"custom_slug"`
	CustomerID     string  `dynamodbav:"customer_id"`
	Clicks         int32   `dynamodbav:"clicks"`
	CreatedAt      string  `dynamodbav:"created_at"`
	UpdatedAt      string  `dynamodbav:"updated_at"`
	ExpirationDate *string `dynamodbav:"expiration_date,omitempty"`
	TTL            *int64  `dynamodbav:"ttl,omitempty"`
}

type LinksRepository struct {
	db *dynamodb.Client
}

// NewLinksRepository creates a new instance of LinksRepository with the provided DynamoDB client.
// It initializes the repository to interact with the DynamoDB database.
//
// Parameters:
//   - db: A pointer to a dynamodb.Client instance used to perform database operations.
//
// Returns:
//   - A pointer to a LinksRepository instance.
func NewLinksRepository(db *dynamodb.Client) *LinksRepository {
	return &LinksRepository{db: db}
}

// GetLinkByShortURL retrieves a link from the DynamoDB table "Links" using the provided short URL.
// It queries the database for an item with the specified short URL as the key.
// If the item is found, it unmarshals the result into a Link struct and returns it.
// If the item is not found or an error occurs during the process, it returns an appropriate error.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - shortURL: The short URL used as the key to query the database.
//
// Returns:
//   - *Link: A pointer to the Link struct containing the retrieved link data.
//   - error: An error if the link is not found or if any other issue occurs during the operation.
func (r *LinksRepository) GetLinkByShortURL(ctx context.Context, shortURL string) (*Link, error) {
	result, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Links"),
		Key: map[string]types.AttributeValue{
			"short_url": &types.AttributeValueMemberS{Value: shortURL},
		},
	})
	if err != nil {
		logger.Log.Error("Failed to get link", zap.Error(err))
		return nil, fmt.Errorf("failed to get link: %v", err)
	}

	if len(result.Item) == 0 {
		logger.Log.Info("Link not found", zap.String("short_url", shortURL))
		return nil, fmt.Errorf("link not found")
	}

	var link Link
	if err := attributevalue.UnmarshalMap(result.Item, &link); err != nil {
		logger.Log.Error("Failed to unmarshal link", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal link: %v", err)
	}

	logger.Log.Info("Link retrieved successfully", zap.String("short_url", shortURL))
	return &link, nil
}

// GetLinkByID retrieves a link from the "Links" table in DynamoDB by its ID.
// It uses a filter expression to search for the link with the specified ID.
// If the link is found, it unmarshals the result into a Link struct and returns it.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - id: The ID of the link to retrieve.
//
// Returns:
//   - A pointer to the Link struct if the link is found.
//   - An error if the link is not found, the expression fails to build,
//     the scan operation fails, or unmarshaling the result fails.
func (r *LinksRepository) GetLinkByID(ctx context.Context, id string) (*Link, error) {
	expr, err := expression.NewBuilder().
		WithFilter(expression.Name("id").Equal(expression.Value(id))).
		Build()
	if err != nil {
		logger.Log.Error("Failed to build expression", zap.Error(err))
		return nil, fmt.Errorf("failed to build expression: %v", err)
	}

	result, err := r.db.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String("Links"),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		logger.Log.Error("Failed to scan for link", zap.Error(err))
		return nil, fmt.Errorf("failed to scan for link: %v", err)
	}

	if len(result.Items) == 0 {
		logger.Log.Info("Link not found", zap.String("id", id))
		return nil, fmt.Errorf("link not found")
	}

	var link Link
	if err := attributevalue.UnmarshalMap(result.Items[0], &link); err != nil {
		logger.Log.Error("Failed to unmarshal link", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal link: %v", err)
	}

	logger.Log.Info("Link retrieved successfully", zap.String("id", id))
	return &link, nil
}

// GetLinkByCustomSlug retrieves a link from the DynamoDB table "Links" using the provided custom slug.
// It queries the "ByCustomSlug" index to find the link that matches the given custom slug.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - customSlug: The custom slug used to query the link.
//
// Returns:
//   - A pointer to the Link object if found.
//   - An error if the query fails, the link is not found, or unmarshalling the result fails.
func (r *LinksRepository) GetLinkByCustomSlug(ctx context.Context, customSlug string) (*Link, error) {
	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("Links"),
		IndexName:              aws.String("ByCustomSlug"),
		KeyConditionExpression: aws.String("custom_slug = :slug"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":slug": &types.AttributeValueMemberS{Value: customSlug},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		logger.Log.Error("Failed to query link by custom slug", zap.Error(err))
		return nil, fmt.Errorf("failed to query link by custom slug: %v", err)
	}

	if len(result.Items) == 0 {
		logger.Log.Info("Link not found", zap.String("custom_slug", customSlug))
		return nil, fmt.Errorf("link not found")
	}

	var link Link
	if err := attributevalue.UnmarshalMap(result.Items[0], &link); err != nil {
		logger.Log.Error("Failed to unmarshal link", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal link: %v", err)
	}

	logger.Log.Info("Link retrieved successfully", zap.String("custom_slug", customSlug))
	return &link, nil
}

// GetCustomerLinks retrieves a list of links associated with a specific customer ID
// from the DynamoDB table. The links are fetched in descending order based on their
// creation time.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - customerID: The unique identifier of the customer whose links are to be retrieved.
//
// Returns:
//   - A slice of pointers to Link objects representing the customer's links.
//   - An error if the query or unmarshalling process fails.
func (r *LinksRepository) GetCustomerLinks(ctx context.Context, customerID string) ([]*Link, error) {
	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("Links"),
		IndexName:              aws.String("ByCustomer"),
		KeyConditionExpression: aws.String("customer_id = :customer"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":customer": &types.AttributeValueMemberS{Value: customerID},
		},
		ScanIndexForward: aws.Bool(false),
	})
	if err != nil {
		logger.Log.Error("Failed to query links by customer", zap.Error(err))
		return nil, fmt.Errorf("failed to query links by customer: %v", err)
	}

	links := make([]*Link, 0, len(result.Items))
	for _, item := range result.Items {
		var link Link
		if err := attributevalue.UnmarshalMap(item, &link); err != nil {
			logger.Log.Error("Failed to unmarshal link", zap.Error(err))
			return nil, fmt.Errorf("failed to unmarshal link: %v", err)
		}
		links = append(links, &link)
	}

	logger.Log.Info("Links retrieved successfully", zap.String("customer_id", customerID))
	return links, nil
}
