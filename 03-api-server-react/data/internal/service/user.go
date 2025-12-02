package service

import (
	"fmt"
	"time"

	"github.com/example/api-server/internal/model"
)

var users = []model.User{
	{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
	{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
}

func ListUsers() []model.User {
	return users
}

func FindUserByID(id int) (*model.User, error) {
	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func RegisterUser(req model.CreateUserRequest) model.User {
	newUser := model.User{
		ID:        len(users) + 1,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}
	users = append(users, newUser)
	return newUser
}
