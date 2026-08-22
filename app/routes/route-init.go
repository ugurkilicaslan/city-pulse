package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"city-pulse/app/auth"
	"city-pulse/app/clients"
	"city-pulse/app/daos"
	"city-pulse/app/handlers"
	"city-pulse/app/services"
	"city-pulse/internal/config"
)

// InitializeRoutes — Gin router'ı kurar, tüm bağımlılıkları oluşturur, HTTP sunucuyu başlatır
func InitializeRoutes(db *mongo.Database, ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	if config.DevMode {
		r.Use(gin.Logger())
	}
	r.Use(corsMiddleware(config.Config.IsLocal()))

	r.Static("/ui", "./ui")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/")
	})

	exchangeClient := clients.NewExchangeClient()
	newsClient := clients.NewNewsClient(config.Config.GNewsAPIKey)
	gameClient := clients.NewGameClient()
	nasaClient := clients.NewNASAClient(config.Config.NASAAPIKey)
	githubClient := clients.NewGitHubClient()

	exchangeDAO := daos.NewExchangeDAO(db)
	newsDAO := daos.NewNewsDAO(db)
	nasaDAO := daos.NewNasaDAO(db)

	exchangeSvc := services.NewExchangeService(exchangeDAO, exchangeClient)
	newsSvc := services.NewNewsService(newsDAO, newsClient)
	gameSvc := services.NewGameService(gameClient)
	nasaSvc := services.NewNasaService(nasaClient, nasaDAO)
	githubSvc := services.NewGithubService(githubClient)
	citySvc := services.NewCityService(exchangeSvc, newsSvc, gameSvc, nasaSvc, githubSvc)

	h := handlers.New(exchangeSvc, newsSvc, gameSvc, citySvc, ctx)

	api := r.Group(config.Config.APIBasePath())
	api.GET("/health", h.Healthz)

	authorized := api.Group("/")
	authorized.Use(auth.APIKeyMiddleware())
	{
		authorized.GET("/exchange/latest", h.ExchangeLatest)
		authorized.GET("/exchange/convert", h.ExchangeConvert)

		authorized.GET("/news/headlines", h.NewsHeadlines)
		authorized.GET("/news/search", h.NewsSearch)

		authorized.GET("/games/deals", h.GameDeals)
		authorized.GET("/games/search", h.GameSearch)

		authorized.GET("/city/snapshot", h.CitySnapshot)
	}

	fmt.Println()
	fmt.Println()
	color.White("...")
	fmt.Println(
		color.GreenString("√"),
		color.YellowString("WiseUP City Pulse Service Started"),
		color.GreenString("√"),
	)
	fmt.Println(color.CyanString("port:%d  base:%s", config.Config.ServerPort, config.Config.APIBasePath()))

	addr := fmt.Sprintf(":%d", config.Config.ServerPort)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	serverErrors := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		log.Fatal(color.RedString("HTTP server başlatılamadı: "), err)
	case <-ctx.Done():
		fmt.Println(color.YellowString("HTTP server durduruluyor..."))

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf(color.RedString("HTTP shutdown hatası: %v"), err)
		}

		fmt.Println(color.GreenString("City Pulse durduruldu."))
	}
}
