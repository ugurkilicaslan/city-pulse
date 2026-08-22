package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"

	"city-pulse/app/collections"
	"city-pulse/app/routes"
	"city-pulse/internal/config"
)

// Run — Uygulamanın ana başlangıç noktası
// 1. Config okur  2. MongoDB bağlar  3. Index oluşturur  4. HTTP server başlatır
func Run() {
	config.FlagConfig()
	config.ViperRead()

	printRunMode()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupGracefulShutdown(cancel)

	mongoClient, err := collections.Connect(ctx, config.Config.MongoURI)
	if err != nil {
		log.Fatal(color.RedString("MongoDB bağlantısı başarısız: "), err)
	}
	defer mongoClient.Disconnect(ctx)

	database := mongoClient.Database(config.Config.MongoDatabase)

	if err := collections.EnsureIndexes(ctx, database); err != nil {
		log.Fatal(color.RedString("Index oluşturma başarısız: "), err)
	}

	fmt.Println(color.GreenString("MongoDB bağlandı → ") + color.CyanString(config.Config.MongoDatabase))

	routes.InitializeRoutes(database, ctx)
}

func printRunMode() {
	color.White("....")
	if config.Config.IsLocal() {
		fmt.Println(color.CyanString("\t\t\tdevelopment mode"))
		return
	}
	fmt.Println(color.CyanString("\t\t\tprod mode"))
}

// setupGracefulShutdown — CTRL+C veya SIGTERM gelince temiz kapanma sağlar
func setupGracefulShutdown(cancel context.CancelFunc) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownCh
		fmt.Println()
		fmt.Println(color.YellowString("City Pulse durduruluyor..."))
		cancel()
	}()
}
