package middlewares

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/otterly-id/otterly/backend/internal/api/models"
	"github.com/otterly-id/otterly/backend/internal/helpers"
	"github.com/otterly-id/otterly/backend/internal/utils"
	"go.uber.org/zap"
)

type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

type UserInfo struct {
	ID   uuid.UUID       `json:"id"`
	Role models.UserRole `json:"role"`
}

type AuthMiddleware struct {
	JWTManager      *utils.JWTManager
	ResponseHandler *helpers.ResponseHandler
	Log             *zap.Logger
}

func NewAuthMiddleware(jwtManager *utils.JWTManager, responseHandler *helpers.ResponseHandler, log *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		JWTManager:      jwtManager,
		ResponseHandler: responseHandler,
		Log:             log,
	}
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := am.getTokenFromCookie(r)
		if err != nil {
			am.Log.Warn("Failed to get token from cookie",
				zap.String("url", r.URL.String()),
				zap.String("method", r.Method),
				zap.Error(err))
			am.ResponseHandler.AuthenticationRequiredError(w, r)
			return
		}

		claims, err := am.JWTManager.ValidateToken(token)
		if err != nil {
			am.Log.Warn("Invalid JWT token",
				zap.String("url", r.URL.String()),
				zap.String("method", r.Method),
				zap.Error(err))
			am.ResponseHandler.CustomError(w, r, http.StatusUnauthorized, "Invalid or expired token", err)
			return
		}

		userID, err := uuid.Parse(claims.ID)
		if err != nil {
			am.ResponseHandler.InvalidIDError(w, r, err)
			return
		}

		userInfo := &UserInfo{
			ID:   userID,
			Role: claims.Role,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) RequireRole(requiredRole models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userInfo, ok := r.Context().Value(UserContextKey).(*UserInfo)
			if !ok {
				am.ResponseHandler.AuthenticationRequiredError(w, r)
				return
			}

			if userInfo.Role != requiredRole {
				am.Log.Warn("Insufficient permissions",
					zap.String("url", r.URL.String()),
					zap.String("method", r.Method),
					zap.String("user_role", string(userInfo.Role)),
					zap.String("required_role", string(requiredRole)))
				am.ResponseHandler.InsufficientPermissionsError(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (am *AuthMiddleware) RequireAnyRole(requiredRoles ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userInfo, ok := r.Context().Value(UserContextKey).(*UserInfo)
			if !ok {
				am.ResponseHandler.AuthenticationRequiredError(w, r)
				return
			}

			hasRole := slices.Contains(requiredRoles, userInfo.Role)

			if !hasRole {
				am.Log.Warn("Insufficient permissions",
					zap.String("url", r.URL.String()),
					zap.String("method", r.Method),
					zap.String("user_role", string(userInfo.Role)),
					zap.String("required_roles", strings.Join(roleStrings(requiredRoles), ", ")))
				am.ResponseHandler.InsufficientPermissionsError(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (am *AuthMiddleware) getTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("otterly_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func roleStrings(roles []models.UserRole) []string {
	result := make([]string, len(roles))
	for i, role := range roles {
		result[i] = string(role)
	}
	return result
}

func GetUserFromContext(ctx context.Context) (*UserInfo, bool) {
	userInfo, ok := ctx.Value(UserContextKey).(*UserInfo)
	return userInfo, ok
}
