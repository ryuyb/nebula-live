package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrTokenClaims  = errors.New("invalid token claims")
)

// TokenConfig JWT配置
type TokenConfig struct {
	SecretKey        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	Issuer           string
}

// DefaultTokenConfig 默认JWT配置
var DefaultTokenConfig = &TokenConfig{
	SecretKey:        "your-secret-key-change-this-in-production",
	AccessTokenTTL:   15 * time.Minute,
	RefreshTokenTTL:  7 * 24 * time.Hour, // 7 days
	Issuer:           "nebula-live",
}

// UserClaims 用户JWT声明
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

// JWTManager JWT管理器
type JWTManager struct {
	config *TokenConfig
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(config *TokenConfig) *JWTManager {
	if config == nil {
		config = DefaultTokenConfig
	}
	return &JWTManager{config: config}
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
func (j *JWTManager) GenerateTokenPair(userID uint, username, email string) (*TokenPair, error) {
	now := time.Now()
	
	// 生成访问令牌
	accessToken, err := j.generateToken(userID, username, email, now.Add(j.config.AccessTokenTTL))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err := j.generateToken(userID, username, email, now.Add(j.config.RefreshTokenTTL))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(j.config.AccessTokenTTL).Unix(),
		TokenType:    "Bearer",
	}, nil
}

// generateToken 生成JWT令牌
func (j *JWTManager) generateToken(userID uint, username, email string, expiresAt time.Time) (string, error) {
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.config.Issuer,
			Subject:   fmt.Sprintf("user_%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.SecretKey))
}

// ValidateToken 验证JWT令牌
func (j *JWTManager) ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenClaims
	}

	return claims, nil
}

// RefreshToken 刷新访问令牌
func (j *JWTManager) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 生成新的令牌对
	return j.GenerateTokenPair(claims.UserID, claims.Username, claims.Email)
}

// ExtractUserID 从令牌中提取用户ID
func (j *JWTManager) ExtractUserID(tokenString string) (uint, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}