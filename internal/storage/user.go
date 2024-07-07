package storage

import (
	"context"
	appModel "dating-app-backend/internal/model"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const usersTableName = "UsersTable"

func (db *DynamoDB) CreateUser(ctx context.Context, user appModel.User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		db.logger.Error("Failed to marshal user", "error", err, "userId", user.ID)
		return err
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(usersTableName),
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
		TableName:              aws.String(usersTableName),
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

func (db *DynamoDB) DiscoverUsers(ctx context.Context, currentUserID string, limit int32) ([]appModel.UserPublicData, error) {
	// Query DynamoDB to get all users except the current user
	input := &dynamodb.ScanInput{
		TableName:        aws.String(usersTableName),
		FilterExpression: aws.String("ID <> :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: currentUserID},
		},
		Limit: aws.Int32(limit),
	}

	result, err := db.client.Scan(ctx, input)
	if err != nil {
		db.logger.Error("Failed to scan users for discovery", "error", err, "currentUserID", currentUserID)
		return nil, err
	}

	// Unmarshal the results
	var users []appModel.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		db.logger.Error("Failed to unmarshal discovered users", "error", err, "currentUserID", currentUserID)
		return nil, err
	}

	publicUsers := make([]appModel.UserPublicData, len(users))
	for i, user := range users {
		publicUsers[i] = user.PublicData()
	}

	// TODO: Filter out users that have already been swiped

	db.logger.Info("Users discovered successfully", "currentUserID", currentUserID, "count", len(users))
	return publicUsers, nil
}
