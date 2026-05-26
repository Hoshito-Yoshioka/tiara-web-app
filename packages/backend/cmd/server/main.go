package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"tiara-web-app/backend/internal/config"
	"tiara-web-app/backend/internal/infrastructure/db"
	handler "tiara-web-app/backend/internal/interface/handler"
	"tiara-web-app/backend/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// --- Config ---
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// --- Database ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}
	fmt.Println("Successfully connected to the database!")

	// --- DI ---
	h := buildHandlers(pool, cfg)

	// --- Echo ---
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		ContentSecurityPolicy: "default-src 'self'",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Static file serving
	e.Static("/uploads", cfg.UploadDir)

	// Route registration
	registerRoutes(e, h, cfg.JWTSecret)

	// --- Graceful Shutdown ---
	// シグナル(Ctrl-C / SIGTERM)を受けたら安全に停止する。
	srvCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := e.Start(":" + cfg.Port); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	<-srvCtx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("Server stopped gracefully")
}

// buildHandlers は全ハンドラーのDI構築を行う。
func buildHandlers(pool *pgxpool.Pool, cfg *config.Config) *handlers {
	queries := db.New(pool)

	shopRepo := db.NewShopRepository(queries)
	shopUsecase := usecase.NewShopUsecase(shopRepo)
	shopHandler := handler.NewShopHandler(shopUsecase)

	staffRepo := db.NewStaffRepository(queries, pool)
	staffUsecase := usecase.NewStaffUsecase(staffRepo)
	staffHandler := handler.NewStaffHandler(staffUsecase)

	adminRepo := db.NewAdminRepository(queries)
	authUsecase := usecase.NewAuthUsecase(adminRepo)
	authHandler := handler.NewAuthHandler(authUsecase, cfg.JWTSecret, cfg.JWTExpiryHours)

	menuRepo := db.NewMenuRepository(queries, pool)
	menuUsecase := usecase.NewMenuUsecase(menuRepo)
	menuHandler := handler.NewMenuHandler(menuUsecase)

	staffAccountRepo := db.NewStaffAccountRepository(queries)
	staffAuthUsecase := usecase.NewStaffAuthUsecase(staffAccountRepo)
	staffDraftRepo := db.NewStaffDraftRepository(queries, pool)
	staffPortalUsecase := usecase.NewStaffPortalUsecase(staffDraftRepo, staffRepo)
	staffPortalHandler := handler.NewStaffPortalHandler(staffAuthUsecase, staffPortalUsecase, staffUsecase,
		handler.WithJWTSecret(cfg.JWTSecret),
		handler.WithJWTExpiryHours(cfg.JWTExpiryHours),
		handler.WithUploadDir(cfg.UploadDir),
	)
	adminReviewUsecase := usecase.NewAdminReviewUsecase(staffDraftRepo, staffRepo)
	adminReviewHandler := handler.NewAdminReviewHandler(adminReviewUsecase)
	adminAccountUsecase := usecase.NewAdminAccountUsecase(staffAccountRepo)
	adminAccountHandler := handler.NewAdminAccountHandler(adminAccountUsecase)

	return &handlers{
		shop:         shopHandler,
		staff:        staffHandler,
		auth:         authHandler,
		menu:         menuHandler,
		staffPortal:  staffPortalHandler,
		adminReview:  adminReviewHandler,
		adminAccount: adminAccountHandler,
	}
}
