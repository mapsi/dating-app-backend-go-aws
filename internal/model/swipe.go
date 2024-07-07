package model

import (
	"time"
)

type SwipePreference string

const (
	SwipeYes SwipePreference = "YES"
	SwipeNo  SwipePreference = "NO"
)

type Swipe struct {
	SwiperId   string          `json:"swiperId" dynamodbav:"SwiperId"`
	SwipedId   string          `json:"swipedId" dynamodbav:"SwipedId"`
	Preference SwipePreference `json:"preference" dynamodbav:"Preference"`
	CreatedAt  time.Time       `json:"createdAt" dynamodbav:"CreatedAt"`
}
