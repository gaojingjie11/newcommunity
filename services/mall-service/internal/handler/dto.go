package handler

import (
	"math"
	"time"

	"smartcommunity-microservices/services/mall-service/internal/model"
)

func yuanFromCents(cents int64) float64 {
	return math.Round(float64(cents)) / 100
}

func centsFromYuan(yuan float64) int64 {
	return int64(math.Round(yuan * 100))
}

type productPayload struct {
	ID            int64   `json:"id"`
	CategoryName  string  `json:"category_name"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	Stock         int     `json:"stock"`
	ImageURL      string  `json:"image_url"`
	Sales         int     `json:"sales"`
	Status        int     `json:"status"`
	CategoryID    int64   `json:"category_id"`
}

func (p productPayload) toModel(id int64) model.Product {
	original := centsFromYuan(p.OriginalPrice)
	price := centsFromYuan(p.Price)
	return model.Product{
		ID:            id,
		CategoryName:  p.CategoryName,
		Name:          p.Name,
		Description:   p.Description,
		Price:         price,
		OriginalPrice: original,
		Stock:         p.Stock,
		ImageURL:      p.ImageURL,
		Sales:         p.Sales,
		Status:        p.Status,
		CategoryID:    p.CategoryID,
	}
}

type productResponse struct {
	ID            int64     `json:"id"`
	CategoryName  string    `json:"category_name"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	Stock         int       `json:"stock"`
	ImageURL      string    `json:"image_url"`
	IsPromotion   int       `json:"is_promotion"`
	Sales         int       `json:"sales"`
	Status        int       `json:"status"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	CategoryID    int64     `json:"category_id"`
}

func toProductResponse(product model.Product) productResponse {
	isPromotion := 0
	if product.OriginalPrice > product.Price && product.Price > 0 {
		isPromotion = 1
	}
	return productResponse{
		ID:            product.ID,
		CategoryName:  product.CategoryName,
		Name:          product.Name,
		Description:   product.Description,
		Price:         yuanFromCents(product.Price),
		OriginalPrice: yuanFromCents(product.OriginalPrice),
		Stock:         product.Stock,
		ImageURL:      product.ImageURL,
		IsPromotion:   isPromotion,
		Sales:         product.Sales,
		Status:        product.Status,
		Version:       product.Version,
		CreatedAt:     product.CreatedAt,
		CategoryID:    product.CategoryID,
	}
}

func toProductResponses(products []model.Product) []productResponse {
	items := make([]productResponse, 0, len(products))
	for _, product := range products {
		items = append(items, toProductResponse(product))
	}
	return items
}

type cartResponse struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	ProductID int64           `json:"product_id"`
	Quantity  int             `json:"quantity"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Product   productResponse `json:"product"`
}

func toCartResponses(carts []model.Cart) []cartResponse {
	items := make([]cartResponse, 0, len(carts))
	for _, cart := range carts {
		items = append(items, cartResponse{
			ID:        cart.ID,
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
			CreatedAt: cart.CreatedAt,
			UpdatedAt: cart.UpdatedAt,
			Product:   toProductResponse(cart.Product),
		})
	}
	return items
}

type favoriteResponse struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	ProductID int64           `json:"product_id"`
	CreatedAt time.Time       `json:"created_at"`
	Product   productResponse `json:"product"`
}

func toFavoriteResponses(favorites []model.Favorite) []favoriteResponse {
	items := make([]favoriteResponse, 0, len(favorites))
	for _, favorite := range favorites {
		items = append(items, favoriteResponse{
			ID:        favorite.ID,
			UserID:    favorite.UserID,
			ProductID: favorite.ProductID,
			CreatedAt: favorite.CreatedAt,
			Product:   toProductResponse(favorite.Product),
		})
	}
	return items
}

type orderItemResponse struct {
	ID              int64           `json:"id"`
	OrderID         int64           `json:"order_id"`
	StoreID         int64           `json:"store_id"`
	ProductID       int64           `json:"product_id"`
	Price           float64         `json:"price"`
	Quantity        int             `json:"quantity"`
	ProductSnapshot string          `json:"product_snapshot,omitempty"`
	Product         productResponse `json:"product"`
}

