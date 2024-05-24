package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arfan21/backend-test/config"
	_ "github.com/arfan21/backend-test/docs"
	"github.com/arfan21/backend-test/internal/middleware"
	"github.com/arfan21/backend-test/pkg/constant"
	"github.com/arfan21/backend-test/pkg/exception"
	"github.com/arfan21/backend-test/pkg/logger"
	"github.com/arfan21/backend-test/pkg/pkgutil"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ctxTimeout = 5
)

type Server struct {
	app     *fiber.App
	dbRedis *redis.Storage
	dbPg    *pgxpool.Pool
}

func New(
	dbRedis *redis.Storage,
	dbPg *pgxpool.Pool,
) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: exception.FiberErrorHandler,
	})

	timeout := time.Duration(config.GetConfig().Service.Timeout) * time.Second

	app.Use(middleware.Timeout(timeout))
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger.Log(context.Background()),
		Fields: append(fiberzerolog.ConfigDefault.Fields, fiberzerolog.FieldRequestID),
	}))

	app.Use(cors.New())
	app.Use(requestid.New())

	app.Use(middleware.RequestIdUser())
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	app.Use(limiter.New(limiter.Config{
		Max:        config.GetConfig().RateLimit.Max,
		Expiration: time.Duration(config.GetConfig().RateLimit.Expiration) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {

			return constant.ErrTooManyRequest
		},
		Storage: dbRedis,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	return &Server{
		app:     app,
		dbRedis: dbRedis,
		dbPg:    dbPg,
	}
}

func (s *Server) Run() error {
	s.Routes()
	ctx := context.Background()
	go func() {
		if err := s.app.Listen(pkgutil.GetPort()); err != nil {
			logger.Log(ctx).Fatal().Err(err).Msg("failed to start server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	_, shutdown := context.WithTimeout(ctx, ctxTimeout*time.Second)
	defer shutdown()

	logger.Log(ctx).Info().Msg("shutting down server")
	return s.app.Shutdown()
}
