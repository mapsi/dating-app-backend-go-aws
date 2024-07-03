package storage

import (
	"context"
	appConfig "dating-app-backend/internal/config"
	appLogger "dating-app-backend/internal/logger"
	appModel "dating-app-backend/internal/model"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB struct {
	client *dynamodb.Client
	logger *appLogger.Logger
}

const tableName = "UsersTable"

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
		TableName: aws.String(tableName),
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

func (db *DynamoDB) CreateUser(ctx context.Context, user appModel.User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		db.logger.Error("Failed to marshal user", "error", err, "userId", user.ID)
		return err
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	if err != nil {
		db.logger.Error("Failed to put item in DynamoDB", "error", err, "userId", user.ID)
		return err
	}

	db.logger.Info("Successfully created user in DynamoDB", "userId", user.ID)
	return nil
}

func (db *DynamoDB) GetUserByEmail(ctx context.Context, email string) (*appModel.User, error) {
	result, err := db.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("EmailIndex"),
		KeyConditionExpression: aws.String("Email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	})

	if err != nil {
		db.logger.Error("Failed to query user by email", "error", err, "email", email)
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("user not found")
	}

	var user appModel.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		db.logger.Error("Failed to unmarshal user", "error", err, "email", email)
		return nil, err
	}

	return &user, nil
}
