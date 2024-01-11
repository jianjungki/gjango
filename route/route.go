package route

import (
	"tiktok_tools/mail"
	mw "tiktok_tools/middleware"
	"tiktok_tools/mobile"
	"tiktok_tools/secret"
	"tiktok_tools/service"
	"tiktok_tools/services"
	"tiktok_tools/services/account"
	"tiktok_tools/services/apify"
	"tiktok_tools/services/auth"
	"tiktok_tools/services/tiktok"
	"tiktok_tools/services/user"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

// NewServices creates a new router services
func NewServices(DB *pg.DB, Log *zap.Logger, JWT *mw.JWT, Mail mail.Service, Mobile mobile.Service, R *gin.Engine) *Services {
	return &Services{DB, Log, JWT, Mail, Mobile, R}
}

// Services lets us bind specific services when setting up routes
type Services struct {
	DB     *pg.DB
	Log    *zap.Logger
	JWT    *mw.JWT
	Mail   mail.Service
	Mobile mobile.Service
	R      *gin.Engine
}

// SetupV1Routes instances various repos and services and sets up the routers
func (s *Services) SetupV1Routes() {
	// database logic
	tiktokShopRepo := services.NewTiktokShopRepo(s.DB, s.Log)
	userRepo := services.NewUserRepo(s.DB, s.Log)
	accountRepo := services.NewAccountRepo(s.DB, s.Log, secret.New())
	rbac := services.NewRBACService(userRepo)

	// service logic
	authService := auth.NewAuthService(userRepo, accountRepo, s.JWT, s.Mail, s.Mobile)
	accountService := account.NewAccountService(userRepo, accountRepo, rbac, secret.New())
	userService := user.NewUserService(userRepo, authService, rbac)

	tiktokService := tiktok.NewTiktokService(tiktokShopRepo, accountRepo)
	// no prefix, no jwt
	service.TiktokRouter(tiktokService, s.R)

	apifyService := apify.NewApifyService()
	// no prefix, no jwt
	service.ApifyRouter(apifyService, s.R)

	// no prefix, no jwt
	service.AuthRouter(authService, s.R)

	// prefixed with /v1 and protected by jwt
	v1Router := s.R.Group("/v1")
	v1Router.Use(s.JWT.MWFunc())
	service.AccountRouter(accountService, v1Router)
	service.UserRouter(userService, v1Router)
}
