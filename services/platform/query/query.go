package query

import (
	"tiktok_tools/apperr"
	"tiktok_tools/model"
)

// List prepares data for list queries
func List(u *model.AuthUser) (*model.ListQuery, error) {
	switch true {
	case int(u.Role) <= 2: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == model.CompanyAdminRole:
		return &model.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	case u.Role == model.LocationAdminRole:
		return &model.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, apperr.Forbidden
	}
}
