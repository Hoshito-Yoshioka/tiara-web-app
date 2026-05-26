package main

import (
	"net/http"

	handler "tiara-web-app/backend/internal/interface/handler"
	authMiddleware "tiara-web-app/backend/internal/interface/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// handlers はDIで構築されたハンドラー群を束ねる構造体。
// main から routes へ渡す唯一の依存。
type handlers struct {
	shop         *handler.ShopHandler
	staff        *handler.StaffHandler
	auth         *handler.AuthHandler
	menu         *handler.MenuHandler
	staffPortal  *handler.StaffPortalHandler
	adminReview  *handler.AdminReviewHandler
	adminAccount *handler.AdminAccountHandler
}

// registerRoutes は全ルートを Echo インスタンスに登録する。
func registerRoutes(e *echo.Echo, h *handlers, jwtSecret string) {
	// ヘルスチェック（バージョニング対象外）
	e.GET("/", healthCheck)

	// API v1 グループ
	v1 := e.Group("/api/v1")
	registerPublicRoutes(v1, h)
	registerAuthRoutes(v1, h)
	registerAdminRoutes(v1, h, jwtSecret)
	registerPortalRoutes(v1, h, jwtSecret)
}

// --- Public ---

func registerPublicRoutes(g *echo.Group, h *handlers) {
	g.GET("/shops", h.shop.ListShops)
	g.GET("/shops/:id", h.shop.GetShopByID)
	g.GET("/staffs", h.staff.ListStaffs)
	g.GET("/staffs/:id", h.staff.GetStaffWithSchedules)
	g.GET("/schedules", h.staff.ListAllStaffsWithSchedules)
	g.GET("/menus", h.menu.ListMenuCategoriesWithItems)
}

// --- Auth ---

// loginRateLimiter はログインエンドポイント用のレート制限ミドルウェア。
// 1秒あたり5リクエスト、バースト10まで許容。
var loginRateLimiter = middleware.RateLimiter(
	middleware.NewRateLimiterMemoryStore(rate.Limit(5)),
)

func registerAuthRoutes(g *echo.Group, h *handlers) {
	g.POST("/auth/login", h.auth.Login, loginRateLimiter)
	g.POST("/staff-auth/login", h.staffPortal.Login, loginRateLimiter)
}

// --- Admin (JWT Protected) ---

func registerAdminRoutes(g *echo.Group, h *handlers, jwtSecret string) {
	admin := g.Group("/admin", authMiddleware.JWTAuth(jwtSecret))

	// Auth verify
	admin.GET("/auth/verify", h.auth.Verify)

	// Shops
	admin.PUT("/shops/:id", h.shop.UpdateShop)

	// Staffs
	admin.POST("/staffs", h.staff.CreateStaff)
	admin.PUT("/staffs/:id", h.staff.UpdateStaff)
	admin.DELETE("/staffs/:id", h.staff.DeleteStaff)
	admin.POST("/staffs/:id/images", h.staff.UploadStaffImage)
	admin.DELETE("/staffs/:id/images/:imageId", h.staff.DeleteStaffImage)
	admin.PUT("/staffs/:id/images/main", h.staff.SetMainImage)
	admin.PUT("/staffs/:id/images/:imageId/crop", h.staff.UpdateImageCropPosition)

	// Menu
	admin.POST("/menu/categories", h.menu.CreateMenuCategory)
	admin.PUT("/menu/categories/:id", h.menu.UpdateMenuCategory)
	admin.DELETE("/menu/categories/:id", h.menu.DeleteMenuCategory)
	admin.POST("/menu/items", h.menu.CreateMenuItem)
	admin.PUT("/menu/items/:id", h.menu.UpdateMenuItem)
	admin.DELETE("/menu/items/:id", h.menu.DeleteMenuItem)

	// Draft reviews
	admin.GET("/reviews/profiles", h.adminReview.ListPendingProfileDrafts)
	admin.GET("/reviews/profiles/:id", h.adminReview.GetProfileDraft)
	admin.PUT("/reviews/profiles/:id", h.adminReview.ReviewProfileDraft)
	admin.PUT("/reviews/profiles/:id/content", h.adminReview.UpdateProfileDraftContent)
	admin.GET("/reviews/schedules", h.adminReview.ListPendingScheduleDrafts)
	admin.GET("/reviews/schedules/approved", h.adminReview.ListApprovedScheduleDrafts)
	admin.GET("/reviews/schedules/:id", h.adminReview.GetScheduleDraft)
	admin.PUT("/reviews/schedules/:id", h.adminReview.ReviewScheduleDraft)
	admin.PUT("/reviews/schedules/:id/content", h.adminReview.UpdateScheduleDraftContent)
	admin.POST("/reviews/schedules/:id/publish", h.adminReview.PublishScheduleDraft)

	// Staff accounts
	admin.GET("/staff-accounts", h.adminAccount.ListStaffAccounts)
	admin.GET("/staff-accounts/staff/:staffId", h.adminAccount.GetStaffAccountByStaffID)
	admin.POST("/staff-accounts", h.adminAccount.CreateStaffAccount)
	admin.PUT("/staff-accounts/:id", h.adminAccount.UpdateStaffAccount)
	admin.DELETE("/staff-accounts/:id", h.adminAccount.DeleteStaffAccount)
}

// --- Staff Portal (Staff JWT Protected) ---

func registerPortalRoutes(g *echo.Group, h *handlers, jwtSecret string) {
	portal := g.Group("/portal", authMiddleware.StaffJWTAuth(jwtSecret))

	portal.GET("/auth/verify", h.staffPortal.Verify)
	portal.GET("/profile", h.staffPortal.GetMyProfileDraft)
	portal.PUT("/profile", h.staffPortal.SaveMyProfileDraft)
	portal.POST("/profile/:id/submit", h.staffPortal.SubmitMyProfileDraft)
	portal.GET("/schedule", h.staffPortal.GetMyScheduleDraft)
	portal.PUT("/schedule", h.staffPortal.SaveMyScheduleDraft)
	portal.POST("/schedule/:id/submit", h.staffPortal.SubmitMyScheduleDraft)

	// Image management (自分の画像のみ操作可能)
	portal.GET("/images", h.staffPortal.ListMyImages)
	portal.POST("/images", h.staffPortal.UploadMyImage)
	portal.DELETE("/images/:imageId", h.staffPortal.DeleteMyImage)
	portal.PUT("/images/:imageId/crop", h.staffPortal.UpdateMyImageCropPosition)
	portal.PUT("/images/main", h.staffPortal.SetMyMainImage)
}

// healthCheck はヘルスチェック用のハンドラー。
func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World! This is Tiara API.")
}
