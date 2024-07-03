package storage

import (
	"context"
	appConfig "dating-app-backend/internal/config"
	appLogger "dating-app-backend/internal/logger"
	appModel "dating-app-backend/internal/model"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB struct {
	client *dynamodb.Client
	logger *appLogger.Logger
}

func NewDynamoDB(cfg *appConfig.Config, logger *appLogger.Logger) (*DynamoDB, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(defaultConfig)

	db := &DynamoDB{client: client, logger: logger}

	if err := db.createUsersTable(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DynamoDB) createUsersTable() error {
	param := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String("Users"),
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

func (db *DynamoDB) CreateUser(ctx context.Context, user appModel.User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		db.logger.Error("Failed to marshal user", "error", err, "userId", user.ID)
		return err
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      av,
	})

	if err != nil {
		db.logger.Error("Failed to put item in DynamoDB", "error", err, "userId", user.ID)
		return err
	}

	db.logger.Info("Successfully created user in DynamoDB", "userId", user.ID)
	return nil
}
