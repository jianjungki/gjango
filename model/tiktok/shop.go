package tiktok

import "tiktok_tools/model"

func init() {
	model.Register(&TiktokShop{})
}

// TiktokShop model
type TiktokShop struct {
	model.Base
	AccessToken          string `json:"access_token"`
	AccessTokenExpireIn  int    `json:"access_token_expire_in"`
	RefreshToken         string `json:"refresh_token"`
	RefreshTokenExpireIn int    `json:"refresh_token_expire_in"`
	OpenID               string `json:"open_id"`
	SellerName           string `json:"seller_name"`
	SellerBaseRegion     string `json:"seller_base_region"`
	UserType             string `json:"user_type"`
}
