package storage

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	appModel "dating-app-backend/internal/model"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jftuga/geodist"
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

func (db *DynamoDB) DiscoverUsers(ctx context.Context, currentUser appModel.User, limit int32, minAge, maxAge int, gender string) ([]appModel.UserPublicData, error) {
	db.logger.Info("Discovering users", "currentUserID", currentUser.ID, "limit", limit, "minAge", minAge, "maxAge", maxAge, "gender", gender)

	// Get all swipes by the current user
	swipedUsers, err := db.getSwipedUsers(ctx, currentUser.ID)
	if err != nil {
		db.logger.Error("Failed to get swiped users", "error", err, "currentUserID", currentUser.ID)
		return nil, err
	}

	// Prepare the filter expression
	filterExp := "ID <> :currentUserId"
	expAttrValues := map[string]types.AttributeValue{
		":currentUserId": &types.AttributeValueMemberS{Value: currentUser.ID},
	}

	// Add swiped users to the filter expression
	for i, swipedID := range swipedUsers {
		filterExp += fmt.Sprintf(" AND ID <> :swipedId%d", i)
		expAttrValues[fmt.Sprintf(":swipedId%d", i)] = &types.AttributeValueMemberS{Value: swipedID}
	}

	if minAge > 0 {
		filterExp += " AND Age >= :minAge"
		expAttrValues[":minAge"] = &types.AttributeValueMemberN{Value: strconv.Itoa(minAge)}
	}

	if maxAge > 0 {
		filterExp += " AND Age <= :maxAge"
		expAttrValues[":maxAge"] = &types.AttributeValueMemberN{Value: strconv.Itoa(maxAge)}
	}

	if gender != "" {
		filterExp += " AND Gender = :gender"
		expAttrValues[":gender"] = &types.AttributeValueMemberS{Value: gender}
	}

	// Perform the scan operation
	result, err := db.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(usersTableName),
		FilterExpression:          aws.String(filterExp),
		ExpressionAttributeValues: expAttrValues,
		Limit:                     aws.Int32(limit),
	})
	if err != nil {
		db.logger.Error("Failed to scan users for discovery", "error", err, "currentUserID", currentUser.ID)
		return nil, err
	}

	var users []appModel.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		db.logger.Error("Failed to unmarshal discovered users", "error", err, "currentUserID", currentUser.ID)
		return nil, err
	}

	publicUsers := make([]appModel.UserPublicData, len(users))
	for i, user := range users {
		publicData := user.PublicData()
		distance, _ := geodist.HaversineDistance(geodist.Coord{Lat: currentUser.Latitude, Lon: currentUser.Longitude},
			geodist.Coord{Lat: user.Latitude, Lon: user.Longitude})
		publicData.DistanceFromMe = distance
		publicUsers[i] = publicData
	}

	// Sort users by distance
	sort.Slice(publicUsers, func(i, j int) bool {
		return publicUsers[i].DistanceFromMe < publicUsers[j].DistanceFromMe
	})

	db.logger.Info("Users discovered successfully", "currentUserID", currentUser.ID, "count", len(publicUsers))
	return publicUsers, nil
}

func (db *DynamoDB) getSwipedUsers(ctx context.Context, swiperId string) ([]string, error) {
	result, err := db.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(swipesTableName),
		KeyConditionExpression: aws.String("SwiperId = :swiperId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":swiperId": &types.AttributeValueMemberS{Value: swiperId},
		},
		ProjectionExpression: aws.String("SwipedId"),
	})
	if err != nil {
		return nil, err
	}

	var swipes []struct {
		SwipedId string `dynamodbav:"SwipedId"`
	}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &swipes)
	if err != nil {
		return nil, err
	}

	swipedIds := make([]string, len(swipes))
	for i, swipe := range swipes {
		swipedIds[i] = swipe.SwipedId
	}

	return swipedIds, nil
}

func (db *DynamoDB) GetUserByID(ctx context.Context, userID string) (*appModel.User, error) {
	db.logger.Info("Getting user by ID", "userID", userID)

	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(usersTableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		db.logger.Error("Failed to get user from DynamoDB", "error", err, "userID", userID)
		return nil, err
	}

	if result.Item == nil {
		db.logger.Warn("User not found", "userID", userID)
		return nil, fmt.Errorf("user not found")
	}

	var user appModel.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		db.logger.Error("Failed to unmarshal user data", "error", err, "userID", userID)
		return nil, err
	}

	db.logger.Info("User retrieved successfully", "userID", userID)
	return &user, nil
}
