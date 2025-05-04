package repository

import (
	"context"
	"errors"
	"fmt"
	"links-service-write/internal/logger"
	"time"

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

// CreateLink inserts a new link into the DynamoDB table "Links".
// It ensures that the custom slug is unique and sets default timestamps if not provided.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - link: The Link object containing the details of the link to be created.
//
// Returns:
//   - A pointer to the created Link object.
//   - An error if the operation fails, such as invalid expiration date format,
//     failure to marshal the link, or if the custom slug already exists.
//
// Behavior:
//   - If the CreatedAt or UpdatedAt fields in the Link object are empty,
//     they are set to the current UTC time in RFC3339 format.
//   - If the ExpirationDate field is provided, it is parsed and used to set the TTL (Time-To-Live) value.
//   - The function ensures that the custom slug is unique by using a conditional expression
//     in the DynamoDB PutItem operation.
//
// Errors:
//   - Returns an error if the ExpirationDate is in an invalid format.
//   - Returns an error if the custom slug already exists in the table.
//   - Returns an error if there is a failure in marshaling the Link object or inserting it into DynamoDB.
func (r *LinksRepository) CreateLink(ctx context.Context, link Link) (*Link, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	if link.CreatedAt == "" {
		link.CreatedAt = now
	}
	if link.UpdatedAt == "" {
		link.UpdatedAt = now
	}

	if link.ExpirationDate != nil && *link.ExpirationDate != "" {
		expTime, err := time.Parse(time.RFC3339, *link.ExpirationDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expiration date format: %v", err)
		}
		ttl := expTime.Unix()
		link.TTL = &ttl
	}

	item, err := attributevalue.MarshalMap(link)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal link: %v", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String("Links"),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(custom_slug) OR custom_slug = :empty"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":empty": &types.AttributeValueMemberS{Value: ""},
		},
	})

	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			logger.Log.Error("custom slug already exists", zap.String("custom_slug", link.CustomSlug))
			return nil, fmt.Errorf("custom slug '%s' already exists", link.CustomSlug)
		}
		logger.Log.Error("failed to create link", zap.Error(err))
		return nil, fmt.Errorf("failed to create link: %v", err)
	}

	logger.Log.Info("link created successfully", zap.String("short_url", link.ShortURL))
	return &link, nil
}

// GetLinkByShortURL retrieves a link from the database using its short URL.
// It queries the "Links" table in DynamoDB and unmarshals the result into a Link struct.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - shortURL: The short URL of the link to retrieve.
//
// Returns:
//   - A pointer to the Link struct if the link is found.
//   - An error if the link is not found, the query fails, or unmarshaling fails.
func (r *LinksRepository) GetLinkByShortURL(ctx context.Context, shortURL string) (*Link, error) {
	result, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Links"),
		Key: map[string]types.AttributeValue{
			"short_url": &types.AttributeValueMemberS{Value: shortURL},
		},
	})
	if err != nil {
		logger.Log.Error("failed to get link", zap.Error(err))
		return nil, fmt.Errorf("failed to get link: %v", err)
	}

	if len(result.Item) == 0 {
		logger.Log.Error("link not found", zap.String("short_url", shortURL))
		return nil, fmt.Errorf("link not found")
	}

	var link Link
	if err := attributevalue.UnmarshalMap(result.Item, &link); err != nil {
		logger.Log.Error("failed to unmarshal link", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal link: %v", err)
	}

	logger.Log.Info("link retrieved successfully", zap.String("short_url", link.ShortURL))
	return &link, nil
}

