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

func TestLoginHandler_ReturnsHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	crypto := new(mocks.AppCryptoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := command.NewLoginHandler(repo, crypto, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestLoginHandler_GetExistingUserError_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	crypto := new(mocks.AppCryptoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()

	cmd := command.LoginCommand{
		Username: "user",
		Password: "123",
	}

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewLoginHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestLoginHandler_HasExistingUser_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	crypto := new(mocks.AppCryptoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()

	cmd := command.LoginCommand{
		Username: "user",
		Password: "123",
	}

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(nil, nil)

	// SUT
	sut := command.NewLoginHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestLoginHandler_VerifyPasswordFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.LoginCommand{
		Username: "user",
		Password: "123",
	}
	user, _ := domain.NewUser("id", cmd.Username, cmd.Password,
		"token", "refresh")

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(user, nil)
	crypto.On("VerifyPassword", user.Password(), cmd.Password).Return(errors.New("error"))

	// SUT
	sut := command.NewLoginHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestLoginHandler_GenerateTokensFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.LoginCommand{
		Username: "user",
		Password: "123",
	}
	user, _ := domain.NewUser("id", cmd.Username, cmd.Password,
		"token", "refresh")

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(user, nil)
	crypto.On("VerifyPassword", user.Password(), cmd.Password).Return(nil)
	matchUsernameFn := func(username string) bool {
		return username == cmd.Username
	}
	crypto.On("GenerateTokens", mock.Anything, mock.MatchedBy(matchUsernameFn)).Return("", "", errors.New("error"))

	// SUT
	sut := command.NewLoginHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestLoginHandler_UpdateUserFails_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.LoginCommand{
		Username: "username",
		Password: "123",
	}
	hash := "hash3"
	token := "token3"
	refreshToken := "refreshToken3"
	existingUser, _ := domain.NewUser("id", cmd.Username, hash,
		"oldtoken", "oldrefreshtoken")

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(existingUser, nil)
	crypto.On("VerifyPassword", existingUser.Password(), cmd.Password).Return(nil)
	matchUsernameFn := func(username string) bool {
		return username == cmd.Username
	}
	crypto.On("GenerateTokens", mock.Anything, mock.MatchedBy(matchUsernameFn)).Return(token, refreshToken, nil)
	matchUserFn := func(user *domain.User) bool {
		return user.Username() == existingUser.Username() && user.Password() == existingUser.Password() &&
			user.Token() == token && user.RefreshToken() == refreshToken
	}
	repo.On("Update", mock.Anything,
		mock.MatchedBy(matchUserFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewLoginHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestLoginHandler_HappyPath_DoesNotThrowError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.UserRepoInterface)
	log := new(mocks.LogInterface)
	crypto := new(mocks.AppCryptoInterface)
	ctx := context.Background()

	cmd := command.LoginCommand{
		Username: "username",
		Password: "123",
	}
	hash := "hash4"
	token := "token4"
	refreshToken := "refreshToken4"
	existingUser, _ := domain.NewUser("id", cmd.Username, hash,
		"oldtoken", "oldrefreshtoken")

	matchExistingUserFn := func(username string) bool {
		return username == cmd.Username
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchExistingUserFn)).Return(existingUser, nil)
	crypto.On("VerifyPassword", existingUser.Password(), cmd.Password).Return(nil)
	matchUsernameFn := func(username string) bool {
		return username == cmd.Username
	}
	crypto.On("GenerateTokens", mock.Anything, mock.MatchedBy(matchUsernameFn)).Return(token, refreshToken, nil)
	matchUserFn := func(user *domain.User) bool {
		return user.Username() == existingUser.Username() && user.Password() == existingUser.Password() &&
			user.Token() == token && user.RefreshToken() == refreshToken
	}
	repo.On("Update", mock.Anything,
		mock.MatchedBy(matchUserFn)).Return(&domain.UpdateResult{}, nil)

	// SUT
	sut := command.NewLoginHandler(repo, crypto, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, result, "Result should not be nil.")
	assert.Nil(t, err, "Error result should be nil.")
	assert.Equal(t, existingUser, result)
}
