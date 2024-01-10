package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jianjungki/tiktok"
)

type TikTokClient struct {
	Cli        *tiktok.Client
	token      *tiktok.AccessTokenResponse
	shopList   map[string]tiktok.Shop
	shopIDList []string
}
type ClientConfig struct {
	AppKey    string
	AppSecret string
}

type ProductItem struct {
	Name        string
	Description string
	CategoryID  string
	BrandID     string
	ImageUrls   []string
	VideoUrls   []string

	Price     ProductPrice
	ExtraInfo ProductInfo
}
type ProductPrice struct {
	ComparePrice string
	RetailPrice  string
}

type ProductInfo struct {
	Weight string `json:"weight"`

	Height int `json:"height"`
	Width  int `json:"width"`
	Length int `json:"length"`

	USCertificationDoc string `json:"us_certification_doc"`
	CommodiyDoc        string `json:"commodiy_doc"`
}

func (p *ProductItem) toCreateProduct() tiktok.CreateProductRequest {

	return tiktok.CreateProductRequest{
		ProductName:   p.Name,
		Description:   p.Description,
		CategoryID:    p.CategoryID,
		BrandID:       p.BrandID,
		PackageLength: p.ExtraInfo.Length,
		PackageWidth:  p.ExtraInfo.Width,
		PackageHeight: p.ExtraInfo.Height,
		PackageWeight: p.ExtraInfo.Weight,
	}
}

var client *TikTokClient

func New(cfg ClientConfig) *TikTokClient {
	tkClient, err := tiktok.New(cfg.AppKey, cfg.AppSecret, tiktok.WithLogger(log.Default()))
	if err != nil {
		return nil
	}
	if client == nil {
		client = &TikTokClient{
			Cli: tkClient,
		}
	}
	return client
}

func (client *TikTokClient) RefreshToken(ctx context.Context) error {
	refreshToken, err := client.Cli.RefreshToken(ctx, client.token.RefreshToken)
	if err != nil {
		log.Fatalf("refresh token query failed: %v\n", err)
		return err
	}
	client.token = &refreshToken
	return nil
}

func (client *TikTokClient) GetShopInfo(ctx context.Context) error {
	shopList, err := client.Cli.GetAuthorizedShop(ctx, client.token.AccessToken, client.token.OpenID)
	if err != nil {
		log.Fatalf("query shop data failed: %v\n", err)
		return err
	}
	for _, shopItem := range shopList.Shops {
		client.shopList[shopItem.ShopID] = shopItem
		client.shopIDList = append(client.shopIDList, shopItem.ShopID)
	}
	return nil
}

func (client *TikTokClient) GetShopListFromCache(ctx context.Context) []tiktok.Shop {
	var shopList = make([]tiktok.Shop, 0)
	for _, shopID := range client.shopIDList {
		shopList = append(shopList, client.shopList[shopID])
	}
	return shopList
}

func (client *TikTokClient) GetAccessTokenFromCache(ctx context.Context) *tiktok.AccessTokenResponse {
	return client.token
}

func (client *TikTokClient) GetOrders(ctx context.Context) (*tiktok.OrdersList, error) {
	if len(client.shopIDList) == 0 {
		return nil, errors.New("shopID list is empty")
	}

	orderList, err := client.Cli.GetOrderList(ctx, tiktok.Param{
		AccessToken: client.token.AccessToken,
		ShopID:      client.shopIDList[0],
	}, tiktok.GetOrderListRequest{})

	if err != nil {
		log.Fatalf("refresh token query failed: %v\n", err)
		return nil, err
	}
	return &orderList, nil
}

func (client *TikTokClient) GetProducts(ctx context.Context) (*tiktok.GetProductListData, error) {
	if len(client.shopIDList) == 0 {
		return nil, errors.New("shopID list is empty")
	}

	productList, err := client.Cli.GetProductList(ctx, tiktok.Param{
		AccessToken: client.token.AccessToken,
		ShopID:      client.shopIDList[0],
	}, tiktok.ProductSearchRequest{})

	if err != nil {
		log.Fatalf("refresh token query failed: %v\n", err)
		return nil, err
	}
	return &productList, nil
}

func (client *TikTokClient) AddProducts(ctx context.Context, productData ProductItem) error {
	if len(client.shopIDList) == 0 {
		return errors.New("shopID list is empty")
	}

	productInfo, err := client.Cli.CreateProduct(ctx, tiktok.Param{
		AccessToken: client.token.AccessToken,
		ShopID:      client.shopIDList[0],
	}, productData.toCreateProduct())
	if err != nil {
		log.Fatalf("create product failed: %v\n", err)
		return err
	}

	fmt.Println(productInfo)
	return nil
}
