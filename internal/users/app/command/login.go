package command

import (
	"context"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/auth"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// LoginCommand defines a login command.
type LoginCommand struct {
	Username string
	Password string
}

// LoginHandler defines a handler to login a user.
type LoginHandler struct {
	repo   adapters.UserRepoInterface
	crypto auth.AppCryptoInterface
	logger logger.LogInterface
}

// LoginHandlerInterface defines a contract to handle command.
type LoginHandlerInterface interface {
	Handle(ctx context.Context, cmd LoginCommand) (*domain.User, error)
}

// NewLoginHandler returns command handler.
func NewLoginHandler(
	repo adapters.UserRepoInterface,
	crypto auth.AppCryptoInterface,
	logger logger.LogInterface,
) LoginHandler {
	return LoginHandler{
		repo:   repo,
		crypto: crypto,
		logger: logger,
	}
}

// Handle handles login command.
func (h LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (*domain.User, error) {
	ctx, span := tracer.NewSpan(ctx, "execute login command")
	defer span.End()

	user, userErr := h.repo.GetOne(ctx, cmd.Username)
	if userErr != nil {
		tracer.AddSpanError(span, userErr)
		return nil, errors.Wrap(userErr, "fetch user failed")
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	validPassErr := h.crypto.VerifyPassword(user.Password(), cmd.Password)
	if validPassErr != nil {
		tracer.AddSpanError(span, validPassErr)
		return nil, domain.ErrWrongPassword
	}

	token, refreshToken, tokenErr := h.crypto.GenerateTokens(user.ID(), user.Username())
	if tokenErr != nil {
		tracer.AddSpanError(span, tokenErr)
		return nil, errors.Wrap(tokenErr, "token generation failed")
	}

	user.UpdateTokens(token, refreshToken)

	_, updateErr := h.repo.Update(ctx, user)
	if updateErr != nil {
		tracer.AddSpanError(span, updateErr)
		return nil, errors.Wrap(updateErr, "user update failed")
	}

	return user, nil
}
