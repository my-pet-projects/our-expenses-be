package command

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/auth"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// SignUpCommand defines a sign-up command.
type SignUpCommand struct {
	Username string
	Password string
}

// SignUpHandler defines a handler to sign-up a user.
type SignUpHandler struct {
	repo   adapters.UserRepoInterface
	crypto auth.AppCryptoInterface
	logger logger.LogInterface
}

// SignUpHandlerInterface defines a contract to handle command.
type SignUpHandlerInterface interface {
	Handle(ctx context.Context, cmd SignUpCommand) (*domain.User, error)
}

// NewSignUpHandler returns command handler.
func NewSignUpHandler(
	repo adapters.UserRepoInterface,
	crypto auth.AppCryptoInterface,
	logger logger.LogInterface,
) SignUpHandler {
	return SignUpHandler{
		repo:   repo,
		crypto: crypto,
		logger: logger,
	}
}

// Handle handles sign-up command.
func (h SignUpHandler) Handle(ctx context.Context, cmd SignUpCommand) (*domain.User, error) {
	ctx, span := tracer.NewSpan(ctx, "execute sign-up command")
	defer span.End()

	existingUser, existingUserErr := h.repo.GetOne(ctx, cmd.Username)
	if existingUserErr != nil {
		tracer.AddSpanError(span, existingUserErr)
		return nil, errors.Wrap(existingUserErr, "fetch user failed")
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	id := primitive.NewObjectID()

	_, hashSpan := tracer.NewSpan(ctx, "generate password hash")
	pwd, pwdErr := h.crypto.HashPassword(cmd.Password)
	if pwdErr != nil {
		tracer.AddSpanError(span, pwdErr)
		return nil, errors.Wrap(pwdErr, "password hash failed")
	}
	hashSpan.End()

	_, tokenSpan := tracer.NewSpan(ctx, "generate jwt tokens")
	token, refreshToken, tokenErr := h.crypto.GenerateTokens(id.Hex(), cmd.Username)
	if tokenErr != nil {
		tracer.AddSpanError(span, tokenErr)
		return nil, errors.Wrap(tokenErr, "token generation failed")
	}
	tokenSpan.End()

	user, userErr := domain.NewUser(id.Hex(), cmd.Username, pwd, token, refreshToken)
	if userErr != nil {
		tracer.AddSpanError(span, userErr)
		return nil, errors.Wrap(userErr, "user creation failed")
	}

	_, userIDErr := h.repo.Insert(ctx, user)
	if userIDErr != nil {
		tracer.AddSpanError(span, userIDErr)
		return nil, errors.Wrap(userIDErr, "user save failed")
	}

	return user, nil
}
