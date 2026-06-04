package logic

import (
	"fmt"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/types/mall"
)

func toProtoProduct(p *model.Product) *mall.ProductInfo {
	if p == nil {
		return nil
	}
	isPromotion := int32(0)
	if p.OriginalPrice > p.Price && p.Price > 0 {
		isPromotion = 1
	}
	return &mall.ProductInfo{
		Id:            p.ID,
		CategoryName:  p.CategoryName,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		OriginalPrice: p.OriginalPrice,
		Stock:         int32(p.Stock),
		ImageUrl:      p.ImageURL,
		IsPromotion:   isPromotion,
		Sales:         int32(p.Sales),
		Status:        int32(p.Status),
		Version:       int32(p.Version),
		CreatedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
		CategoryId:    p.CategoryID,
	}
}

func toProtoStore(s *model.Store) *mall.StoreInfo {
	if s == nil {
		return nil
	}
	return &mall.StoreInfo{
		Id:      s.ID,
		Name:    s.Name,
		Address: s.Address,
		Phone:   s.Phone,
	}
}

func toProtoPromotion(pr *model.Promotion) *mall.PromotionInfo {
	if pr == nil {
		return nil
	}
	return &mall.PromotionInfo{
		Id:        pr.ID,
		Title:     pr.Title,
		Type:      int32(pr.Type),
		StartDate: pr.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:   pr.EndDate.Format("2006-01-02 15:04:05"),
		Status:    int32(pr.Status),
	}
}

func toProtoCategory(c *model.ProductCategory) *mall.CategoryInfo {
	if c == nil {
		return nil
	}
	return &mall.CategoryInfo{
		Id:   c.ID,
		Name: c.Name,
		Icon: c.Icon,
		Sort: int32(c.Sort),
	}
}

func toProtoServiceArea(sa *model.ServiceArea) *mall.ServiceAreaInfo {
	if sa == nil {
		return nil
	}
	return &mall.ServiceAreaInfo{
		Id:     sa.ID,
		Name:   sa.Name,
		Status: int32(sa.Status),
	}
}

func toProtoCartItem(c *model.Cart) *mall.CartItem {
	if c == nil {
		return nil
	}
	return &mall.CartItem{
		Id:           c.ID,
		UserId:       c.UserID,
		ProductId:    c.ProductID,
		Quantity:     int32(c.Quantity),
		CreatedAt:    c.CreatedAt.Format("2006-01-02 15:04:05"),
		ProductName:  c.Product.Name,
		ProductPrice: c.Product.Price,
		ProductImage: c.Product.ImageURL,
	}
}

func toProtoFavorite(f *model.Favorite) *mall.FavoriteInfo {
	if f == nil {
		return nil
	}
	return &mall.FavoriteInfo{
		Id:        f.ID,
		UserId:    f.UserID,
		ProductId: f.ProductID,
		CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
		Product:   toProtoProduct(&f.Product),
	}
}

func toProtoComment(co *model.ProductComment) *mall.CommentInfo {
	if co == nil {
		return nil
	}
	return &mall.CommentInfo{
		Id:        co.ID,
		ProductId: co.ProductID,
		UserId:    co.UserID,
		Username:  co.User.Username,
		Avatar:    co.User.Avatar,
		Content:   co.Content,
		Rating:    int32(co.Rating),
		CreatedAt: co.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toProtoOrder(o *model.Order) *mall.OrderInfo {
	if o == nil {
		return nil
	}
	items := make([]*mall.OrderItemInfo, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, &mall.OrderItemInfo{
			Id:           item.ID,
			OrderId:      item.OrderID,
			ProductId:    item.ProductID,
			ProductName:  item.Product.Name,
			ProductImage: item.Product.ImageURL,
			Price:        item.Price,
			Quantity:     int32(item.Quantity),
		})
	}
	var expireAtStr string
	var expiresIn int64
	if o.ExpireAt != nil {
		expireAtStr = o.ExpireAt.Format("2006-01-02 15:04:05")
		expiresIn = int64(time.Until(*o.ExpireAt).Seconds())
		if expiresIn < 0 {
			expiresIn = 0
		}
	}
	return &mall.OrderInfo{
		Id:               o.ID,
		OrderNo:          o.OrderNo,
		UserId:           o.UserID,
		TotalAmount:      o.TotalAmount,
		Status:           int32(o.Status),
		StoreId:          o.StoreID,
		StoreName:        o.Store.Name,
		CreatedAt:        o.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        o.UpdatedAt.Format("2006-01-02 15:04:05"),
		Items:            items,
		UsedBalance:      o.UsedBalance,
		ExpireAt:         expireAtStr,
		ExpiresInSeconds: expiresIn,
	}
}

func toProtoWalletTx(t *model.WalletTransaction) *mall.WalletTransactionInfo {
	if t == nil {
		return nil
	}
	typeStr := "unknown"
	switch t.Type {
	case 1:
		typeStr = "recharge"
	case 2:
		typeStr = "payment"
	case 3:
		typeStr = "refund"
	case 4:
		typeStr = "transfer_in"
	case 5:
		typeStr = "transfer_out"
	}
	txNo := ""
	if t.IdempotencyKey != nil {
		txNo = *t.IdempotencyKey
	} else {
		txNo = fmt.Sprintf("TX%d", t.ID)
	}
	return &mall.WalletTransactionInfo{
		Id:            t.ID,
		WalletId:      t.UserID,
		TransactionNo: txNo,
		Type:          typeStr,
		Amount:        t.Amount,
		Balance:       t.BalanceAfter,
		Remark:        t.Remark,
		CreatedAt:     t.CreatedAt.Format("2006-01-02 15:04:05"),
		BizType:       t.BizType,
		BizId:         t.BizID,
	}
}

func toProtoStoreProduct(sp *model.StoreProduct) *mall.StoreProductInfo {
	if sp == nil {
		return nil
	}
	productName := ""
	var price int64 = 0
	if sp.Product != nil {
		productName = sp.Product.Name
		price = sp.Product.Price
	}
	return &mall.StoreProductInfo{
		Id:          sp.ID,
		StoreId:     sp.StoreID,
		ProductId:   sp.ProductID,
		ProductName: productName,
		Price:       price,
		Stock:       int32(sp.Stock),
		SoldCount:   int32(sp.SoldCount),
		Status:      int32(sp.Status),
	}
}
