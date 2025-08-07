package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/otterly-id/otterly/backend/db"
	"github.com/otterly-id/otterly/backend/internal/api/models"
	"github.com/otterly-id/otterly/backend/internal/delivery/middlewares"
	"github.com/otterly-id/otterly/backend/internal/helpers"
	"github.com/otterly-id/otterly/backend/internal/utils"
	"go.uber.org/zap"
)

type AuthController struct {
	Log             *zap.Logger
	Validate        *validator.Validate
	ResponseHandler *helpers.ResponseHandler
	DB              *db.Queries
	JWTManager      *utils.JWTManager
}

func NewAuthController(logger *zap.Logger, validator *validator.Validate, db *db.Queries, jwtManager *utils.JWTManager) *AuthController {
	return &AuthController{
		Log:             logger,
		Validate:        validator,
		ResponseHandler: helpers.NewHandler(logger),
		DB:              db,
		JWTManager:      jwtManager,
	}
}

// Register func register new user.
// @Summary      Register
// @Description  Register new user.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body   models.RegisterRequest true "Register request"
// @Success      200  {object}  models.SuccessResponse[models.RegisterResponse]
// @Failure      400  {object}  models.FailureResponse[string]
// @Failure      404  {object}  models.FailureResponse[string]
// @Failure      500  {object}  models.FailureResponse[string]
// @Router       /api/auth/register [post]
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	newUser := &models.RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		ac.ResponseHandler.JSONDecodeError(w, r, err)
		return
	}

	if err := ac.Validate.Struct(newUser); err != nil {
		ac.ResponseHandler.ValidationError(w, r, err)
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		ac.ResponseHandler.HashPasswordError(w, r, err)
		return
	}
	newUser.Password = string(hashedPassword)

	user, err := ac.DB.Register(newUser)
	if err != nil {
		ac.ResponseHandler.CreateItemError(w, r, err, "user")
		return
	}

	ac.ResponseHandler.Success(w, r, http.StatusCreated, "User registered successfully", user)
}

// Login func login with credentials.
// @Summary      Login
// @Description  Login using email and password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body   models.LoginRequest true "Login request"
// @Success      200  {object}  models.SuccessResponse[models.RoleResponse]
// @Failure      400  {object}  models.FailureResponse[string]
// @Failure      404  {object}  models.FailureResponse[string]
// @Failure      500  {object}  models.FailureResponse[string]
// @Router       /api/auth/login [post]
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		ac.ResponseHandler.JSONDecodeError(w, r, err)
		return
	}

	if err := ac.Validate.Struct(user); err != nil {
		ac.ResponseHandler.ValidationError(w, r, err)
		return
	}

	foundUser, err := ac.DB.Login(user.Email)
	if err != nil {
		ac.ResponseHandler.NotFoundError(w, r, err, "User")
		return
	}

	if ok := utils.ComparePassword(user.Password, foundUser.Password); !ok {
		ac.ResponseHandler.AuthenticationFailedError(w, r, fmt.Errorf("invalid credentials provided"))
		return
	}

	token, duration, err := ac.JWTManager.GenerateToken(foundUser.ID.String(), foundUser.Email, foundUser.Role)
	if err != nil {
		ac.ResponseHandler.TokenGenerationError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "otterly_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(duration),
	})

	roleResponse := models.RoleResponse{
		Role: foundUser.Role,
	}

	ac.ResponseHandler.Success(w, r, http.StatusOK, "Login successful", roleResponse)
}

// GetAuthenticatedUser func get current authenticated user.
// @Summary      Get Authenticated User
// @Description  Get current authenticated user data.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  models.SuccessResponse[models.UserResponse]
// @Failure      400  {object}  models.FailureResponse[string]
// @Failure      401  {object}  models.FailureResponse[string]
// @Failure      404  {object}  models.FailureResponse[string]
// @Failure      500  {object}  models.FailureResponse[string]
// @Router       /api/auth/me [get]
func (ac *AuthController) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		ac.ResponseHandler.AuthenticationRequiredError(w, r)
		return
	}

	user, err := ac.DB.GetUser(userInfo.ID)
	if err != nil {
		ac.ResponseHandler.NotFoundError(w, r, err, "User")
		return
	}

	ac.ResponseHandler.Success(w, r, http.StatusOK, "User found", user)
}

// Logout func logs out the current user.
// @Summary      Logout
// @Description  Logout the current authenticated user by removing the JWT cookie.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  models.SuccessResponseWithoutData
// @Failure      401  {object}  models.FailureResponse[string]
// @Router       /api/auth/logout [post]
func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		ac.ResponseHandler.AuthenticationRequiredError(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "otterly_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	ac.ResponseHandler.Success(w, r, http.StatusOK, "Logout successful", nil)
}
