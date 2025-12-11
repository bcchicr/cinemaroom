package app

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/config"
	"github.com/bcchicr/cinemaroom/services/authn/internal/application/handlers"
	grpcendpoint "github.com/bcchicr/cinemaroom/services/authn/internal/endpoint/grpc"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure/gateway"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure/persistence"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure/services/password"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure/services/token"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Container struct {
	GrpcServer *grpc.Server
	DB         *sql.DB
}

func NewContainer(cfg *config.Config) (*Container, error) {
	db, err := sql.Open("pgx", cfg.Postgres.Url)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatalf("Failed to ping database: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}

	accountRepository := persistence.NewPostgresAccountRepository(db)
	passwordService := password.NewService()

	refreshTokenRepository := persistence.NewPostgresRefreshTokenRepository(db)
	tokenService := token.NewService(
		cfg.Jwt.AccessSecret,
		time.Duration(cfg.Jwt.AccessTTLInMinutes)*time.Minute,
		time.Duration(cfg.Jwt.RefreshTTLInDays)*time.Hour*24,
		refreshTokenRepository,
		redisClient,
	)
	userServiceClient, err := gateway.NewGrpcUserServiceClient(cfg.UserService.Url)
	if err != nil {
		log.Fatalf("Failed to create user service client")
	}

	registrationHandler := handlers.NewRegistrationHandler(
		passwordService,
		tokenService,
		accountRepository,
		refreshTokenRepository,
		userServiceClient,
	)

	loginHandler := handlers.NewLoginHandler(
		tokenService,
		accountRepository,
		refreshTokenRepository,
	)

	logoutHandler := handlers.NewLogoutHandler(
		tokenService,
		accountRepository,
		refreshTokenRepository,
		redisClient,
	)

	refreshHandler := handlers.NewRefreshHandler(
		tokenService,
		accountRepository,
		refreshTokenRepository,
	)

	editHandler := handlers.NewEditHandler(
		passwordService,
		tokenService,
		accountRepository,
		refreshTokenRepository,
	)

	errorHandlerInterceptor := grpcendpoint.NewErrorHandlerInterceptor()

	grpcServer := grpcendpoint.NewServer(
		registrationHandler,
		loginHandler,
		logoutHandler,
		refreshHandler,
		editHandler,
		errorHandlerInterceptor,
	)

	return &Container{
		GrpcServer: grpcServer,
		DB:         db,
	}, nil
}

func (c *Container) Close() error {
	c.GrpcServer.Stop()
	return c.DB.Close()
}
