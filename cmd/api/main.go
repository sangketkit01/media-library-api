package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sangketkit01/media-library-api/internal/config"
	"github.com/sangketkit01/media-library-api/internal/handlers"
	"github.com/sangketkit01/media-library-api/internal/middleware"
	"github.com/sangketkit01/media-library-api/internal/routes"
	"github.com/sangketkit01/media-library-api/internal/token"
)

const (
	webPort = "8099"
)

type App struct {
	handler    *handlers.Handler
	config     *config.Config
	tokenMaker token.Maker
	routes     *routes.Route
}

func init() {
	if err := godotenv.Load("../../.env.local"); err != nil {
		log.Println("failed to load .env.local. Try to load .env.production...")

		if err := godotenv.Load("../../.env.production"); err != nil {
			log.Println("failed to load .env.production")
		} else {
			log.Println("Success Loading .env.production")
		}
	} else {
		log.Println("Succeed Loading .env.local")
	}
}

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local"
	}

	config, err := config.NewConfig("../../", env)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("environment:", config.Environment)

	tokenMaker, err := token.NewPasetoMaker(config.Secretkey)
	if err != nil {
		log.Panic(err)
	}

	middleware := middleware.NewMiddleware(tokenMaker)

	handler, err := handlers.NewHandler(config, tokenMaker)
	if err != nil {
		log.Panic(err)
	}

	router := routes.NewRoute(middleware, handler)

	app := App{
		config:     config,
		tokenMaker: tokenMaker,
		handler:    handler,
		routes:     router,
	}

	uploadDir := "../../uploads"

	// Create shutdown-aware context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start background task with graceful shutdown
	go startBackgroundTask(ctx, uploadDir)

	// Start Fiber server
	go func() {
		if err := app.routes.Router.Listen(fmt.Sprintf(":%s", webPort)); err != nil {
			log.Printf("Fiber server error: %v\n", err)
		}
	}()

	log.Println("Server running...")

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown initiated...")

	// Shutdown Fiber server
	if err := app.routes.Router.Shutdown(); err != nil {
		log.Printf("Fiber shutdown error: %v\n", err)
	}

	// Close database pool
	app.handler.Pool.Close()

	log.Println("Graceful shutdown complete.")
}

func logTotalSizePerUser(uploadDir string) {
	users, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Println("Failed to read upload directory:", err)
		return
	}

	for _, userDir := range users {
		if !userDir.IsDir() {
			continue
		}

		userPath := filepath.Join(uploadDir, userDir.Name())

		var totalSize int64 = 0

		err := filepath.Walk(userPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				totalSize += info.Size()
			}
			return nil
		})

		if err != nil {
			log.Println("Error walking user folder", userDir.Name(), ":", err)
			continue
		}

		totalMB := float64(totalSize) / (1024 * 1024)
		log.Printf("User %s: Total size = %.2f MB\n", userDir.Name(), totalMB)
	}
}

func startBackgroundTask(ctx context.Context, uploadDir string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logTotalSizePerUser(uploadDir)

		case <-ctx.Done():
			log.Println("Background task stopping...")
			return
		}
	}
}
