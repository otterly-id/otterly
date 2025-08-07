package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/otterly-id/otterly/backend/db"
	"github.com/otterly-id/otterly/backend/internal/api/models"
	"github.com/otterly-id/otterly/backend/internal/helpers"
	"github.com/otterly-id/otterly/backend/internal/utils"
	"go.uber.org/zap"
)

type UserController struct {
	Log             *zap.Logger
	Validate        *validator.Validate
	ResponseHandler *helpers.ResponseHandler
	DB              *db.Queries
}

func NewUserController(logger *zap.Logger, validator *validator.Validate, db *db.Queries) *UserController {
	return &UserController{
		Log:             logger,
		Validate:        validator,
		ResponseHandler: helpers.NewHandler(logger),
		DB:              db,
	}
}

// CreateUser func create single user.
// @Summary      Create User
// @Description  Add new user data.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        request body   models.CreateUserRequest true "Create user request"
// @Success      200  {object}  models.SuccessResponse[models.CreateUserResponse]
// @Failure      400  {object}  models.FailureResponse[string]
// @Failure      404  {object}  models.FailureResponse[string]
// @Failure      500  {object}  models.FailureResponse[string]
// @Router       /api/users [post]
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := &models.CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		uc.ResponseHandler.JSONDecodeError(w, r, err)
		return
	}

	if err := uc.Validate.Struct(newUser); err != nil {
		uc.ResponseHandler.ValidationError(w, r, err)
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		uc.ResponseHandler.HashPasswordError(w, r, err)
		return
	}
	newUser.Password = string(hashedPassword)

	user, err := uc.DB.CreateUser(newUser)
	if err != nil {
		uc.ResponseHandler.CreateItemError(w, r, err, "user")
		return
	}

	uc.ResponseHandler.Success(w, r, http.StatusCreated, "User created successfully", user)
}

// GetUsers func get all users.
// @Summary      Get All Users
// @Description  Get all users data.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  models.SuccessResponse[[]models.UserResponse]
// @Failure      400  {object}  models.FailureResponse[string]
// @Failure      404  {object}  models.FailureResponse[string]
// @Failure      500  {object}  models.FailureResponse[string]
// @Router       /api/users [get]
func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.DB.GetUsers()
	if err != nil {
		uc.ResponseHandler.NotFoundError(w, r, err, "Users")
		return
	}

	uc.ResponseHandler.Success(w, r, http.StatusOK, "Users found", users)
}

// GetUser func get user by ID.
// @Summary      Get User by ID
// @Description  Get user data based on provided ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param id 	 path string true "User ID"
// @Success      200  {object}  models.SuccessResponse[models.UserResponse]
// @Failure      400  {object}  models.FailureResponse[string]
// @Failure      404  {object}  models.FailureResponse[string]
// @Failure      500  {object}  models.FailureResponse[string]
// @Router       /api/users/{id} [get]
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		uc.ResponseHandler.InvalidIDError(w, r, err)
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		uc.ResponseHandler.InvalidIDError(w, r, err)
		return
	}

	user, err := uc.DB.GetUser(parsedId)
	if err != nil {
		uc.ResponseHandler.NotFoundError(w, r, err, "User")
		return
	}

	uc.ResponseHandler.Success(w, r, http.StatusOK, "User found", user)
}

// UpdateUser func update single user.
// @Summary      Update User
// @Description  Edit user data based on provided ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param id 	 path string true "User ID"
// @Param		 request body   models.UpdateUserRequest true "Update user request"
// @Success      200  {object}  models.SuccessResponse[models.UpdateUserResponse]
// @Failure      400  {object}  models.FailureResponse[string]
// @Failure      404  {object}  models.FailureResponse[string]
// @Failure      500  {object}  models.FailureResponse[string]
// @Router       /api/users/{id} [patch]
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		uc.ResponseHandler.InvalidIDError(w, r, err)
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		uc.ResponseHandler.InvalidIDError(w, r, err)
		return
	}

	selectedUser := &models.UpdateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(selectedUser); err != nil {
		uc.ResponseHandler.JSONDecodeError(w, r, err)
		return
	}

	if err := uc.Validate.Struct(selectedUser); err != nil {
		uc.ResponseHandler.ValidationError(w, r, err)
		return
	}

	user, err := uc.DB.UpdateUser(parsedId, selectedUser)
	if err != nil {
		uc.ResponseHandler.UpdateItemError(w, r, err, "User")
		return
	}

	uc.ResponseHandler.Success(w, r, http.StatusCreated, "User updated successfully", user)
}

// DeleteUser func delete single user.
// @Summary      Delete User
// @Description  Remove user data based on provided ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param id	 path string true "User ID"
// @Success      200  {object}  models.SuccessResponseWithoutData
// @Failure      400  {object}  models.FailureResponse[string]
// @Failure      404  {object}  models.FailureResponse[string]
// @Failure      500  {object}  models.FailureResponse[string]
// @Router       /api/users/{id} [delete]
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		uc.ResponseHandler.InvalidIDError(w, r, err)
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		uc.ResponseHandler.InvalidIDError(w, r, err)
		return
	}

	if err := uc.DB.DeleteUser(parsedId); err != nil {
		uc.ResponseHandler.DeleteItemError(w, r, err, "User")
		return
	}

	uc.ResponseHandler.Success(w, r, http.StatusOK, "User deleted successfully", nil)
}
