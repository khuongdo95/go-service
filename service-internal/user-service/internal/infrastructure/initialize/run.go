package initialize

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/khuongdo95/go-pkg/appLogger"
	"github.com/khuongdo95/go-pkg/common"
	"github.com/khuongdo95/go-pkg/common/response"
	"github.com/khuongdo95/go-pkg/database/mysql"
	"github.com/khuongdo95/go-pkg/extractor"
	"github.com/khuongdo95/go-service/internal/adapter/controllers"
	"github.com/khuongdo95/go-service/internal/infrastructure/core"
	"github.com/khuongdo95/go-service/internal/infrastructure/global"
	"github.com/khuongdo95/go-service/internal/infrastructure/routers"
	"github.com/khuongdo95/go-service/internal/usecase/services"
	"github.com/khuongdo95/go-service/internal/usecase/services/cache"
	"github.com/khuongdo95/go-service/internal/usecase/services/signin"
	"github.com/khuongdo95/go-service/internal/usecase/services/signup"
	tokenPkg "github.com/khuongdo95/go-service/internal/usecase/services/token"

	redisClient "github.com/khuongdo95/go-pkg/caching/redis"
	"go.uber.org/zap"
)

func must[T any](value T, err *response.AppError) T {
	if err != nil {
		panic(err)
	}
	return value
}

func Run() {
	global.AppConfig = must(common.LoadConfig[global.Config](""))
	global.Log = must(appLogger.NewLogger(global.AppConfig.Logger, global.AppConfig.Server))
	global.SQLDB = must(mysql.NewConnection(global.AppConfig.SQL))
	global.Ent = must(core.NewEntClient(global.SQLDB, global.AppConfig.SQL))
	global.Cache = must(redisClient.NewRedisClient(global.AppConfig.Cache))
	global.App = &global.HttpServer{
		App:  must(core.NewHttpServer(global.AppConfig.AllowOrigins, global.Log)),
		Port: global.AppConfig.Server.Port,
		Name: global.AppConfig.Server.ServiceName,
	}

	global.Ent = must(core.NewEntClient(global.SQLDB, global.AppConfig.SQL))
	global.Ext = extractor.New()

	cache := cache.New(global.Ent.Client(), global.Cache)

	token := must(tokenPkg.New(global.AppConfig.AccessTokenSigning, global.AppConfig.IDTokenSigning))
	userSer := services.NewUserService(
		global.Ent.Client(),
		signup.New(global.Ent.Client(), token, cache),
		signin.New(global.Ent.Client(), token, cache),
	)
	ipSer := services.NewIpAccess(global.Ent.Client())

	IpController := controllers.NewIpAccessController(ipSer)
	UserController := controllers.NewUserController(userSer)
	routers.UserRouters(global.App, UserController)
	routers.IpAccessRouters(global.App, IpController)

	Start(global.App)
}

func Start(s *global.HttpServer) {
	defer func() {
		if err := global.Log.Logger().Sync(); err != nil {
			global.Log.Info(fmt.Sprintf("failed to sync logger: %v", err.Error()))
		}
	}()

	// Graceful shutdown setup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%d", s.Port)
		if err := s.Listen(addr); err != nil {
			serverErr <- err
		}
	}()
	// Wait for interrupt signal or server error
	select {
	case err := <-serverErr:
		global.Log.Info(fmt.Sprintf("server error: %v", err.Error()))
		return
	case sig := <-sigChan:
		global.Log.Info("received signal", zap.String("signal", sig.String()))
	}
	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Shutdown server

	if err := s.ShutdownWithContext(ctx); err != nil {
		global.Log.Info(fmt.Sprintf("graceful shutdown error: %v", err.Error()))
		return
	}

	global.Log.Info("graceful shutdown completed", nil)
}