// GetLinkByID retrieves a link from the "Links" table by its ID using the "ByID" global secondary index (GSI).
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - id: The unique identifier of the link to retrieve.
//
// Returns:
//   - A pointer to the Link struct if the link is found.
//   - An error if the link is not found, if there is an issue building the query expression,
//     querying the database, or unmarshaling the result.
func (r *LinksRepository) GetLinkByID(ctx context.Context, id string) (*Link, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(
			expression.Key("id").Equal(expression.Value(id)),
		).
		Build()
	if err != nil {
		logger.Log.Error("failed to build expression", zap.Error(err))
		return nil, fmt.Errorf("building expression: %w", err)
	}

	out, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String("Links"),
		IndexName:                 aws.String("ByID"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		logger.Log.Error("failed to query GSI ByID", zap.Error(err))
		return nil, fmt.Errorf("querying GSI ByID: %w", err)
	}
	if len(out.Items) == 0 {
		logger.Log.Error("link not found", zap.String("id", id))
		return nil, fmt.Errorf("link not found")
	}

	var link Link
	if err := attributevalue.UnmarshalMap(out.Items[0], &link); err != nil {
		logger.Log.Error("failed to unmarshal link", zap.Error(err))
		return nil, fmt.Errorf("unmarshal link: %w", err)
	}

	logger.Log.Info("link retrieved successfully", zap.String("id", id))
	return &link, nil
}

// GetLinkByCustomSlug retrieves a link from the DynamoDB table "Links" using the provided custom slug.
// It queries the "ByCustomSlug" index to find the link associated with the given custom slug.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - customSlug: The custom slug used to identify the link.
//
// Returns:
//   - A pointer to the Link struct if a matching link is found.
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
		logger.Log.Error("failed to query link by custom slug", zap.Error(err))
		return nil, fmt.Errorf("failed to query link by custom slug: %v", err)
	}

	if len(result.Items) == 0 {
		logger.Log.Warn("link not found", zap.String("custom_slug", customSlug))
		return nil, fmt.Errorf("link not found")
	}

	var link Link
	if err := attributevalue.UnmarshalMap(result.Items[0], &link); err != nil {
		logger.Log.Error("failed to unmarshal link", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal link: %v", err)
	}

	logger.Log.Info("link retrieved successfully", zap.String("custom_slug", customSlug))
	return &link, nil
}

// GetCustomerLinks retrieves a list of links associated with a specific customer ID
// from the DynamoDB table. The links are fetched in descending order based on the
// "created_at" attribute.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - customerID: The ID of the customer whose links are to be retrieved.
//
// Returns:
//   - A slice of pointers to Link objects representing the customer's links.
//   - An error if the query fails or if unmarshalling the data encounters an issue.
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
		logger.Log.Error("failed to query links by customer", zap.Error(err))
		return nil, fmt.Errorf("failed to query links by customer: %v", err)
	}

	links := make([]*Link, 0, len(result.Items))
	for _, item := range result.Items {
		var link Link
		if err := attributevalue.UnmarshalMap(item, &link); err != nil {
			logger.Log.Error("failed to unmarshal link", zap.Error(err))
			return nil, fmt.Errorf("failed to unmarshal link: %v", err)
		}
		links = append(links, &link)
	}

	logger.Log.Info("customer links retrieved successfully", zap.String("customer_id", customerID))
	return links, nil
}

// DeleteLink deletes a link from the database based on its ID and associated customer ID.
// It first retrieves the link using the provided ID to obtain the primary key (short_url),
// and then deletes the link using this key.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - id: The unique identifier of the link to be deleted.
//   - customerID: The ID of the customer associated with the link.
//
// Returns:
//   - error: An error if the operation fails, or nil if the deletion is successful.
func (r *LinksRepository) DeleteLink(ctx context.Context, id, customerID string) error {
	link, err := r.GetLinkByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("Links"),
		Key: map[string]types.AttributeValue{
			"short_url": &types.AttributeValueMemberS{Value: link.ShortURL},
		},
	})
	if err != nil {
		logger.Log.Error("failed to delete link", zap.Error(err))
		return fmt.Errorf("failed to delete link: %v", err)
	}

	logger.Log.Info("link deleted successfully", zap.String("short_url", link.ShortURL))
	return nil
}

