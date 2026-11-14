package user

import (
	"errors"
	"gotickets/internal/domain/user/dto"
	"gotickets/internal/httpresponse"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service}
}

func (h *handler) CreateUser(c *echo.Context) error {
	var req dto.CreateUserRequest // input

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.CreateUser(req)
	if err != nil {

		if errors.Is(err, ErrorAlreadyExists) {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code:    http.StatusConflict,
				Message: "Failed to create User",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest // input

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{
				Code:    http.StatusUnauthorized,
				Message: "Invalid email or password",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to login user",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) LoginWithGoogle(c *echo.Context) error {
	var req dto.GoogleLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Code: http.StatusBadRequest, Message: "Invalid request payload", Details: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Code: http.StatusBadRequest, Message: "Validation failed", Details: err.Error()})
	}
	response, err := h.service.LoginWithGoogle(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{Code: http.StatusUnauthorized, Message: "Invalid Google credentials"})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{Code: http.StatusInternalServerError, Message: "Google login failed", Details: err.Error()})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *handler) UpdateRole(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Code: http.StatusBadRequest, Message: "Invalid user id"})
	}
	var req dto.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Code: http.StatusBadRequest, Message: "Invalid request payload", Details: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{Code: http.StatusBadRequest, Message: "Validation failed", Details: err.Error()})
	}
	currentRole, _ := c.Get("user_role").(string)
	if currentRole == RoleAdmin && (req.Role == RoleAdmin || req.Role == RoleSuperAdmin) {
		return c.JSON(http.StatusForbidden, httpresponse.Error{Code: http.StatusForbidden, Message: "Admins cannot assign this role"})
	}
	if err := h.service.UpdateRole(uint(id), req.Role); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{Code: http.StatusNotFound, Message: "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{Code: http.StatusInternalServerError, Message: "Failed to update role", Details: err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Role updated"})
}

func (h *handler) GetMe(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Cannot get user information",
			Details: "missing user id in context",
		})
	}

	email, _ := c.Get("user_email").(string)
	name, _ := c.Get("user_name").(string)
	role, _ := c.Get("user_role").(string)

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:    userID,
		Name:  name,
		Email: email,
		Role:  role,
	})
}
