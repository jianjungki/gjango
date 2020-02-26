package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/mock"
	"github.com/calvinchengx/gin-go-pg/mock/mockdb"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository/auth"
	"github.com/calvinchengx/gin-go-pg/service"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		wantResp    *model.AuthToken
		userRepo    *mockdb.User
		accountRepo *mockdb.Account
		jwt         *mock.JWT
		m           *mock.Mail
		mobile      *mock.Mobile
	}{
		{
			name:       "Invalid request",
			req:        `{"username":"juzernejm"}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Fail on FindByUsername",
			req:        `{"username":"juzernejm","password":"hunter123"}`,
			wantStatus: http.StatusInternalServerError,
			userRepo: &mockdb.User{
				FindByUsernameFn: func(string) (*model.User, error) {
					return nil, apperr.DB
				},
			},
		},
		{
			name:       "Success",
			req:        `{"username":"juzernejm","password":"hunter123"}`,
			wantStatus: http.StatusOK,
			userRepo: &mockdb.User{
				FindByUsernameFn: func(string) (*model.User, error) {
					return &model.User{
						Password: auth.HashPassword("hunter123"),
						Active:   true,
					}, nil
				},
				UpdateLoginFn: func(*model.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(*model.User) (string, string, error) {
					return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
				},
			},
			wantResp: &model.AuthToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339)},
		},
	}
	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			authService := auth.NewAuthService(tt.userRepo, tt.accountRepo, tt.jwt, tt.m, tt.mobile)
			service.AuthRouter(authService, r)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/login"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(model.AuthToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				tt.wantResp.RefreshToken = response.RefreshToken
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestRefresh(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		wantResp    *model.RefreshToken
		userRepo    *mockdb.User
		accountRepo *mockdb.Account
		jwt         *mock.JWT
		m           *mock.Mail
		mobile      *mock.Mobile
	}{
		{
			name:       "Fail on FindByToken",
			req:        "refreshtoken",
			wantStatus: http.StatusInternalServerError,
			userRepo: &mockdb.User{
				FindByTokenFn: func(string) (*model.User, error) {
					return nil, apperr.DB
				},
			},
		},
		{
			name:       "Success",
			req:        "refreshtoken",
			wantStatus: http.StatusOK,
			userRepo: &mockdb.User{
				FindByTokenFn: func(string) (*model.User, error) {
					return &model.User{
						Username: "johndoe",
						Active:   true,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(*model.User) (string, string, error) {
					return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
				},
			},
			wantResp: &model.RefreshToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339)},
		},
	}
	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			authService := auth.NewAuthService(tt.userRepo, tt.accountRepo, tt.jwt, tt.m, tt.mobile)
			service.AuthRouter(authService, r)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/refresh/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(model.RefreshToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestSignup(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		userRepo    *mockdb.User
		accountRepo *mockdb.Account
		jwt         *mock.JWT
		m           *mock.Mail
		mobile      *mock.Mobile
	}{
		{
			name:       "Success",
			req:        `{"email":"juzernejm","password":"hunter123","password_confirm":"hunter123"}`,
			wantStatus: http.StatusCreated,
			userRepo: &mockdb.User{ // no such user, so create
				FindByEmailFn: func(string) (*model.User, error) {
					return nil, apperr.DB
				},
			},
			accountRepo: &mockdb.Account{
				CreateAndVerifyFn: func(*model.User) (*model.Verification, error) {
					return &model.Verification{
						Token:  "some-random-token-for-verification",
						UserID: 1,
					}, nil
				},
			},
			m: &mock.Mail{
				SendVerificationEmailFn: func(string, *model.Verification) error {
					return nil
				},
			},
		},
		{
			name:       "Failure because no password",
			req:        `{"email":"calvin","password":"","password_confirm":""}`,
			wantStatus: http.StatusInternalServerError,
			userRepo: &mockdb.User{ // no such user, so create
				FindByUsernameFn: func(string) (*model.User, error) {
					return nil, apperr.DB
				},
			},
			accountRepo: &mockdb.Account{
				CreateAndVerifyFn: func(*model.User) (*model.Verification, error) {
					return &model.Verification{
						Token:  "some-random-token-for-verification",
						UserID: 1,
					}, nil
				},
			},
			m: &mock.Mail{
				SendVerificationEmailFn: func(string, *model.Verification) error {
					return nil
				},
			},
		},
		{
			name:       "Failure because user already exists",
			req:        `{"email":"calvin","password":"whatever123","password_confirm":"whatever123"}`,
			wantStatus: http.StatusConflict,
			userRepo: &mockdb.User{ // user already exists
				FindByEmailFn: func(string) (*model.User, error) {
					return &model.User{
						Username: "calvin",
						Active:   true,
					}, nil
				},
			},
			accountRepo: &mockdb.Account{
				CreateAndVerifyFn: func(*model.User) (*model.Verification, error) {
					return &model.Verification{
						Token:  "some-random-token-for-verification",
						UserID: 1,
					}, nil
				},
			},
			m: &mock.Mail{
				SendVerificationEmailFn: func(string, *model.Verification) error {
					return nil
				},
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			authService := auth.NewAuthService(tt.userRepo, tt.accountRepo, tt.jwt, tt.m, tt.mobile)
			service.AuthRouter(authService, r)
			ts := httptest.NewServer(r)
			defer ts.Close()
			// signup
			path := ts.URL + "/signup"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestVerification(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		userRepo    *mockdb.User
		accountRepo *mockdb.Account
		jwt         *mock.JWT
		m           *mock.Mail
		mobile      *mock.Mobile
	}{
		{
			name:       "Success",
			req:        "some-random-verification-token",
			wantStatus: http.StatusOK,
			accountRepo: &mockdb.Account{
				FindVerificationTokenFn: func(context.Context, string) (*model.Verification, error) {
					return &model.Verification{
						Token:  "some-random-token-for-verification",
						UserID: 1,
					}, nil
				},
				DeleteVerificationTokenFn: func(context.Context, *model.Verification) error {
					return nil
				},
			},
		},
		{
			name:       "Failed",
			req:        "some-random-verification-token",
			wantStatus: http.StatusNotFound,
			accountRepo: &mockdb.Account{
				FindVerificationTokenFn: func(context.Context, string) (*model.Verification, error) {
					return nil, apperr.NotFound
				},
			},
		},
		{
			name:       "Failed",
			req:        "some-random-verification-token",
			wantStatus: http.StatusInternalServerError,
			accountRepo: &mockdb.Account{
				FindVerificationTokenFn: func(context.Context, string) (*model.Verification, error) {
					return &model.Verification{
						Token:  "some-random-token-for-verification",
						UserID: 1,
					}, nil
				},
				DeleteVerificationTokenFn: func(context.Context, *model.Verification) error {
					return apperr.DB
				},
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			authService := auth.NewAuthService(tt.userRepo, tt.accountRepo, tt.jwt, tt.m, tt.mobile)
			service.AuthRouter(authService, r)
			ts := httptest.NewServer(r)
			defer ts.Close()

			path := ts.URL + "/verification/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestSignupMobile(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		userRepo    *mockdb.User
		accountRepo *mockdb.Account
		jwt         *mock.JWT
		m           *mock.Mail
		mobile      *mock.Mobile
	}{
		{
			name:       "Success",
			req:        `{"country_code":"+65","mobile":"91919191"}`,
			wantStatus: http.StatusCreated,
			userRepo: &mockdb.User{
				FindByMobileFn: func(string, string) (*model.User, error) {
					return nil, apperr.DB // no such user, so create
				},
			},
			accountRepo: &mockdb.Account{
				CreateWithMobileFn: func(*model.User) error {
					return nil
				},
			},
			mobile: &mock.Mobile{
				GenerateSMSTokenFn: func(string, string) error {
					return nil
				},
			},
		},
		{
			name:       "Failure: no country code",
			req:        `{"mobile":"91919191}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Failure: no mobile",
			req:        `{"country_code":"+1}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Failure: user with mobile number already exists",
			req:        `{"country_code":"+65","mobile":"91919191"}`,
			wantStatus: http.StatusConflict,
			userRepo: &mockdb.User{
				FindByMobileFn: func(string, string) (*model.User, error) {
					return &model.User{ // user already exists
						CountryCode: "+65",
						Mobile:      "91919191",
					}, nil
				},
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			authService := auth.NewAuthService(tt.userRepo, tt.accountRepo, tt.jwt, tt.m, tt.mobile)
			service.AuthRouter(authService, r)
			ts := httptest.NewServer(r)
			defer ts.Close()
			// signup
			path := ts.URL + "/signup/m"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestVerifyMobile(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		userRepo    *mockdb.User
		accountRepo *mockdb.Account
		jwt         *mock.JWT
		m           *mock.Mail
		mobile      *mock.Mobile
	}{
		{
			name:       "Success",
			req:        `{"country_code":"+65","mobile":"91919191","code":"324567"}`,
			wantStatus: http.StatusOK,
			mobile: &mock.Mobile{
				CheckCodeFn: func(string, string, string) error {
					return nil
				},
			},
			userRepo: &mockdb.User{
				FindByMobileFn: func(string, string) (*model.User, error) {
					return &model.User{
						CountryCode: "+65",
						Mobile:      "91919191",
					}, nil
				},
				UpdateFn: func(*model.User) (*model.User, error) {
					return &model.User{
						CountryCode: "+65",
						Mobile:      "91919191",
						Active:      true,
					}, nil
				},
			},
		},
		{
			name:       "Failure: no country code",
			req:        `{"mobile":"91919191}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Failure: no mobile",
			req:        `{"country_code":"+1}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Failure: code not verified",
			req:        `{"country_code":"+65","mobile":"91919191","code":"324567"}`,
			wantStatus: http.StatusNotFound,
			mobile: &mock.Mobile{
				CheckCodeFn: func(string, string, string) error {
					return apperr.NewStatus(http.StatusNotFound)
				},
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			authService := auth.NewAuthService(tt.userRepo, tt.accountRepo, tt.jwt, tt.m, tt.mobile)
			service.AuthRouter(authService, r)
			ts := httptest.NewServer(r)
			defer ts.Close()
			// signup
			path := ts.URL + "/verifycode"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
