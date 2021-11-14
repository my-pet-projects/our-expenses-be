package command_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestSignUpHandler_ReturnsHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	crypto := new(mocks.AppCryptoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := command.NewSignUpHandler(repo, crypto, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestSignUpHandler_GetExistingUserError_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	crypto := new(mocks.AppCryptoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()

	cmd := command.SignUpCommand{
		Username: "user",
		Password: "123",
	}

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewSignUpHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestSignUpHandler_HasExistingUser_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	crypto := new(mocks.AppCryptoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()

	cmd := command.SignUpCommand{
		Username: "user",
		Password: "123",
	}

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(&domain.User{}, nil)

	// SUT
	sut := command.NewSignUpHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestSignUpHandler_HashPasswordFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.SignUpCommand{
		Username: "user",
		Password: "123",
	}

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(nil, nil)

	matchHashFn := func(password string) bool {
		return password == cmd.Password
	}
	crypto.On("HashPassword", mock.MatchedBy(matchHashFn)).Return("", errors.New("error"))

	// SUT
	sut := command.NewSignUpHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestSignUpHandler_GenerateTokensFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.SignUpCommand{
		Username: "user",
		Password: "123",
	}
	hash := "hash1"

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(nil, nil)

	matchHashFn := func(password string) bool {
		return password == cmd.Password
	}
	crypto.On("HashPassword", mock.MatchedBy(matchHashFn)).Return(hash, nil)
	matchUsernameFn := func(username string) bool {
		return username == cmd.Username
	}
	crypto.On("GenerateTokens", mock.Anything, mock.MatchedBy(matchUsernameFn)).Return("", "", errors.New("error"))

	// SUT
	sut := command.NewSignUpHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestSignUpHandler_CreateUserFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.SignUpCommand{
		Username: "",
		Password: "123",
	}
	hash := "hash2"
	token := "token2"
	refreshToken := "refreshToken2"

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(nil, nil)

	matchHashFn := func(password string) bool {
		return password == cmd.Password
	}
	crypto.On("HashPassword", mock.MatchedBy(matchHashFn)).Return(hash, nil)
	matchUsernameFn := func(username string) bool {
		return username == cmd.Username
	}
	crypto.On("GenerateTokens", mock.Anything, mock.MatchedBy(matchUsernameFn)).Return(token, refreshToken, nil)

	// SUT
	sut := command.NewSignUpHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestSignUpHandler_InsertUserFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.SignUpCommand{
		Username: "username",
		Password: "123",
	}
	hash := "hash3"
	token := "token3"
	refreshToken := "refreshToken3"

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(nil, nil)

	matchHashFn := func(password string) bool {
		return password == cmd.Password
	}
	crypto.On("HashPassword", mock.MatchedBy(matchHashFn)).Return(hash, nil)
	matchUsernameFn := func(username string) bool {
		return username == cmd.Username
	}
	crypto.On("GenerateTokens", mock.Anything, mock.MatchedBy(matchUsernameFn)).Return(token, refreshToken, nil)

	matchUserFn := func(user *domain.User) bool {
		return user.Username() == cmd.Username && user.Password() == hash &&
			user.Token() == token && user.RefreshToken() == refreshToken
	}
	repo.On("Insert", mock.Anything,
		mock.MatchedBy(matchUserFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewSignUpHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestSignUpHandler_HappyPath_DoesNotThrowError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.SignUpCommand{
		Username: "username",
		Password: "123",
	}
	hash := "hash4"
	token := "token4"
	refreshToken := "refreshToken4"

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(nil, nil)

	matchHashFn := func(password string) bool {
		return password == cmd.Password
	}
	crypto.On("HashPassword", mock.MatchedBy(matchHashFn)).Return(hash, nil)
	matchUsernameFn := func(username string) bool {
		return username == cmd.Username
	}
	crypto.On("GenerateTokens", mock.Anything, mock.MatchedBy(matchUsernameFn)).Return(token, refreshToken, nil)

	matchUserFn := func(user *domain.User) bool {
		return user.Username() == cmd.Username && user.Password() == hash &&
			user.Token() == token && user.RefreshToken() == refreshToken
	}
	repo.On("Insert", mock.Anything,
		mock.MatchedBy(matchUserFn)).Return(&domain.InsertResult{}, nil)

	// SUT
	sut := command.NewSignUpHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, result, "Result should not be nil.")
	assert.Nil(t, err, "Error result should be nil.")
	assert.Equal(t, cmd.Username, result.Username())
	assert.Equal(t, hash, result.Password())
	assert.Equal(t, token, result.Token())
	assert.Equal(t, refreshToken, result.RefreshToken())
}