type orderResponse struct {
	ID           int64               `json:"id"`
	OrderNo      string              `json:"order_no"`
	UserID       int64               `json:"user_id"`
	StoreID      int64               `json:"store_id"`
	TotalAmount  float64             `json:"total_amount"`
	UsedPoints   int                 `json:"used_points"`
	UsedBalance  float64             `json:"used_balance"`
	Status       int                 `json:"status"`
	ExpireAt     *time.Time          `json:"expire_at"`
	ExpiresIn    int64               `json:"expires_in_seconds"`
	CancelReason string              `json:"cancel_reason,omitempty"`
	CancelledAt  *time.Time          `json:"cancelled_at,omitempty"`
	PaidAt       *time.Time          `json:"paid_at"`
	Version      int                 `json:"version"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	Items        []orderItemResponse `json:"items"`
	Store        model.Store         `json:"store"`
}

func toOrderResponse(order model.Order) orderResponse {
	items := make([]orderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, orderItemResponse{
			ID:              item.ID,
			OrderID:         item.OrderID,
			StoreID:         item.StoreID,
			ProductID:       item.ProductID,
			Price:           yuanFromCents(item.Price),
			Quantity:        item.Quantity,
			ProductSnapshot: item.ProductSnapshot,
			Product:         toProductResponse(item.Product),
		})
	}
	return orderResponse{
		ID:           order.ID,
		OrderNo:      order.OrderNo,
		UserID:       order.UserID,
		StoreID:      order.StoreID,
		TotalAmount:  yuanFromCents(order.TotalAmount),
		UsedPoints:   order.UsedPoints,
		UsedBalance:  yuanFromCents(order.UsedBalance),
		Status:       order.Status,
		ExpireAt:     order.ExpireAt,
		ExpiresIn:    secondsUntil(order.ExpireAt),
		CancelReason: order.CancelReason,
		CancelledAt:  order.CancelledAt,
		PaidAt:       order.PaidAt,
		Version:      order.Version,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		Items:        items,
		Store:        order.Store,
	}
}

func secondsUntil(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	seconds := int64(time.Until(*t).Seconds())
	if seconds < 0 {
		return 0
	}
	return seconds
}

func toOrderResponses(orders []model.Order) []orderResponse {
	items := make([]orderResponse, 0, len(orders))
	for _, order := range orders {
		items = append(items, toOrderResponse(order))
	}
	return items
}

type walletTransactionResponse struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Type          int       `json:"type"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	RelatedID     int64     `json:"related_id"`
	Remark        string    `json:"remark"`
	BizType       string    `json:"biz_type"`
	BizID         string    `json:"biz_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func toWalletTransactionResponses(txs []model.WalletTransaction) []walletTransactionResponse {
	items := make([]walletTransactionResponse, 0, len(txs))
	for _, tx := range txs {
		items = append(items, walletTransactionResponse{
			ID:            tx.ID,
			UserID:        tx.UserID,
			Type:          tx.Type,
			Amount:        yuanFromCents(tx.Amount),
			BalanceBefore: yuanFromCents(tx.BalanceBefore),
			BalanceAfter:  yuanFromCents(tx.BalanceAfter),
			RelatedID:     tx.RelatedID,
			Remark:        tx.Remark,
			BizType:       tx.BizType,
			BizID:         tx.BizID,
			CreatedAt:     tx.CreatedAt,
		})
	}
	return items
}

type storeProductResponse struct {
	ID          int64            `json:"id"`
	StoreID     int64            `json:"store_id"`
	ProductID   int64            `json:"product_id"`
	Stock       int              `json:"stock"`
	LockedStock int              `json:"locked_stock"`
	SoldCount   int              `json:"sold_count"`
	Version     int              `json:"version"`
	Status      int              `json:"status"`
	Product     *productResponse `json:"product,omitempty"`
}

func toStoreProductResponses(items []model.StoreProduct) []storeProductResponse {
	responses := make([]storeProductResponse, 0, len(items))
	for _, item := range items {
		var product *productResponse
		if item.Product != nil {
			p := toProductResponse(*item.Product)
			product = &p
		}
		responses = append(responses, storeProductResponse{
			ID:          item.ID,
			StoreID:     item.StoreID,
			ProductID:   item.ProductID,
			Stock:       item.Stock,
			LockedStock: item.LockedStock,
			SoldCount:   item.SoldCount,
			Version:     item.Version,
			Status:      item.Status,
			Product:     product,
		})
	}
	return responses
}
