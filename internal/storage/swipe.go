package storage

import (
	"context"
	"dating-app-backend/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const swipesTableName = "SwipesTable"

func (db *DynamoDB) RecordSwipe(ctx context.Context, swipe model.Swipe) (bool, string, error) {
	db.logger.Info("Recording swipe", "swiperId", swipe.SwiperId, "swipedId", swipe.SwipedId, "preference", swipe.Preference)

	// Update swipe statistics for the swiped user
	swipedUser, err := db.GetUserByID(ctx, swipe.SwipedId)
	if err != nil {
		db.logger.Error("Failed to get swiped user", "error", err, "swipedId", swipe.SwipedId)
		return false, "", err
	}

	swipedUser.TotalSwipes++
	if swipe.Preference == model.SwipeYes {
		swipedUser.YesSwipes++
	}
	swipedUser.UpdateAttractivenessScore()

	// TODO: maybe do this using the listener on the DynamoDB stream
	// Update the swiped user in the database
	err = db.UpdateUser(ctx, swipedUser)
	if err != nil {
		db.logger.Error("Failed to update swiped user", "error", err, "swipedId", swipe.SwipedId)
		return false, "", err
	}

	item, err := attributevalue.MarshalMap(swipe)
	if err != nil {
		db.logger.Error("Failed to marshal swipe", "error", err)
		return false, "", err
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(swipesTableName),
		Item:      item,
	})
	if err != nil {
		db.logger.Error("Failed to put swipe in DynamoDB", "error", err)
		return false, "", err
	}

	// TODO: A lambda handler on the DynamoDB stream could be used to check for matches

	// Check for a match
	if swipe.Preference == model.SwipeYes {
		matchResult, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(swipesTableName),
			Key: map[string]types.AttributeValue{
				"SwiperId": &types.AttributeValueMemberS{Value: swipe.SwipedId},
				"SwipedId": &types.AttributeValueMemberS{Value: swipe.SwiperId},
			},
		})
		if err != nil {
			db.logger.Error("Failed to check for match", "error", err)
			return false, "", err
		}

		if matchResult.Item != nil {
			var matchSwipe model.Swipe
			err = attributevalue.UnmarshalMap(matchResult.Item, &matchSwipe)
			if err != nil {
				db.logger.Error("Failed to unmarshal match swipe", "error", err)
				return false, "", err
			}

			if matchSwipe.Preference == model.SwipeYes {
				db.logger.Info("Match found", "swiperId", swipe.SwiperId, "swipedId", swipe.SwipedId)
				return true, swipe.SwipedId, nil
			}
		}
	}

	db.logger.Info("Swipe recorded successfully", "swiperId", swipe.SwiperId, "swipedId", swipe.SwipedId)
	return false, "", nil
}
