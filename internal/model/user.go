package model

import (
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	"github.com/oklog/ulid/v2"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func GenerateRandomUser() User {

	age, _ := faker.RandomInt(18, 100)

	return User{
		// Generate a ULID for the ID
		// https://github.com/oklog/ulid "Care should be taken when providing a source of entropy."
		ID:    ulid.MustNew(ulid.Now(), nil).String(),
		Email: faker.Email(),
		// TODO: Use a more secure password hashing algorithm
		Password: faker.Password(),
		Name:     faker.Name(),
		Gender:   randomGender(),
		Age:      age[0],
	}
}
func (u *User) CheckPassword(password string) bool {
	return u.Password == password
}

func randomGender() string {
	genders := []string{"Male", "Female"}

	return genders[rand.IntN(len(genders))]
}
