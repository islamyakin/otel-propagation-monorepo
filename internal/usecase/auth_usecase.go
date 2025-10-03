package usecase

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/islamyakin/otel-propagation-monorepo/internal/config"
	"github.com/islamyakin/otel-propagation-monorepo/internal/entity"
	"github.com/islamyakin/otel-propagation-monorepo/internal/model"
	"github.com/islamyakin/otel-propagation-monorepo/internal/repository"
)

type AuthUseCase interface {
	Register(req *model.RegisterRequest) (*entity.User, error)
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	VerifyToken(tokenString string) (*JWTClaims, error)
}

type JWTClaims struct {
	UserID   int         `json:"user_id"`
	Username string      `json:"username"`
	Role     entity.Role `json:"role"`
	jwt.RegisteredClaims
}

type authUseCase struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthUseCase(userRepo repository.UserRepository, config *config.Config) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
		config:   config,
	}
}

func (uc *authUseCase) Register(req *model.RegisterRequest) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with default role "user"
	user := &entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     entity.UserRole,
	}

	createdUser, err := uc.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Remove password from response
	createdUser.Password = ""
	return createdUser, nil
}

func (uc *authUseCase) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// Get user by username
	user, err := uc.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := uc.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Remove password from response
	user.Password = ""

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (uc *authUseCase) VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (uc *authUseCase) generateToken(user *entity.User) (string, error) {
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    uc.config.JWT.Issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.config.JWT.SecretKey))
}
