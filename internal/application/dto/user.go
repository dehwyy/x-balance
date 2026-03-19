package dto

import (
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type UserCreateRequest struct {
	Name           user.Name
	OverdraftLimit user.OverdraftLimit
}

type UserCreateResponse struct {
	User user.User
}

type UserGetByIDRequest struct {
	ID user.ID
}

type UserGetByIDResponse struct {
	User user.User
}

type UserUpdateRequest struct {
	User user.User
}

type UserUpdateResponse struct {
	User user.User
}

type UserDeleteRequest struct {
	ID user.ID
}