// UpdateLink updates an existing link in the repository.
// It retrieves the original link to ensure it belongs to the client and uses its primary key.
// The function preserves certain fields from the original link, such as ShortURL, CreatedAt, and Clicks,
// while updating the UpdatedAt field to the current time.
//
// If the CustomerID field is empty, an error is returned.
//
// If the ExpirationDate field is provided, it validates the format and calculates the TTL (time-to-live).
// If the ExpirationDate is invalid, an error is returned. If no ExpirationDate is provided, the TTL is set to nil.
//
// The updated link is marshaled into a DynamoDB-compatible attribute map and stored in the "Links" table.
// If the update operation fails, an error is returned.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - link: The Link object containing the updated data.
//
// Returns:
//   - A pointer to the updated Link object.
//   - An error if the update operation fails or if validation errors occur.
func (r *LinksRepository) UpdateLink(ctx context.Context, link Link) (*Link, error) {
	existingLink, err := r.GetLinkByID(ctx, link.ID)
	if err != nil {
		return nil, err
	}

	link.ShortURL = existingLink.ShortURL
	link.CreatedAt = existingLink.CreatedAt
	link.Clicks = existingLink.Clicks
	link.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if link.CustomerID == "" {
		logger.Log.Error("customer_id cannot be empty")
		return nil, fmt.Errorf("customer_id cannot be empty")
	}

	if link.ExpirationDate != nil && *link.ExpirationDate != "" {
		expTime, err := time.Parse(time.RFC3339, *link.ExpirationDate)
		if err != nil {
			logger.Log.Error("invalid expiration date format", zap.Error(err))
			return nil, fmt.Errorf("invalid expiration date format: %v", err)
		}
		ttl := expTime.Unix()
		link.TTL = &ttl
	} else {
		link.TTL = nil
	}

	item, err := attributevalue.MarshalMap(link)
	if err != nil {
		logger.Log.Error("failed to marshal updated link", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal updated link: %v", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Links"),
		Item:      item,
	})
	if err != nil {
		logger.Log.Error("failed to update link", zap.Error(err))
		return nil, fmt.Errorf("failed to update link: %v", err)
	}

	logger.Log.Info("link updated successfully", zap.String("short_url", link.ShortURL))
	return &link, nil
}

// UpdateLinkClicks increments the click count of a link and updates its
// "updated_at" timestamp in the database. It retrieves the link by its ID,
// constructs an update expression to modify the "clicks" and "updated_at" fields,
// and applies the update to the DynamoDB table.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - id: The unique identifier of the link to be updated.
//
// Returns:
//   - A pointer to the updated Link object if the operation is successful.
//   - An error if the link cannot be retrieved, the update expression cannot be
//     built, the update operation fails, or the updated link cannot be unmarshaled.
func (r *LinksRepository) UpdateLinkClicks(ctx context.Context, id string) (*Link, error) {
	link, err := r.GetLinkByID(ctx, id)
	if err != nil {
		return nil, err
	}

	expr, err := expression.NewBuilder().
		WithUpdate(
			expression.Set(
				expression.Name("clicks"),
				expression.Plus(expression.Name("clicks"), expression.Value(1)),
			).Set(
				expression.Name("updated_at"),
				expression.Value(time.Now().UTC().Format(time.RFC3339)),
			),
		).Build()
	if err != nil {
		logger.Log.Error("failed to build update expression", zap.Error(err))
		return nil, fmt.Errorf("failed to build update expression: %v", err)
	}

	result, err := r.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("Links"),
		Key: map[string]types.AttributeValue{
			"short_url": &types.AttributeValueMemberS{Value: link.ShortURL},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		logger.Log.Error("failed to update link clicks", zap.Error(err))
		return nil, fmt.Errorf("failed to update link clicks: %v", err)
	}

	var updatedLink Link
	if err := attributevalue.UnmarshalMap(result.Attributes, &updatedLink); err != nil {
		logger.Log.Error("failed to unmarshal updated link", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal updated link: %v", err)
	}

	logger.Log.Info("link clicks updated successfully", zap.String("short_url", updatedLink.ShortURL))
	return &updatedLink, nil
}
