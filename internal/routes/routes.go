package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/handlers"
	"github.com/unify-z/go-surl/internal/handlers/admin"
	"github.com/unify-z/go-surl/internal/middlewares"
	ShortURLService "github.com/unify-z/go-surl/internal/service/surl"
	UserService "github.com/unify-z/go-surl/internal/service/user"
	verifyCodeService "github.com/unify-z/go-surl/internal/service/verify_code"
)

func SetupRouter(
	surlSvc *ShortURLService.ShortURLManager,
	userSvc *UserService.UserManager,
	verifyCodeSvc *verifyCodeService.VerifyCodeManager,
) *gin.Engine {
	router := gin.Default()
	surlHandler := handlers.NewSURLHandler(surlSvc)
	userHandler := handlers.NewUserHandler(userSvc, verifyCodeSvc)
	configHandler := handlers.NewConfigHandler()
	public := router.Group("/")
	public.GET("/s/:short_code", surlHandler.RedirectToOriginalURL)
	public.GET("/api/config", configHandler.GetSiteConfig)
	userApi := router.Group("/api/user")
	userApi.POST("/login", userHandler.Login)
	userApi.POST("/register", userHandler.Register)
	userApi.POST("/create_email_code", userHandler.SendEmailVerifyCode)
	userApi.GET("/info", middlewares.AuthMiddleware(userSvc), userHandler.GetUserInfo)
	surlApi := router.Group("/api/surl")
	surlApi.Use(middlewares.AuthMiddleware(userSvc))
	surlApi.POST("/create", surlHandler.CreateShortURL)
	surlApi.GET("/list", surlHandler.ListUserShortURLs)
	surlApi.POST("/delete", surlHandler.DeleteShortURL)
	surlApi.POST("/update", surlHandler.EditShortURL)
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(middlewares.AuthMiddleware(userSvc))
	adminGroup.Use(middlewares.RequireAdmin())
	adminSURLHandler := admin.NewAdminSURLHandler(surlSvc)
	adminUserHandler := admin.NewAdminUserHandler(userSvc)
	adminGroup.POST("/surl/delete", adminSURLHandler.DeleteSURL)
	adminGroup.GET("/surl/list", adminSURLHandler.ListShortURL)
	adminGroup.POST("/surl/update", adminSURLHandler.EditShortURL)
	adminGroup.POST("/user/delete", adminUserHandler.DeleteUser)
	adminGroup.GET("/user/list", adminUserHandler.GetUser)
	adminGroup.POST("/user/ban", adminUserHandler.BanUser)
	adminGroup.POST("/user/unban", adminUserHandler.UnbanUser)
	return router
}
