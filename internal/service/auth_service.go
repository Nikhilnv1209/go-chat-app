package service

import (
	"context"
	apperrors "chat-app/internal/errors"
	"chat-app/internal/models"
	"chat-app/internal/repository"
	"chat-app/pkg/jwt"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtService       jwt.Service
}

func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, jwtService jwt.Service) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password string) (string, string, *models.User, error) {
	// 1. Check if user already exists
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", nil, err
	}
	if existing != nil {
		return "", "", nil, apperrors.ErrEmailExists
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", nil, err
	}

	// 3. Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		LastSeen: time.Now(),
		IsOnline: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", "", nil, err
	}

	// 4. Generate Tokens
	accessToken, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := s.createRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, *models.User, error) {
	// 1. Find user
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", nil, err
	}
	if user == nil {
		return "", "", nil, apperrors.ErrNotFound
	}

	// 2. Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, apperrors.ErrInvalidCredentials
	}

	// 3. Generate Tokens
	accessToken, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := s.createRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	// 1. Hash the incoming token
	hash := hashToken(refreshToken)

	// 2. Find in DB
	storedToken, err := s.refreshTokenRepo.GetByHash(ctx, hash)
	if err != nil {
		// Could check if err is "record not found" -> return ErrInvalidToken
		return "", "", apperrors.ErrInvalidCredentials // Or a specific ErrInvalidToken
	}

	// 3. Validate
	if storedToken.Revoked {
		// Security alert: Someone tried to use a revoked token!
		// We could revoke ALL user tokens here for safety.
		return "", "", errors.New("token revoked")
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return "", "", errors.New("token expired")
	}

	// 4. Generate new Access Token
	accessToken, err := s.jwtService.GenerateToken(storedToken.UserID)
	if err != nil {
		return "", "", err
	}

	// 5. Rotate Refresh Token (Optional but good)
	// For this MVP, let's KEEP the same refresh token to avoid cookie complexity
	// OR we can rotate it. Let's start with basic: Reuse same refresh token until expiry.
	// If we want rotation, we'd revoke old and create new.
	// Let's just return the SAME refresh token for now, or "" if not changing.
	// The interface returns (string, string, error). Let's return new access, SAME refresh.

	return accessToken, refreshToken, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	hash := hashToken(refreshToken)
	storedToken, err := s.refreshTokenRepo.GetByHash(ctx, hash)
	if err != nil {
		return nil // Already gone or invalid, just ignore
	}
	return s.refreshTokenRepo.Revoke(ctx, storedToken.ID)
}

func (s *authService) ValidateToken(tokenString string) (uuid.UUID, error) {
	return s.jwtService.ValidateToken(tokenString)
}

func (s *authService) SearchUsers(ctx context.Context, query string, excludeUserID uuid.UUID) ([]models.User, error) {
	return s.userRepo.Search(ctx, query, excludeUserID)
}

func (s *authService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// Helpers

func (s *authService) createRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	// Generate 32 bytes of random entropy
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	rawToken := base64.URLEncoding.EncodeToString(b)

	// Hash it for storage
	hash := hashToken(rawToken)

	token := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		Revoked:   false,
	}

	if err := s.refreshTokenRepo.Create(ctx, token); err != nil {
		return "", err
	}

	return rawToken, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
