package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
)

// SignedDetails holds details that will be signed into the token.
type SignedDetails struct {
	ID   string
	User string
	jwt.StandardClaims
}

// AppCryptoInterface is a contract for application crypto methods.
type AppCryptoInterface interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword string, providedPassword string) error
	GenerateTokens(id string, user string) (signedToken string, signedRefreshToken string, err error)
	ValidateToken(signedToken string) (*SignedDetails, error)
}

// AppCrypto holds application crypto methods.
type AppCrypto struct {
	config config.Security
}

// NewAppCrypto instantinates AppCrypto.
func NewAppCrypto(config config.Security) AppCryptoInterface {
	return AppCrypto{
		config: config,
	}
}

// HashPassword encrypts the password.
func (c AppCrypto) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func (c AppCrypto) VerifyPassword(hashedPassword string, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
}

// GenerateTokens generates JWT tokens.
func (c AppCrypto) GenerateTokens(id string, user string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		ID:   id,
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().
				Add(time.Hour * time.Duration(c.config.Jwt.TokenExpiration)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().
				Add(time.Hour * time.Duration(c.config.Jwt.RefreshTokenExpiration)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(c.config.Jwt.SecretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(c.config.Jwt.SecretKey))

	return token, refreshToken, errors.Wrap(err, "token sign failed")
}

// ValidateToken validates the JWT token.
func (c AppCrypto) ValidateToken(signedToken string) (*SignedDetails, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(c.config.Jwt.SecretKey), nil
	}

	token, tokenErr := jwt.ParseWithClaims(signedToken, &SignedDetails{}, keyFunc)
	if tokenErr != nil {
		return nil, errors.Wrap(tokenErr, "invalid token")
	}

	payload, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return payload, nil
}
