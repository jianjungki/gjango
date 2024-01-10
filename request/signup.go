package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tiktok_tools/apperr"
)

// EmailSignup contains the user signup request
type EmailSignup struct {
	Email           string `json:"email" binding:"required,min=3,email"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

// AccountSignup validates user signup request
func AccountSignup(c *gin.Context) (*EmailSignup, error) {
	var r EmailSignup
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	if r.Password != r.PasswordConfirm {
		err := apperr.New(http.StatusBadRequest, "passwords do not match")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return nil, err
	}
	return &r, nil
}

// MobileSignup contains the user signup request with a mobile number
type MobileSignup struct {
	CountryCode string `json:"country_code" binding:"required,min=2"`
	Mobile      string `json:"mobile" binding:"required"`
}

// Mobile validates user signup request via mobile
func Mobile(c *gin.Context) (*MobileSignup, error) {
	var r MobileSignup
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return &r, nil
}

// MobileVerify contains the user's mobile verification country code, mobile number and verification code
type MobileVerify struct {
	CountryCode string `json:"country_code" binding:"required,min=2"`
	Mobile      string `json:"mobile" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Signup      bool   `json:"signup" binding:"required"`
}

// AccountVerifyMobile validates user mobile verification
func AccountVerifyMobile(c *gin.Context) (*MobileVerify, error) {
	var r MobileVerify
	if err := c.ShouldBindJSON(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
