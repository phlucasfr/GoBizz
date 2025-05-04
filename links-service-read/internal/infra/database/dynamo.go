package database

import (
	"context"
	"errors"
	"fmt"
	"links-service-read/internal/logger"
	"links-service-read/utils"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

// NewDynamoClient initializes and returns a new DynamoDB client.
// It configures the AWS SDK based on the environment (production or development/local).
// In production, it uses the default AWS configuration with the "2" region.
// In development/local, it uses a custom endpoint specified by the DYNAMODB_ENDPOINT environment variable.
//
// The function also ensures the existence of the "Links" table in DynamoDB.
//
// Returns:
// - *dynamodb.Client: A pointer to the initialized DynamoDB client.
// - error: An error if the client initialization or table verification fails.
func NewDynamoClient() (*dynamodb.Client, error) {
	var cfg aws.Config
	var err error

	ctx := context.Background()
	endpoint := utils.ConfigInstance.DynamoEndpoint
	if endpoint == "" {
		logger.Log.Error("Failed to create DynamoDB client: DYNAMODB_ENDPOINT is not set")
		return nil, fmt.Errorf("DYNAMODB_ENDPOINT is not set")
	}

	if os.Getenv("ENVIRONMENT") == "production" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion("us-east-2"),
		)
		if err != nil {
			logger.Log.Error("Failed to load AWS config (production)", zap.Error(err))
			return nil, fmt.Errorf("failed to load AWS config (production): %v", err)
		}
	} else {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion("us-east-2"),
			config.WithEndpointResolver(aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint, SigningRegion: "us-east-2"}, nil
				})),
		)
		if err != nil {
			logger.Log.Error("Failed to load AWS config (development/local)", zap.Error(err))
			return nil, fmt.Errorf("failed to load AWS config (development/local): %v", err)
		}
	}

	client := dynamodb.NewFromConfig(cfg)

	if err := ensureLinksTable(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

func ensureLinksTable(ctx context.Context, db *dynamodb.Client) error {
	_, err := db.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String("Links"),
	})
	if err == nil {
		logger.Log.Info("DynamoDB table 'Links' already exists")
		return nil
	}

	var rnfe *types.ResourceNotFoundException
	if !errors.As(err, &rnfe) {
		logger.Log.Error("Error describing DynamoDB table", zap.Error(err))
		return fmt.Errorf("error describing DynamoDB table: %v", err)
	}

	attributeDefinitions := []types.AttributeDefinition{
		{AttributeName: aws.String("short_url"), AttributeType: types.ScalarAttributeTypeS},
		{AttributeName: aws.String("custom_slug"), AttributeType: types.ScalarAttributeTypeS},
		{AttributeName: aws.String("customer_id"), AttributeType: types.ScalarAttributeTypeS},
		{AttributeName: aws.String("created_at"), AttributeType: types.ScalarAttributeTypeS},
		{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
	}

	createInput := &dynamodb.CreateTableInput{
		TableName:            aws.String("Links"),
		AttributeDefinitions: attributeDefinitions,
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("short_url"), KeyType: types.KeyTypeHash},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("ByCustomSlug"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("custom_slug"), KeyType: types.KeyTypeHash},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(2),
					WriteCapacityUnits: aws.Int64(2),
				},
			},
			{
				IndexName: aws.String("ByCustomer"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("customer_id"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("created_at"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(2),
					WriteCapacityUnits: aws.Int64(2),
				},
			},
			{
				IndexName: aws.String("ByID"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(2),
					WriteCapacityUnits: aws.Int64(2),
				},
			},
		},
	}

	_, err = db.CreateTable(ctx, createInput)
	if err != nil {
		logger.Log.Error("Failed to create DynamoDB table", zap.Error(err))
		return fmt.Errorf("failed to create DynamoDB table: %v", err)
	}

	waiter := dynamodb.NewTableExistsWaiter(db)
	maxWaitTime := 5 * time.Minute
	describeInput := &dynamodb.DescribeTableInput{
		TableName: aws.String("Links"),
	}

	if err := waiter.Wait(ctx, describeInput, maxWaitTime); err != nil {
		logger.Log.Error("Failed to wait for table creation", zap.Error(err))
		return fmt.Errorf("failed to wait for table creation: %v", err)
	}

	_, err = db.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String("Links"),
		TimeToLiveSpecification: &types.TimeToLiveSpecification{
			Enabled:       aws.Bool(true),
			AttributeName: aws.String("ttl"),
		},
	})
	if err != nil {
		logger.Log.Error("Failed to enable TTL on DynamoDB table", zap.Error(err))
		return fmt.Errorf("failed to enable TTL: %v", err)
	}

	logger.Log.Info("DynamoDB table 'Links' created successfully with TTL enabled")
	return nil
}
