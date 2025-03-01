package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tiktok_tools/apperr"
	"tiktok_tools/model"
)

// RegisterAdmin contains admin registration request
type RegisterAdmin struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Username        string `json:"username" binding:"required,min=3,alphanum"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	Email           string `json:"email" binding:"required,email"`

	CompanyID  int `json:"company_id" binding:"required"`
	LocationID int `json:"location_id" binding:"required"`
	RoleID     int `json:"role_id" binding:"required"`
}

// AccountCreate validates account creation request
func AccountCreate(c *gin.Context) (*RegisterAdmin, error) {
	var r RegisterAdmin
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	if r.Password != r.PasswordConfirm {
		err := apperr.New(http.StatusBadRequest, "passwords do not match")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return nil, err
	}
	if r.RoleID < int(model.SuperAdminRole) || r.RoleID > int(model.UserRole) {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, apperr.BadRequest
	}
	return &r, nil
}

// Password contains password change request
type Password struct {
	ID                 int    `json:"-"`
	OldPassword        string `json:"old_password" binding:"required,min=8"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" binding:"required"`
}

// PasswordChange validates password change request
func PasswordChange(c *gin.Context) (*Password, error) {
	var p Password
	id, err := ID(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	if p.NewPassword != p.NewPasswordConfirm {
		err := apperr.New(http.StatusBadRequest, "passwords do not match")
		apperr.Response(c, err)
		return nil, err
	}
	p.ID = id
	return &p, nil
}
