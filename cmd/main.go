package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Asliddin3/energy-maximum/config"
	"github.com/Asliddin3/energy-maximum/controller"
	"github.com/Asliddin3/energy-maximum/docs"
	"github.com/Asliddin3/energy-maximum/migrate"
	"github.com/Asliddin3/energy-maximum/pkg/hash"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/Asliddin3/energy-maximum/pkg/middleware"
	postgresdb "github.com/Asliddin3/energy-maximum/pkg/postgres"
	"github.com/allegro/bigcache/v3"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/mvrilo/go-redoc"
	ginredoc "github.com/mvrilo/go-redoc/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go.uber.org/automaxprocs"
)

// @title           Golang CRM Swagger Documentation
// @version         1.0
// @description     This is a sample server CRM server.

// @securityDefinitions.apikey  BearerAuth
// @host localhost:8000
// @in header
// @name Authorization
// @Description									AUTH.
func main() {
	cfg := config.Load()
	// log := logger.New(cfg.LogLevel, "crm-go-service")
	log := logger.NewLogger()
	doc := redoc.Redoc{
		Title:       "Example API",
		Description: "Example API Description",
		SpecFile:    "./docs/swagger.json", // "./openapi.yaml"
		SpecPath:    "/docs/swagger.json",  // "
		DocsPath:    "/redoc",
	}

	server := gin.Default()
	server.Use(
		ginredoc.New(doc),
		gin.Recovery(),
		gin.Logger(),
		middleware.New(middleware.GinCorsMiddleware()),
		gzip.Gzip(gzip.DefaultCompression),
	)

	db, err := postgresdb.NewClient(cfg)
	if err != nil {
		log.Error("postgresdb connection error", logger.Error(err))
		return
	}
	err = migrate.Migrate(db)
	if err != nil {
		log.Error("failed to migrate db", logger.Error(err))
		return
	}
	ctx := context.Background()
	cache, err := bigcache.New(
		ctx,
		bigcache.Config{
			Shards:             1024,
			LifeWindow:         10,
			MaxEntriesInWindow: 0,
			MaxEntrySize:       60,
			Verbose:            false,
		},
	)
	if err != nil {
		log.Error("failed to initialize bigcashe", logger.Error(err))
		return
	}
	hash := hash.NewHasher()
	handler := controller.NewHandler(db, log, cache, cfg, hash)

	handler.Init(server)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	err = server.Run(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Error("router running error", logger.Error(err))
		return
	}
	// err = migration.A
}
