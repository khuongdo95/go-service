package global

import (
	"github.com/gofiber/fiber/v3"
	"github.com/khuongdo95/go-pkg/appLogger"
	redisClient "github.com/khuongdo95/go-pkg/caching/redis"
	"github.com/khuongdo95/go-pkg/database/mysql"
	"github.com/khuongdo95/go-pkg/extractor"
	"github.com/khuongdo95/go-pkg/settings"
	"github.com/khuongdo95/go-service/internal/infrastructure/core"
)

type (
	HttpServer struct {
		*fiber.App
		Port int
		Name string
	}

	Config struct {
		AllowOrigins       string                   `mapstructure:"ALLOW_ORIGINS"`
		Server             *settings.ServerSetting  `mapstructure:"SERVER"`
		Logger             *appLogger.LoggerConfig  `mapstructure:"LOGGER"`
		Cache              *redisClient.CacheConfig `mapstructure:"CACHE"`
		SQL                *mysql.SQLConfig         `mapstructure:"SQL"`
		AccessTokenSigning *JwtSigning              `mapstructure:"ACCESS_TOKEN_SIGNING"`
		IDTokenSigning     *JwtSigning              `mapstructure:"ID_TOKEN_SIGNING"`
	}
	JwtSigning struct {
		PrivateKey  string `mapstructure:"PRIVATE_KEY"`
		Issuer      string `mapstructure:"ISSUER"`
		ExpiresTime int    `mapstructure:"EXPIRES_TIME"`
	}
)

var (
	Log                *appLogger.Logger
	AppConfig          *Config
	App                *HttpServer
	SQLDB              *mysql.Connection
	Ent                *core.EntClient
	Cache              *redisClient.CacheClient
	AccessTokenSigning *JwtSigning
	IDTokenSigning     *JwtSigning
	Ext                extractor.Extractor
)
