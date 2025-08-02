package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/otterly-id/otterly/backend/app/models"
	"github.com/otterly-id/otterly/backend/app/validators"
	"github.com/otterly-id/otterly/backend/pkg/helpers"
	"github.com/otterly-id/otterly/backend/platform/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	logger         *zap.Logger
	errorHandler   *helpers.ErrorHandler
	successHandler *helpers.SuccessHandler
}

func NewUserController(logger *zap.Logger) *UserController {
	return &UserController{
		logger:         logger,
		errorHandler:   helpers.NewErrorHandler(logger),
		successHandler: helpers.NewSuccessHandler(logger),
	}
}

// CreateUser func create single user.
// @Summary      Create User
// @Description  Add new user data.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param		 request body   models.CreateUserRequest true "Create user request"
// @Success      200  {object}  helpers.SuccessResponse[models.CreateUserResponse]
// @Failure      400  {object}  helpers.FailureResponse[string]
// @Failure      404  {object}  helpers.FailureResponse[string]
// @Failure      500  {object}  helpers.FailureResponse[string]
// @Router       /api/users [post]
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := &models.CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		uc.errorHandler.JSONDecodeError(w, r, err)
		return
	}

	validator := validators.UserValidator()
	if err := validator.Struct(newUser); err != nil {
		uc.errorHandler.ValidationError(w, r, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.errorHandler.PasswordError(w, r, err)
		return
	}
	newUser.Password = string(hashedPassword)

	db, err := database.GetDBConnection()
	if err != nil {
		uc.errorHandler.DBConnectionError(w, r, err)
		return
	}

	user, err := db.CreateUser(newUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			uc.errorHandler.DuplicateKeyError(w, r, err, "user")
		} else {
			uc.errorHandler.CreateError(w, r, err, "user")
		}
		return
	}

	uc.successHandler.WithData(w, r, "User created successfully", user, &helpers.ResponseOptions{StatusCode: http.StatusCreated})
}

// GetUsers func get all users.
// @Summary      Get All Users
// @Description  Get all users data.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  helpers.SuccessResponse[[]models.UserResponse]
// @Failure      400  {object}  helpers.FailureResponse[string]
// @Failure      404  {object}  helpers.FailureResponse[string]
// @Failure      500  {object}  helpers.FailureResponse[string]
// @Router       /api/users [get]
func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.GetDBConnection()
	if err != nil {
		uc.errorHandler.DBConnectionError(w, r, err)
		return
	}

	users, err := db.GetUsers()
	if err != nil {
		uc.errorHandler.NotFoundError(w, r, err, "Users")
		return
	}

	uc.successHandler.WithData(w, r, "Users found", users, nil)
}

// GetUser func get user by ID.
// @Summary      Get User by ID
// @Description  Get user data based on provided ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param id 	 path string true "User ID"
// @Success      200  {object}  helpers.SuccessResponse[models.UserResponse]
// @Failure      400  {object}  helpers.FailureResponse[string]
// @Failure      404  {object}  helpers.FailureResponse[string]
// @Failure      500  {object}  helpers.FailureResponse[string]
// @Router       /api/users/{id} [get]
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		uc.errorHandler.InvalidIDError(w, r, err)
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		uc.errorHandler.InvalidIDError(w, r, err)
		return
	}

	db, err := database.GetDBConnection()
	if err != nil {
		uc.errorHandler.DBConnectionError(w, r, err)
		return
	}

	user, err := db.GetUser(parsedId)
	if err != nil {
		uc.errorHandler.NotFoundError(w, r, err, "User")
		return
	}

	uc.successHandler.WithData(w, r, "User found", user, nil)
}

// UpdateUser func update single user.
// @Summary      Update User
// @Description  Edit user data based on provided ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param id 	 path string true "User ID"
// @Param		 request body   models.UpdateUserRequest true "Update user request"
// @Success      200  {object}  helpers.SuccessResponse[models.UpdateUserResponse]
// @Failure      400  {object}  helpers.FailureResponse[string]
// @Failure      404  {object}  helpers.FailureResponse[string]
// @Failure      500  {object}  helpers.FailureResponse[string]
// @Router       /api/users/{id} [patch]
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		uc.errorHandler.InvalidIDError(w, r, err)
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		uc.errorHandler.InvalidIDError(w, r, err)
		return
	}

	selectedUser := &models.UpdateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(selectedUser); err != nil {
		uc.errorHandler.JSONDecodeError(w, r, err)
		return
	}

	validator := validators.UserValidator()
	if err := validator.Struct(selectedUser); err != nil {
		uc.errorHandler.ValidationError(w, r, err)
		return
	}

	db, err := database.GetDBConnection()
	if err != nil {
		uc.errorHandler.DBConnectionError(w, r, err)
		return
	}

	user, err := db.UpdateUser(parsedId, selectedUser)
	if err != nil {
		uc.errorHandler.UpdateError(w, r, err, "User")
		return
	}

	uc.successHandler.WithData(w, r, "User updated successfully", user, nil)
}

// DeleteUser func delete single user.
// @Summary      Delete User
// @Description  Remove user data based on provided ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param id	 path string true "User ID"
// @Success      200  {object}  helpers.SuccessResponseWithoutData
// @Failure      400  {object}  helpers.FailureResponse[string]
// @Failure      404  {object}  helpers.FailureResponse[string]
// @Failure      500  {object}  helpers.FailureResponse[string]
// @Router       /api/users/{id} [delete]
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		uc.errorHandler.InvalidIDError(w, r, err)
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		uc.errorHandler.InvalidIDError(w, r, err)
		return
	}

	db, err := database.GetDBConnection()
	if err != nil {
		uc.errorHandler.DBConnectionError(w, r, err)
		return
	}

	if err := db.DeleteUser(parsedId); err != nil {
		uc.errorHandler.DeleteError(w, r, err, "User")
		return
	}

	uc.successHandler.WithoutData(w, r, "User deleted successfully", nil)
}
