package queries

import (
	"github.com/jmoiron/sqlx"
	"github.com/otterly-id/otterly/backend/internal/api/models"
)

type AuthQueries struct {
	*sqlx.DB
}

func (q *AuthQueries) Register(u *models.RegisterRequest) (models.RegisterResponse, error) {
	var user models.RegisterResponse

	if err := q.QueryRowx(
		`INSERT INTO users (name, email, password_hash, role)
         VALUES ($1, $2, $3, $4, $5, 'USER')
         RETURNING id, name, email, created_at`,
		u.Name,
		u.Email,
		u.Password,
	).StructScan(&user); err != nil {
		return models.RegisterResponse{}, err
	}

	return user, nil
}

func (q *AuthQueries) Login(email string) (models.LoginResponse, error) {
	var user models.LoginResponse

	if err := q.Get(&user, `SELECT id, password, email, role FROM users WHERE email = $1 AND deleted_at IS NULL`, email); err != nil {
		return models.LoginResponse{}, err
	}

	return user, nil
}
