package model

import (
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/oklog/ulid/v2"
)

type User struct {
	ID                  string  `json:"id" dynamodbav:"ID"`
	Email               string  `json:"email" dynamodbav:"Email"`
	Password            string  `json:"-" dynamodbav:"Password"`
	Name                string  `json:"name" dynamodbav:"Name"`
	Gender              string  `json:"gender" dynamodbav:"Gender"`
	Age                 int     `json:"age" dynamodbav:"Age"`
	Latitude            float64 `json:"latitude" dynamodbav:"Latitude"`
	Longitude           float64 `json:"longitude" dynamodbav:"Longitude"`
	YesSwipes           int     `json:"yesSwipes" dynamodbav:"YesSwipes"`
	TotalSwipes         int     `json:"totalSwipes" dynamodbav:"TotalSwipes"`
	AttractivenessScore float64 `json:"attractivenessScore" dynamodbav:"AttractivenessScore"`
}

type UserPublicData struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	Gender              string  `json:"gender"`
	Age                 int     `json:"age"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	DistanceFromMe      float64 `json:"distanceFromMe"`
	AttractivenessScore float64 `json:"attractivenessScore"`
}

func (u *User) PublicData() UserPublicData {
	return UserPublicData{
		ID:                  u.ID,
		Name:                u.Name,
		Gender:              u.Gender,
		Age:                 u.Age,
		Latitude:            u.Latitude,
		Longitude:           u.Longitude,
		AttractivenessScore: u.AttractivenessScore,
	}
}

func (u *User) UpdateAttractivenessScore() {
	if u.TotalSwipes > 0 {
		u.AttractivenessScore = float64(u.YesSwipes) / float64(u.TotalSwipes)
	} else {
		u.AttractivenessScore = 0.5 // Default score for new users
	}
}

func GenerateRandomUser() User {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	return User{
		ID:    id,
		Email: faker.Email(),
		// TODO: Use a more secure password hashing algorithm
		Password:  faker.Password(),
		Name:      faker.Name(),
		Gender:    randomGender(),
		Age:       rand.Intn(62) + 18,
		Latitude:  randomLatitude(),
		Longitude: randomLongitude(),
	}
}

func randomLatitude() float64 {
	return rand.Float64()*180 - 90
}

func randomLongitude() float64 {
	return rand.Float64()*360 - 180
}

func (u *User) CheckPassword(password string) bool {
	return u.Password == password
}

func randomGender() string {
	genders := []string{"Male", "Female"}

	return genders[rand.Intn(len(genders))]
}
