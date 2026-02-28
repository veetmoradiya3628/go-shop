package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/veetmoradiya3628/go-shop/internal/config"
	"github.com/veetmoradiya3628/go-shop/internal/dto"
	"github.com/veetmoradiya3628/go-shop/internal/models"
	"github.com/veetmoradiya3628/go-shop/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// check if user exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("User already exists")
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// create user
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	// create a cart
	cart := models.Cart{
		UserID: user.ID,
	}
	if err := s.db.Create(&cart).Error; err != nil {
		fmt.Println("Unable to create cart")
		return nil, err
	}

	// generate token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// find user by email
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("Invalid email or password")
	}
	// check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("Invalid email or password")
	}
	// generate token
	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("Invalid refresh token")
	}
	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("Refresh token not found")
	}

	var user models.User
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}
	s.db.Delete(&refreshToken)
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
}
func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(&s.config.JWT, user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}
	// save refresh token
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}
	if err := s.db.Create(&refreshTokenModel).Error; err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
	}, nil
}
