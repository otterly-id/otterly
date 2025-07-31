package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
	RoleOwner UserRole = "OWNER"
)

type User struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	FullName    string     `json:"full_name"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	PhoneNumber string     `json:"phone_number"`
	Role        UserRole   `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type CreateUserRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50,alpha_space"`
	FullName    string `json:"full_name" validate:"max=100"`
	Email       string `json:"email" validate:"required,email,max=254"`
	Password    string `json:"password" validate:"required,min=8,max=255,password_strength"`
	PhoneNumber string `json:"phone_number" validate:"max=20,phone"`
	Role        string `json:"role" validate:"required,oneof=USER OWNER"`
}

type CreateUserResponse struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	FullName    string    `db:"full_name" json:"full_name,omitempty"`
	Email       string    `db:"email" json:"email"`
	PhoneNumber string    `db:"phone_number" json:"phone_number,omitempty"`
	Role        UserRole  `db:"role" json:"role"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type UserResponse struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	FullName    string    `db:"full_name" json:"full_name"`
	Email       string    `db:"email" json:"email"`
	PhoneNumber string    `db:"phone_number" json:"phone_number"`
	Role        UserRole  `db:"role" json:"role"`
}

type UpdateUserRequest struct {
	Name        *string `json:"name" validate:"omitempty,max=50,alpha_space"`
	FullName    *string `json:"full_name" validate:"omitempty,max=100"`
	Email       *string `json:"email" validate:"omitempty,email,max=254"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,max=20,phone"`
}

type UpdateUserResponse struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	FullName    string    `db:"full_name" json:"full_name,omitempty"`
	Email       string    `db:"email" json:"email"`
	PhoneNumber string    `db:"phone_number" json:"phone_number,omitempty"`
	Role        UserRole  `db:"role" json:"role"`
	UpdatedAt   string    `db:"updated_at" json:"updated_at"`
}
