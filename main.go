package main

import (
	"fmt"
	"log"

	"github.com/unify-z/go-surl/internal/config"
	"github.com/unify-z/go-surl/internal/helpers"
	"github.com/unify-z/go-surl/internal/logger"
	"github.com/unify-z/go-surl/internal/repository/surl"
	"github.com/unify-z/go-surl/internal/repository/user"
	"github.com/unify-z/go-surl/internal/routes"
	ShortURLService "github.com/unify-z/go-surl/internal/service/surl"
	UserService "github.com/unify-z/go-surl/internal/service/user"
	verifyCodeService "github.com/unify-z/go-surl/internal/service/verify_code"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfgMgr, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg := cfgMgr.GetAppConfig()
	_ = logger.InitLogger(logger.GetLogLevelInt(cfg.Log.Level), cfg.Log.FilePath)
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&user.UserModel{}, &surl.ShortURL{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	jwtHelper := helpers.NewJWTHelper(cfg.Jwt.SecretKey, cfg.Jwt.TokenDuration)
	smtpHelper := helpers.NewSMTPHelper(
		cfg.SMTP.Server,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		false,
	)
	surlRepo := surl.NewShortURLRepo(db)
	userRepo := user.NewUserRepo(db)
	surlSvc := ShortURLService.NewShortURLManager(surlRepo)
	verifyCodeSvc := verifyCodeService.NewVerifyCodeManager(*smtpHelper)
	userSvc := UserService.NewUserManager(userRepo, jwtHelper, verifyCodeSvc)
	router := routes.SetupRouter(surlSvc, userSvc, verifyCodeSvc)
	addr := fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
