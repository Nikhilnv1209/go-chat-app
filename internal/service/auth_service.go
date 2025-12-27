package service

import (
	"chat-app/internal/errors"
	"chat-app/internal/models"
	"chat-app/internal/repository"
	"chat-app/pkg/jwt"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo   repository.UserRepository
	jwtService jwt.Service
}

func NewAuthService(userRepo repository.UserRepository, jwtService jwt.Service) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Register(username, email, password string) (*models.User, error) {
	// 1. Check if user already exists
	existing, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.ErrEmailExists
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	// 1. Find user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.ErrInvalidCredentials
	}

	// 2. Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.ErrInvalidCredentials
	}

	// 3. Generate Token
	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) ValidateToken(tokenString string) (uuid.UUID, error) {
	return s.jwtService.ValidateToken(tokenString)
}

func (s *authService) SearchUsers(query string, excludeUserID uuid.UUID) ([]models.User, error) {
	return s.userRepo.Search(query, excludeUserID)
}

func (s *authService) GetUser(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
