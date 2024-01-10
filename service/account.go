package service

import (
	"net/http"

	"tiktok_tools/apperr"
	"tiktok_tools/model"
	"tiktok_tools/request"
	"tiktok_tools/services/account"

	"github.com/gin-gonic/gin"
)

// AccountService represents the account http service
type AccountService struct {
	svc *account.Service
}

// AccountRouter sets up all the controller functions to our router
func AccountRouter(svc *account.Service, r *gin.RouterGroup) {
	a := AccountService{
		svc: svc,
	}
	ar := r.Group("/users")
	ar.POST("", a.create)
	ar.PATCH("/:id/password", a.changePassword)
}

func (a *AccountService) create(c *gin.Context) {
	r, err := request.AccountCreate(c)
	if err != nil {
		return
	}
	user := &model.User{
		Username:   r.Username,
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		CompanyID:  r.CompanyID,
		LocationID: r.LocationID,
		RoleID:     r.RoleID,
	}
	if err := a.svc.Create(c, user); err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (a *AccountService) changePassword(c *gin.Context) {
	p, err := request.PasswordChange(c)
	if err != nil {
		return
	}
	if err := a.svc.ChangePassword(c, p.OldPassword, p.NewPassword, p.ID); err != nil {
		apperr.Response(c, err)
		return
	}
	c.Status(http.StatusOK)
}
