package services

import (
	"fmt"
	"net/http"
	"tiktok_tools/apperr"
	"tiktok_tools/model/tiktok"

	"github.com/go-pg/pg/v10/orm"
	"go.uber.org/zap"
)

// NewTiktokShopRepo returns a new TiktokShopRepo instance
func NewTiktokShopRepo(db orm.DB, log *zap.Logger) *TiktokShopRepo {
	return &TiktokShopRepo{db, log}
}

// TiktokShopRepo is the client for our user model
type TiktokShopRepo struct {
	db  orm.DB
	log *zap.Logger
}

// Create creates a new tiktok shop in our database
func (a *TiktokShopRepo) Create(s *tiktok.TiktokShop) (*tiktok.TiktokShop, error) {
	shop := new(tiktok.TiktokShop)
	sql := `SELECT id FROM tiktok_shops WHERE open_id = ? and deleted_at IS NULL`
	res, err := a.db.Query(shop, sql, s.OpenID)
	if err != nil {
		fmt.Println(err.Error())
		a.log.Error("TiktokShopRepo Error: ", zap.Error(err))
		return nil, apperr.DB
	}
	if res.RowsReturned() != 0 {
		fmt.Println("user exists in database")
		return nil, apperr.New(http.StatusBadRequest, "TiktokShop already exists.")
	}
	if _, err := a.db.Model(s).Insert(); err != nil {
		a.log.Warn("TiktokShopRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return s, nil
}

// DeleteVerificationToken sets deleted_at for an existing verification token
func (a *TiktokShopRepo) View(openID string) (*tiktok.TiktokShop, error) {
	var s = new(tiktok.TiktokShop)
	sql := `SELECT * FROM tiktok_shops WHERE (open_id = ? and deleted_at IS NULL)`
	_, err := a.db.QueryOne(s, sql, openID)
	if err != nil {
		a.log.Warn("TiktokShopRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.NotFound
	}
	return s, nil
}

// DeleteVerificationToken sets deleted_at for an existing verification token
func (a *TiktokShopRepo) RefreshToken(s *tiktok.TiktokShop) error {
	_, err := a.db.Model().Column("refresh_token", "refresh_token_expire_in", "update_at").WherePK().Update()
	if err != nil {
		a.log.Warn("TiktokShopRepo Error", zap.Error(err))
		return apperr.DB
	}
	return err
}
