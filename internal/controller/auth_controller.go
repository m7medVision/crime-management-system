package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/dto"
	"github.com/m7medVision/crime-management-system/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginDTO true "Login Credentials"
// @Success 200 {object} map[string]interface{} "token and user information"
// @Failure 400 {object} map[string]string "Invalid login request"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var loginDTO dto.LoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		errorResponse := dto.ErrorDTO{
			Message: "Invalid login request",
			Code:    http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	token, user, err := ctrl.authService.Login(loginDTO.Username, loginDTO.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorDTO{
			Message: "Invalid credentials",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponseDTO{
		Token: token,
		User:  user.Email,
	})
}
