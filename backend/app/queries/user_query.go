package queries

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/otterly-id/otterly/backend/app/models"
)

type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) CreateUser(u *models.CreateUserRequest) (models.CreateUserResponse, error) {
	var user models.CreateUserResponse

	if err := q.QueryRowx(
		`INSERT INTO users (name, full_name, email, password, phone_number, role)
         VALUES ($1, $2, $3, $4, $5, $6)
         RETURNING id, name, full_name, email, phone_number, role`,
		u.Name,
		u.FullName,
		u.Email,
		u.Password,
		u.PhoneNumber,
		u.Role,
	).StructScan(&user); err != nil {
		return models.CreateUserResponse{}, err
	}

	return user, nil
}

func (q *UserQueries) GetUsers() ([]models.UserResponse, error) {
	var user []models.UserResponse

	if err := q.Select(&user, `SELECT id, name, full_name, email, phone_number, role FROM users WHERE deleted_at IS NULL`); err != nil {
		return []models.UserResponse{}, err
	}

	return user, nil
}

func (q *UserQueries) GetUser(id uuid.UUID) (models.UserResponse, error) {
	var user models.UserResponse

	if err := q.Get(&user, `SELECT id, name, full_name, email, phone_number, role FROM users WHERE id = $1 AND deleted_at IS NULL`, id); err != nil {
		return models.UserResponse{}, err
	}

	return user, nil
}

func (q *UserQueries) UpdateUser(id uuid.UUID, u *models.UpdateUserRequest) (models.UpdateUserResponse, error) {
	setParts := []string{}
	args := []interface{}{id}
	argIndex := 2

	if u.Name != nil && *u.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *u.Name)
		argIndex++
	}

	if u.FullName != nil && *u.FullName != "" {
		setParts = append(setParts, fmt.Sprintf("full_name = $%d", argIndex))
		args = append(args, *u.FullName)
		argIndex++
	}

	if u.Email != nil && *u.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *u.Email)
		argIndex++
	}

	if u.PhoneNumber != nil && *u.PhoneNumber != "" {
		setParts = append(setParts, fmt.Sprintf("phone_number = $%d", argIndex))
		args = append(args, *u.PhoneNumber)
	}

	if len(setParts) == 0 {
		return models.UpdateUserResponse{}, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, "updated_at = NOW()")

	query := fmt.Sprintf(
		`UPDATE users SET %s
         WHERE id = $1 AND deleted_at IS NULL
         RETURNING id, name, full_name, email, phone_number, role, updated_at`,
		strings.Join(setParts, ", "))

	var user models.UpdateUserResponse
	if err := q.QueryRowx(query, args...).StructScan(&user); err != nil {
		return models.UpdateUserResponse{}, err
	}

	return user, nil
}

func (q *UserQueries) DeleteUser(id uuid.UUID) error {
	if _, err := q.Exec(`UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id); err != nil {
		return err
	}

	return nil
}
