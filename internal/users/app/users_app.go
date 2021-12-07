package app

import (
	"context"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/auth"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// Application provides an application.
type Application struct {
	Commands Commands
	Queries  Queries
	Logger   logger.LogInterface
	Config   config.Config
	Tracer   tracer.TraceInterface
}

// Commands struct holds available application commands.
type Commands struct {
	SignUp command.SignUpHandlerInterface
	Login  command.LoginHandlerInterface
}

// Queries struct holds available application queries.
type Queries struct{}

// NewApplication returns application instance.
func NewApplication(
	ctx context.Context,
	cancel context.CancelFunc,
	config *config.Config,
	logger *logger.Logger,
	tracer *tracer.Tracer,
	mongoClient *database.MongoClient,
) (*Application, error) {
	crypto := auth.NewAppCrypto(config.Server.Security)
	userRepo := adapters.NewUserRepo(mongoClient, logger)

	return &Application{
		Commands: Commands{
			SignUp: command.NewSignUpHandler(userRepo, crypto, logger),
			Login:  command.NewLoginHandler(userRepo, crypto, logger),
		},
		Queries: Queries{},
		Logger:  logger,
		Config:  *config,
		Tracer:  tracer,
	}, nil
}
