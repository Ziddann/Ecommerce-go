package services

import (
    "errors"
    "strings"
    "time"

    "ecommerce/dto"
    "ecommerce/models"
    "ecommerce/repositories"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

var (
    ErrEmailAlreadyRegistered = errors.New("email already registered")
    ErrInvalidCredentials     = errors.New("invalid email or password")
    ErrUserNotFound           = errors.New("user not found")
    ErrPasswordTooLong        = errors.New("password exceeds 72 bytes")
)

type JWTClaims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

type AuthService struct {
    UserRepository *repositories.UserRepository
    JWTSecret      string
    ExpiresInHours int
}

func NewAuthService(
    userRepository *repositories.UserRepository,
    jwtSecret string,
    expiresInHours int,
) *AuthService {
    return &AuthService{
        UserRepository: userRepository,
        JWTSecret:      jwtSecret,
        ExpiresInHours: expiresInHours,
    }
}

func (s *AuthService) Register(
    request dto.RegisterRequest,
) (*models.User, error) {

    email := strings.ToLower(strings.TrimSpace(request.Email))

    if len([]byte(request.Password)) > 72 {
        return nil, ErrPasswordTooLong
    }

    existingUser, err := s.UserRepository.GetByEmail(email)

    if err == nil && existingUser != nil {
        return nil, ErrEmailAlreadyRegistered
    }

    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }

    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(request.Password),
        bcrypt.DefaultCost,
    )

    if err != nil {
        return nil, err
    }

    user := &models.User{
        Name:     strings.TrimSpace(request.Name),
        Email:    email,
        Password: string(hashedPassword),
        Role:     "customer",
    }

    if err := s.UserRepository.Create(user); err != nil {
        if errors.Is(err, gorm.ErrDuplicatedKey) {
            return nil, ErrEmailAlreadyRegistered
        }

        return nil, err
    }

    return user, nil
}

func (s *AuthService) Login(
    email string,
    password string,
) (string, *models.User, error) {

    email = strings.ToLower(strings.TrimSpace(email))

    if len([]byte(password)) > 72 {
        return "", nil, ErrInvalidCredentials
    }

    user, err := s.UserRepository.GetByEmail(email)

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return "", nil, ErrInvalidCredentials
    }

    if err != nil {
        return "", nil, err
    }

    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password),
        []byte(password),
    )

    if err != nil {
        return "", nil, ErrInvalidCredentials
    }

    token, err := s.GenerateToken(user)

    if err != nil {
        return "", nil, err
    }

    return token, user, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {

    claims := JWTClaims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(
                time.Now().Add(
                    time.Duration(s.ExpiresInHours) * time.Hour,
                ),
            ),
            IssuedAt: jwt.NewNumericDate(time.Now()),
            Issuer:   "ecommerce-api",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString([]byte(s.JWTSecret))
}

func (s *AuthService) GetProfile(id uint) (*models.User, error) {
    user, err := s.UserRepository.GetByID(id)

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }

    if err != nil {
        return nil, err
    }

    return user, nil
}