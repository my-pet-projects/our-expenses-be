package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
)

func TestNewAppCrypto_ReturnsAppCrypto(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Security{
		Jwt: config.Jwt{
			SecretKey: "key",
		},
	}

	// Act
	result := NewAppCrypto(config)

	// Assert
	assert.NotNil(t, result, "Result should not be nil.")
}

func TestHashPassword_ReturnsHash(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Security{
		Jwt: config.Jwt{
			SecretKey: "key",
		},
	}

	// SUT
	sut := NewAppCrypto(config)

	// Act
	res, resErr := sut.HashPassword("hash")

	// Assert
	assert.NotNil(t, res, "Result should not be nil.")
	assert.Nil(t, resErr, "Result error should be nil.")
}

func TestVerifyPassword_ReturnsVerificationResult(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Security{
		Jwt: config.Jwt{
			SecretKey: "key",
		},
	}

	// SUT
	sut := NewAppCrypto(config)

	// Act
	resErr := sut.VerifyPassword("hash", "plain")

	// Assert
	assert.NotNil(t, resErr, "Result error should not be nil.")
}

func TestGenerateTokens_ReturnsVerificationResult(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Security{
		Jwt: config.Jwt{
			SecretKey: "key",
		},
	}

	// SUT
	sut := NewAppCrypto(config)

	// Act
	token, refreshToken, err := sut.GenerateTokens("id", "username")

	// Assert
	assert.Nil(t, err, "Result error should be nil.")
	assert.NotNil(t, token)
	assert.NotNil(t, refreshToken)
}

func TestValidateToken_InvalidToken_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Security{
		Jwt: config.Jwt{
			SecretKey: "key",
		},
	}

	// SUT
	sut := NewAppCrypto(config)

	// Act
	res, resErr := sut.ValidateToken("token")

	// Assert
	assert.NotNil(t, resErr, "Result error should not be nil.")
	assert.Nil(t, res)
}

func TestValidateToken_ExpiredToken_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Security{
		Jwt: config.Jwt{
			SecretKey:       "key",
			TokenExpiration: -1,
		},
	}

	// SUT
	sut := NewAppCrypto(config)
	token, _, _ := sut.GenerateTokens("id", "string")

	// Act
	res, resErr := sut.ValidateToken(token)

	// Assert
	assert.NotNil(t, resErr, "Result error should not be nil.")
	assert.Nil(t, res)
}

func TestValidateToken_ValidToken_ReturnClaims(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Security{
		Jwt: config.Jwt{
			SecretKey:       "key",
			TokenExpiration: 1,
		},
	}

	// SUT
	sut := NewAppCrypto(config)
	token, _, _ := sut.GenerateTokens("id", "string")

	// Act
	res, resErr := sut.ValidateToken(token)

	// Assert
	assert.Nil(t, resErr, "Result error should not be nil.")
	assert.NotNil(t, res)
}
