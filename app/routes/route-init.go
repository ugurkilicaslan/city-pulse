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
	"city-pulse/app/middleware"
	"city-pulse/app/services"
	"city-pulse/internal/cache"
	"city-pulse/internal/config"
)

// InitializeRoutes — Gin router'ı kurar, tüm bağımlılıkları oluşturur, HTTP sunucuyu başlatır
func InitializeRoutes(db *mongo.Database, ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.RateLimit(100))
	r.Use(corsMiddleware(config.Config.IsLocal()))

	r.Static("/ui", "./ui")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/")
	})

	globalCache := cache.New()

	exchangeClient := clients.NewExchangeClient()
	newsClient := clients.NewNewsClient(config.Config.GNewsAPIKey)
	gameClient := clients.NewGameClient()
	nasaClient := clients.NewNASAClient(config.Config.NASAAPIKey)
	githubClient := clients.NewGitHubClient()
	weatherClient := clients.NewWeatherClient()
	cryptoClient := clients.NewCryptoClient()

	exchangeDAO := daos.NewExchangeDAO(db)
	newsDAO := daos.NewNewsDAO(db)
	nasaDAO := daos.NewNasaDAO(db)
	userDAO := daos.NewUserDAO(db)
	analyticsDAO := daos.NewAnalyticsDAO(db)
	prefDAO := daos.NewPreferenceDAO(db)
	bookmarkDAO := daos.NewBookmarkDAO(db)
	alertDAO := daos.NewAlertDAO(db)

	exchangeSvc := services.NewExchangeService(exchangeDAO, exchangeClient)
	newsSvc := services.NewNewsService(newsDAO, newsClient)
	gameSvc := services.NewGameService(gameClient)
	nasaSvc := services.NewNasaService(nasaClient, nasaDAO)
	githubSvc := services.NewGithubService(githubClient)
	citySvc := services.NewCityService(exchangeSvc, newsSvc, gameSvc, nasaSvc, githubSvc)
	userSvc := services.NewUserService(userDAO)
	weatherSvc := services.NewWeatherService(weatherClient, globalCache)
	cryptoSvc := services.NewCryptoService(cryptoClient, globalCache)
	analyticsSvc := services.NewAnalyticsService(analyticsDAO)
	prefSvc := services.NewPreferenceService(prefDAO)
	bookmarkSvc := services.NewBookmarkService(bookmarkDAO)
	alertSvc := services.NewAlertService(alertDAO)

	h := handlers.New(exchangeSvc, newsSvc, gameSvc, citySvc, userSvc, weatherSvc, cryptoSvc, analyticsSvc, prefSvc, bookmarkSvc, alertSvc, ctx)

	api := r.Group(config.Config.APIBasePath())
	api.GET("/health", h.Healthz)

	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)

	protected := api.Group("/")
	protected.Use(auth.JWTMiddleware())
	{
		protected.GET("/exchange/latest", h.ExchangeLatest)
		protected.GET("/exchange/convert", h.ExchangeConvert)

		protected.GET("/news/headlines", h.NewsHeadlines)
		protected.GET("/news/search", h.NewsSearch)

		protected.GET("/games/deals", h.GameDeals)
		protected.GET("/games/search", h.GameSearch)

		protected.GET("/city/snapshot", h.CitySnapshot)

		protected.GET("/weather/current", h.WeatherCurrent)
		protected.GET("/crypto/prices", h.CryptoPrices)

		protected.GET("/me", h.MeProfile)
		protected.GET("/me/preferences", h.MePreferences)
		protected.PUT("/me/preferences", h.MeUpdatePreferences)
		protected.DELETE("/me/preferences", h.MeResetPreferences)

		protected.GET("/me/bookmarks", h.BookmarkList)
		protected.POST("/me/bookmarks", h.BookmarkCreate)
		protected.DELETE("/me/bookmarks/:id", h.BookmarkDelete)
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
