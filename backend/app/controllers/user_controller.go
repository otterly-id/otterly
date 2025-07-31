package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/otterly-id/otterly/backend/app/models"
	"github.com/otterly-id/otterly/backend/app/validators"
	"github.com/otterly-id/otterly/backend/pkg/utils"
	"github.com/otterly-id/otterly/backend/platform/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var logger = utils.NewLogger()

// CreateUser func create single user.
// @Summary      Create User
// @Description  Add new user data.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponse[models.CreateUserResponse]
// @Failure      400  {object}  utils.FailureResponse[string]
// @Failure      404  {object}  utils.FailureResponse[string]
// @Failure      500  {object}  utils.FailureResponse[string]
// @Router       /api/users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := &models.CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		errMessage := "Failed to parse JSON body"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusBadRequest})
		return
	}

	validator := validators.UserValidator()
	if err := validator.Struct(newUser); err != nil {
		errMessage := "Validation failed"
		validationErrors := utils.ValidatorErrors(err)
		logger.Error(errMessage, zap.String("error", err.Error()),
			zap.Strings("validation_errors", validationErrors))
		utils.NewFailureResponse(errMessage, validationErrors).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusBadRequest})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		errMessage := "Failed to process password"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(w, nil)
		return
	}
	newUser.Password = string(hashedPassword)

	db, err := database.OpenDBConnection()
	if err != nil {
		errMessage := "Database connection failed"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(w, nil)
		return
	}
	defer db.Close()

	user, err := db.CreateUser(newUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			errMessage := "User with this email already exists"
			logger.Error(errMessage, zap.String("error", err.Error()))
			utils.NewFailureResponse(errMessage, err.Error()).Write(
				w, &utils.ResponseOptions{StatusCode: http.StatusConflict})
		} else {
			errMessage := "Failed to create user"
			logger.Error(errMessage, zap.String("error", err.Error()))
			utils.NewFailureResponse(errMessage, err.Error()).Write(
				w, nil)
		}
		return
	}

	utils.NewSuccessResponse("User created successfully", user).Write(w, nil)
}

// GetUsers func get all users.
// @Summary      Get All Users
// @Description  Get all users data.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponse[[]models.UserResponse]
// @Failure      400  {object}  utils.FailureResponse[string]
// @Failure      404  {object}  utils.FailureResponse[string]
// @Failure      500  {object}  utils.FailureResponse[string]
// @Router       /api/users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {

	db, err := database.OpenDBConnection()
	if err != nil {
		errMessage := "Database connection failed"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(w, nil)
		return
	}
	defer db.Close()

	users, err := db.GetUsers()
	if err != nil {
		errMessage := "Users not found"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusNotFound})
		return
	}

	utils.NewSuccessResponse("Users found", users).Write(w, nil)
}

// GetUser func get user by ID.
// @Summary      Get User by ID
// @Description  Get user data based on provided ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponse[models.UserResponse]
// @Failure      400  {object}  utils.FailureResponse[string]
// @Failure      404  {object}  utils.FailureResponse[string]
// @Failure      500  {object}  utils.FailureResponse[string]
// @Router       /api/users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		errMessage := "Invalid ID"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusBadRequest})
		return
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		errMessage := "Database connection failed"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(w, nil)
		return
	}
	defer db.Close()

	user, err := db.GetUser(id)
	if err != nil {
		errMessage := "User not found"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusNotFound})
		return
	}

	utils.NewSuccessResponse("User found", user).Write(w, nil)
}

// UpdateUser func update single user.
// @Summary      Update User
// @Description  Edit user data based on provided ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponse[models.UpdateUserResponse]
// @Failure      400  {object}  utils.FailureResponse[string]
// @Failure      404  {object}  utils.FailureResponse[string]
// @Failure      500  {object}  utils.FailureResponse[string]
// @Router       /api/users/{id} [patch]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		errMessage := "Invalid ID"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusBadRequest})
		return
	}

	selectedUser := &models.UpdateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(selectedUser); err != nil {
		errMessage := "Failed to parse JSON body"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusBadRequest})
		return
	}

	validator := validators.UserValidator()
	if err := validator.Struct(selectedUser); err != nil {
		errMessage := "Validation failed"
		validationErrors := utils.ValidatorErrors(err)
		logger.Error(errMessage, zap.String("error", err.Error()),
			zap.Strings("validation_errors", validationErrors))
		utils.NewFailureResponse(errMessage, validationErrors).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusBadRequest})
		return
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		errMessage := "Database connection failed"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(w, nil)
		return
	}
	defer db.Close()

	user, err := db.UpdateUser(id, selectedUser)
	if err != nil {
		errMessage := "Failed to update user"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(w, nil)
		return
	}

	utils.NewSuccessResponse("User updated successfully", user).Write(w, nil)
}

// DeleteUser func delete single user.
// @Summary      Delete User
// @Description  Remove user data based on provided ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponseWithoutData
// @Failure      400  {object}  utils.FailureResponse[string]
// @Failure      404  {object}  utils.FailureResponse[string]
// @Failure      500  {object}  utils.FailureResponse[string]
// @Router       /api/users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		errMessage := "Invalid ID"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusBadRequest})
		return
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		errMessage := "Database connection failed"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(w, nil)
		return
	}
	defer db.Close()

	if err := db.DeleteUser(id); err != nil {
		errMessage := "User not found"
		logger.Error(errMessage, zap.String("error", err.Error()))
		utils.NewFailureResponse(errMessage, err.Error()).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusNotFound})
		return
	}

	utils.NewSuccessResponse[interface{}]("User deleted successfully", nil).Write(w, nil)
}
