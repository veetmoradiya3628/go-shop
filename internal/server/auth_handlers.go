package server

import (
	"github.com/gin-gonic/gin"
	"github.com/veetmoradiya3628/go-shop/internal/dto"
	"github.com/veetmoradiya3628/go-shop/internal/services"
	"github.com/veetmoradiya3628/go-shop/internal/utils"
)

func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	authService := services.NewAuthService(s.db, s.config)
	response, err := authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to register user", err)
		return
	}
	utils.CreatedResponse(c, "User registered successfully", response)
}

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := services.NewAuthService(s.db, s.config)
	response, err := authService.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid email or password")
		return
	}
	utils.SuccessResponse(c, "Login successful", response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := services.NewAuthService(s.db, s.config)
	response, err := authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid refresh token")
		return
	}
	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	authService := services.NewAuthService(s.db, s.config)
	if err := authService.Logout(req.RefreshToken); err != nil {
		utils.BadRequestResponse(c, "Failed to logout", err)
		return
	}
	utils.SuccessResponse(c, "Logout successful", nil)
}

func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	userService := services.NewUserService(s.db)
	response, err := userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "Failed to get user profile")
		return
	}
	utils.SuccessResponse(c, "User profile retrieved successfully", response)
}

func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	userService := services.NewUserService(s.db)
	response, err := userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to update user profile", err)
		return
	}
	utils.SuccessResponse(c, "User profile updated successfully", response)
}
