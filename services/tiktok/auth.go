package tiktok

import (
	"log"
	"net/http"
	"tiktok_tools/apperr"
	"tiktok_tools/client"
	"tiktok_tools/model"
	"tiktok_tools/model/tiktok"
	"tiktok_tools/services"

	"github.com/gin-gonic/gin"
)

// EmailSignup contains the user signup request
type TiktokSignup struct {
	OpenID          string `json:"open_id" binding:"required"`
	Email           string `json:"email" binding:"required,min=5,email"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

// NewUserService create a new user application service
func NewTiktokService(shopRepo *services.TiktokShopRepo, account *services.AccountRepo) *Service {
	return &Service{
		shopRepo:    shopRepo,
		accountRepo: account,
		tkCli: client.New(client.ClientConfig{
			AppKey:    "6b0h9912nnv8t",
			AppSecret: "f448132eb3c311e4fc65dc39bac08be332fb4736",
		}),
	}
}

// Service represents the user application service
type Service struct {
	shopRepo    *services.TiktokShopRepo
	accountRepo *services.AccountRepo
	tkCli       *client.TikTokClient
}

func (s *Service) AuthLink(c *gin.Context) {
	authCode := c.Query("code")
	log.Println(authCode)
	resp, err := s.tkCli.Cli.GetAccessToken(c.Request.Context(), authCode)
	if err != nil {
		log.Fatalf("get access token failed, err: %s", err.Error())
		return
	}

	shopItem, err := s.shopRepo.Create(&tiktok.TiktokShop{
		AccessToken:          resp.AccessToken,
		AccessTokenExpireIn:  resp.AccessTokenExpireIn,
		RefreshToken:         resp.RefreshToken,
		RefreshTokenExpireIn: resp.AccessTokenExpireIn,
		OpenID:               resp.OpenID,
		SellerName:           resp.SellerName,
	})
	if err != nil {
		log.Fatalf("create shop from tiktok failed: %s", err.Error())
		return
	}

	log.Printf("resp: %v", resp)
	c.JSON(200, gin.H{
		"message": shopItem,
	})
}

func (s *Service) BindUser(c *gin.Context) {
	var r TiktokSignup
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return
	}
	if r.Password != r.PasswordConfirm {
		err := apperr.New(http.StatusBadRequest, "passwords do not match")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	user, err := s.accountRepo.Create(&model.User{
		Email: r.Email,
	})
	if err != nil {
		log.Fatalf("create shop from tiktok failed: %s", err.Error())
		return
	}
	c.JSON(200, gin.H{
		"message": user,
	})
}

func (s *Service) WebHook(c *gin.Context) {
	authCode := c.Query("code")
	resp, err := s.tkCli.Cli.GetAccessToken(c.Request.Context(), authCode)
	if err != nil {
		log.Fatalf("get access token failed, err: %s", err.Error())
		return
	}
	log.Printf("resp: %v", resp)
}
