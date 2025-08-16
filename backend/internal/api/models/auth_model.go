package models

import "github.com/google/uuid"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50,alpha_space"`
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,max=255,password_strength"`
}

type RegisterResponse struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt string    `db:"created_at" json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Password string    `db:"password_hash" json:"-"`
	Email    string    `db:"email" json:"email"`
	Role     UserRole  `db:"role" json:"role"`
}

type RoleResponse struct {
	Role UserRole `json:"role"`
}
