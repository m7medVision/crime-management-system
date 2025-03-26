package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/auth"
	"github.com/m7medVision/crime-management-system/internal/dto"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Admin endpoint to create a new user with specified role and clearance
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserDTO true "User details"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string "Invalid user data"
// @Failure 403 {object} map[string]string "Permission denied"
// @Security ApiKeyAuth
// @Router /users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var userDTO dto.UserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password if provided
	hashedPassword, err := auth.HashPassword(userDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user := &model.User{
		Username:       userDTO.Username,
		Password:       hashedPassword,
		Email:          userDTO.Email,
		FullName:       userDTO.FullName,
		Role:           model.Role(userDTO.Role),
		ClearanceLevel: model.ClearanceLevel(userDTO.ClearanceLevel),
		IsActive:       userDTO.IsActive,
	}

	result, err := ctrl.userService.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Don't return the password hash in the response
	result.Password = ""
	c.JSON(http.StatusCreated, result)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Admin endpoint to update user details including role and clearance
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UpdateUserDTO true "Updated user details"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string "Invalid user data"
// @Failure 403 {object} map[string]string "Permission denied"
// @Failure 404 {object} map[string]string "User not found"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateDTO dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing user
	existingUser, err := ctrl.userService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields if provided
	if updateDTO.Email != "" {
		existingUser.Email = updateDTO.Email
	}
	if updateDTO.FullName != "" {
		existingUser.FullName = updateDTO.FullName
	}
	if updateDTO.Password != "" {
		hashedPassword, err := auth.HashPassword(updateDTO.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
		existingUser.Password = hashedPassword
	}
	if updateDTO.Role != "" {
		existingUser.Role = model.Role(updateDTO.Role)
	}
	if updateDTO.ClearanceLevel != "" {
		existingUser.ClearanceLevel = model.ClearanceLevel(updateDTO.ClearanceLevel)
	}
	if updateDTO.IsActive != nil {
		existingUser.IsActive = *updateDTO.IsActive
	}

	result, err := ctrl.userService.UpdateUser(existingUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Don't return the password hash in the response
	result.Password = ""
	c.JSON(http.StatusOK, result)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Admin endpoint to delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 403 {object} map[string]string "Permission denied"
// @Failure 404 {object} map[string]string "User not found"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := ctrl.userService.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListUsers godoc
// @Summary List all users
// @Description Admin endpoint to get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param offset query int false "Pagination offset" default(0)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "users and total count"
// @Failure 403 {object} map[string]string "Permission denied"
// @Security ApiKeyAuth
// @Router /users [get]
func (ctrl *UserController) ListUsers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := ctrl.userService.ListUsers(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't return password hashes in the response
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
	})
}

// GetUser godoc
// @Summary Get user details
// @Description Admin endpoint to get details of a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 403 {object} map[string]string "Permission denied"
// @Failure 404 {object} map[string]string "User not found"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctrl.userService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't return the password hash in the response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}
