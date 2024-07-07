package model

import (
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/oklog/ulid/v2"
)

type User struct {
	ID       string `json:"id" dynamodbav:"ID"`
	Email    string `json:"email" dynamodbav:"Email"`
	Password string `json:"-" dynamodbav:"Password"`
	Name     string `json:"name" dynamodbav:"Name"`
	Gender   string `json:"gender" dynamodbav:"Gender"`
	Age      int    `json:"age" dynamodbav:"Age"`
}

type UserPublicData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

func (u *User) PublicData() UserPublicData {
	return UserPublicData{
		ID:     u.ID,
		Name:   u.Name,
		Gender: u.Gender,
		Age:    u.Age,
	}
}

func GenerateRandomUser() User {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	return User{
		ID:    id,
		Email: faker.Email(),
		// TODO: Use a more secure password hashing algorithm
		Password: faker.Password(),
		Name:     faker.Name(),
		Gender:   randomGender(),
		Age:      rand.Intn(82) + 18,
	}
}
func (u *User) CheckPassword(password string) bool {
	return u.Password == password
}

func randomGender() string {
	genders := []string{"Male", "Female"}

	return genders[rand.Intn(len(genders))]
}
