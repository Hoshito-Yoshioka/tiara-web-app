package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"tiara-web-app/backend/internal/infrastructure/db"
	handler "tiara-web-app/backend/internal/interface/handler" // handlerパッケージをエイリアス
	"tiara-web-app/backend/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database Connection
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		e.Logger.Fatal("DATABASE_URL environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		e.Logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Ping database to verify connection
	err = pool.Ping(ctx)
	if err != nil {
		e.Logger.Fatalf("Failed to ping database: %v\n", err)
	}
	fmt.Println("Successfully connected to the database!")

	// Dependencies
	queries := db.New(pool)
	shopRepo := db.NewShopRepository(queries)
	shopUsecase := usecase.NewShopUsecase(shopRepo)
	shopHandler := handler.NewShopHandler(shopUsecase)

	// Routes
	e.GET("/", healthCheck)
	e.GET("/shops", shopHandler.ListShops)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World! This is Tiara API.")
}
