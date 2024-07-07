package storage

import (
	"context"
	appConfig "dating-app-backend/internal/config"
	appLogger "dating-app-backend/internal/logger"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB struct {
	client *dynamodb.Client
	logger *appLogger.Logger
}

func NewDynamoDB(cfg *appConfig.Config, logger *appLogger.Logger) (*DynamoDB, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.AWSEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           cfg.AWSEndpoint,
				SigningRegion: cfg.AWSRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	defaultConfig, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(cfg.AWSRegion),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretKey, "")),
	)

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(defaultConfig)

	db := &DynamoDB{client: client, logger: logger}

	if err := db.createUsersTable(); err != nil {
		return nil, err
	}

	if err := db.createSwipesTable(); err != nil {
		return nil, err
	}

	return db, nil
}

// TODO: maybe move this to IaC
func (db *DynamoDB) createUsersTable() error {
	param := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("Email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String(usersTableName),
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("EmailIndex"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("Email"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}

	_, err := db.client.CreateTable(context.TODO(), param)
	if err != nil {
		var resourceInUseErr *types.ResourceInUseException
		if errors.As(err, &resourceInUseErr) {
			db.logger.Warn("Users table already exists")
			return nil
		}
		db.logger.Error("Failed to create Users table", "error", err)
		return err
	}

	db.logger.Info("Successfully created Users table")
	return nil
}

func (db *DynamoDB) createSwipesTable() error {
	_, err := db.client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("SwiperId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SwipedId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("SwiperId"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SwipedId"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(swipesTableName),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		var resourceInUseErr *types.ResourceInUseException
		if errors.As(err, &resourceInUseErr) {
			db.logger.Warn("Swipes table already exists")
			return nil
		}
		db.logger.Error("Failed to create Swipes table", "error", err)
		return err
	}
	db.logger.Info("Successfully created Swipes table")
	return nil
}
