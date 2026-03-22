package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"tiara-web-app/backend/internal/infrastructure/db"
	handler "tiara-web-app/backend/internal/interface/handler"
	authMiddleware "tiara-web-app/backend/internal/interface/middleware"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3001"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

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

	staffRepo := db.NewStaffRepository(queries, pool)
	staffUsecase := usecase.NewStaffUsecase(staffRepo)
	staffHandler := handler.NewStaffHandler(staffUsecase)

	adminRepo := db.NewAdminRepository(queries)
	authUsecase := usecase.NewAuthUsecase(adminRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	menuRepo := db.NewMenuRepository(queries)
	menuUsecase := usecase.NewMenuUsecase(menuRepo)
	menuHandler := handler.NewMenuHandler(menuUsecase)

	// Staff Portal Dependencies
	staffAccountRepo := db.NewStaffAccountRepository(queries)
	staffAuthUsecase := usecase.NewStaffAuthUsecase(staffAccountRepo)
	staffDraftRepo := db.NewStaffDraftRepository(queries, pool)
	staffPortalUsecase := usecase.NewStaffPortalUsecase(staffDraftRepo, staffRepo)
	staffPortalHandler := handler.NewStaffPortalHandler(staffAuthUsecase, staffPortalUsecase)
	adminReviewUsecase := usecase.NewAdminReviewUsecase(staffDraftRepo, staffRepo)
	adminReviewHandler := handler.NewAdminReviewHandler(adminReviewUsecase)
	adminAccountUsecase := usecase.NewAdminAccountUsecase(staffAccountRepo)
	adminAccountHandler := handler.NewAdminAccountHandler(adminAccountUsecase)

	// Static file serving for uploaded images
	e.Static("/uploads", "uploads")

	// Public Routes
	e.GET("/", healthCheck)
	e.GET("/shops", shopHandler.ListShops)
	e.GET("/shops/:id", shopHandler.GetShopByID)
	e.GET("/staffs", staffHandler.ListStaffs)
	e.GET("/staffs/:id", staffHandler.GetStaffWithSchedules)
	e.GET("/schedules", staffHandler.ListAllStaffsWithSchedules)
	e.GET("/menus", menuHandler.ListMenuCategoriesWithItems)

	// Auth Routes
	e.POST("/auth/login", authHandler.Login)

	// Staff Auth Routes (スタッフポータル認証)
	e.POST("/staff-auth/login", staffPortalHandler.Login)

	// Admin Routes (JWT Protected)
	admin := e.Group("/admin", authMiddleware.JWTAuth())
	admin.GET("/auth/verify", authHandler.Verify)
	admin.PUT("/shops/:id", shopHandler.UpdateShop)
	admin.POST("/staffs", staffHandler.CreateStaff)
	admin.PUT("/staffs/:id", staffHandler.UpdateStaff)
	admin.DELETE("/staffs/:id", staffHandler.DeleteStaff)
	admin.POST("/staffs/:id/images", staffHandler.UploadStaffImage)
	admin.DELETE("/staffs/:id/images/:imageId", staffHandler.DeleteStaffImage)
	admin.PUT("/staffs/:id/images/main", staffHandler.SetMainImage)
	admin.POST("/menu/categories", menuHandler.CreateMenuCategory)
	admin.PUT("/menu/categories/:id", menuHandler.UpdateMenuCategory)
	admin.DELETE("/menu/categories/:id", menuHandler.DeleteMenuCategory)
	admin.POST("/menu/items", menuHandler.CreateMenuItem)
	admin.PUT("/menu/items/:id", menuHandler.UpdateMenuItem)
	admin.DELETE("/menu/items/:id", menuHandler.DeleteMenuItem)
	// Admin Review Routes (下書きレビュー)
	admin.GET("/reviews/profiles", adminReviewHandler.ListPendingProfileDrafts)
	admin.PUT("/reviews/profiles/:id", adminReviewHandler.ReviewProfileDraft)
	admin.GET("/reviews/schedules", adminReviewHandler.ListPendingScheduleDrafts)
	admin.PUT("/reviews/schedules/:id", adminReviewHandler.ReviewScheduleDraft)
	// Admin Staff Account Management (スタッフアカウント管理)
	admin.GET("/staff-accounts", adminAccountHandler.ListStaffAccounts)
	admin.GET("/staff-accounts/staff/:staffId", adminAccountHandler.GetStaffAccountByStaffID)
	admin.POST("/staff-accounts", adminAccountHandler.CreateStaffAccount)
	admin.PUT("/staff-accounts/:id", adminAccountHandler.UpdateStaffAccount)
	admin.DELETE("/staff-accounts/:id", adminAccountHandler.DeleteStaffAccount)

	// Staff Portal Routes (スタッフ専用、Staff JWT Protected)
	portal := e.Group("/portal", authMiddleware.StaffJWTAuth())
	portal.GET("/auth/verify", staffPortalHandler.Verify)
	portal.GET("/profile", staffPortalHandler.GetMyProfileDraft)
	portal.PUT("/profile", staffPortalHandler.SaveMyProfileDraft)
	portal.POST("/profile/:id/submit", staffPortalHandler.SubmitMyProfileDraft)
	portal.GET("/schedule", staffPortalHandler.GetMyScheduleDraft)
	portal.PUT("/schedule", staffPortalHandler.SaveMyScheduleDraft)
	portal.POST("/schedule/:id/submit", staffPortalHandler.SubmitMyScheduleDraft)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World! This is Tiara API.")
}
