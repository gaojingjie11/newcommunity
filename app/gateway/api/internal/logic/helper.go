package logic

import (
	"context"
	"encoding/json"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"
	"smartcommunity-microservices/app/workorder/rpc/workorderrpc"
)

func getUserIDFromCtx(ctx context.Context) int64 {
	val := ctx.Value("user_id")
	if val == nil {
		return 0
	}
	if num, ok := val.(json.Number); ok {
		id, _ := num.Int64()
		return id
	}
	if f, ok := val.(float64); ok {
		return int64(f)
	}
	if i, ok := val.(int64); ok {
		return i
	}
	return 0
}

func toAPIProductInfo(p *mall.ProductInfo) types.ProductInfo {
	if p == nil {
		return types.ProductInfo{}
	}
	return types.ProductInfo{
		Id:            p.Id,
		CategoryName:  p.CategoryName,
		Name:          p.Name,
		Description:   p.Description,
		Price:         float64(p.Price) / 100.0,
		OriginalPrice: float64(p.OriginalPrice) / 100.0,
		Stock:         p.Stock,
		ImageUrl:      p.ImageUrl,
		IsPromotion:   p.IsPromotion,
		Sales:         p.Sales,
		Status:        p.Status,
		Version:       p.Version,
		CreatedAt:     p.CreatedAt,
		CategoryId:    p.CategoryId,
		ViewCount:     p.ViewCount,
	}
}

func toAPIStoreInfo(s *mall.StoreInfo) types.StoreInfo {
	if s == nil {
		return types.StoreInfo{}
	}
	return types.StoreInfo{
		Id:        s.Id,
		Name:      s.Name,
		Address:   s.Address,
		Phone:     s.Phone,
		CreatedAt: s.CreatedAt,
	}
}

func toAPICategoryInfo(c *mall.CategoryInfo) types.CategoryInfo {
	if c == nil {
		return types.CategoryInfo{}
	}
	return types.CategoryInfo{
		Id:          c.Id,
		Name:        c.Name,
		Description: "",
		Sort:        c.Sort,
	}
}

func toAPICartItem(c *mall.CartItem) types.CartItem {
	if c == nil {
		return types.CartItem{}
	}
	return types.CartItem{
		Id:           c.Id,
		UserId:       c.UserId,
		ProductId:    c.ProductId,
		Quantity:     c.Quantity,
		CreatedAt:    c.CreatedAt,
		ProductName:  c.ProductName,
		ProductPrice: float64(c.ProductPrice) / 100.0,
		ProductImage: c.ProductImage,
	}
}

func toAPIFavoriteInfo(f *mall.FavoriteInfo) types.FavoriteInfo {
	if f == nil {
		return types.FavoriteInfo{}
	}
	return types.FavoriteInfo{
		Id:        f.Id,
		UserId:    f.UserId,
		ProductId: f.ProductId,
		CreatedAt: f.CreatedAt,
		Product:   toAPIProductInfo(f.Product),
	}
}

func toAPICommentInfo(c *mall.CommentInfo) types.CommentInfo {
	if c == nil {
		return types.CommentInfo{}
	}
	return types.CommentInfo{
		Id:        c.Id,
		ProductId: c.ProductId,
		UserId:    c.UserId,
		Username:  c.Username,
		Avatar:    c.Avatar,
		Content:   c.Content,
		Rating:    c.Rating,
		CreatedAt: c.CreatedAt,
	}
}

func toAPIOrderItemInfo(item *mall.OrderItemInfo) types.OrderItemInfo {
	if item == nil {
		return types.OrderItemInfo{}
	}
	return types.OrderItemInfo{
		Id:           item.Id,
		OrderId:      item.OrderId,
		ProductId:    item.ProductId,
		ProductName:  item.ProductName,
		ProductImage: item.ProductImage,
		Price:        float64(item.Price) / 100.0,
		Quantity:     item.Quantity,
	}
}

func toAPIOrderInfo(o *mall.OrderInfo) types.OrderInfo {
	if o == nil {
		return types.OrderInfo{}
	}
	items := make([]types.OrderItemInfo, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, toAPIOrderItemInfo(item))
	}
	return types.OrderInfo{
		Id:               o.Id,
		OrderNo:          o.OrderNo,
		UserId:           o.UserId,
		TotalAmount:      float64(o.TotalAmount) / 100.0,
		Status:           o.Status,
		StoreId:          o.StoreId,
		StoreName:        o.StoreName,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
		Items:            items,
		UsedBalance:      float64(o.UsedBalance) / 100.0,
		ExpireAt:         o.ExpireAt,
		ExpiresInSeconds: o.ExpiresInSeconds,
		UsedPoints:       o.UsedPoints,
		StoreAddress:     o.StoreAddress,
		StorePhone:       o.StorePhone,
	}
}

func toAPIWalletTransactionInfo(t *mall.WalletTransactionInfo) types.WalletTransactionInfo {
	if t == nil {
		return types.WalletTransactionInfo{}
	}
	return types.WalletTransactionInfo{
		Id:            t.Id,
		WalletId:      t.WalletId,
		TransactionNo: t.TransactionNo,
		Type:          t.Type,
		Amount:        float64(t.Amount) / 100.0,
		Balance:       float64(t.Balance) / 100.0,
		Remark:        t.Remark,
		CreatedAt:     t.CreatedAt,
		BizType:       t.BizType,
		BizId:         t.BizId,
	}
}

func toAPIStoreProductInfo(sp *mall.StoreProductInfo) types.StoreProductInfo {
	if sp == nil {
		return types.StoreProductInfo{}
	}
	return types.StoreProductInfo{
		Id:          sp.Id,
		StoreId:     sp.StoreId,
		ProductId:   sp.ProductId,
		ProductName: sp.ProductName,
		Price:       float64(sp.Price) / 100.0,
		Stock:       sp.Stock,
		SoldCount:   sp.SoldCount,
		Status:      sp.Status,
	}
}

func toAPINoticeInfo(n *communityrpc.NoticeInfo) types.NoticeInfo {
	if n == nil {
		return types.NoticeInfo{}
	}
	return types.NoticeInfo{
		Id:        n.Id,
		Title:     n.Title,
		Content:   n.Content,
		Publisher: n.Publisher,
		ViewCount: n.ViewCount,
		Status:    n.Status,
		CreatedAt: n.CreatedAt,
	}
}

func toAPIVisitorInfo(v *communityrpc.VisitorInfo) types.VisitorInfo {
	if v == nil {
		return types.VisitorInfo{}
	}
	return types.VisitorInfo{
		Id:           v.Id,
		UserId:       v.UserId,
		VisitorName:  v.VisitorName,
		VisitorPhone: v.VisitorPhone,
		VisitPurpose: v.VisitPurpose,
		ReleaseTime:  v.ReleaseTime,
		ValidDate:    v.ValidDate,
		Status:       v.Status,
		AuditRemark:  v.AuditRemark,
		AuditAt:      v.AuditAt,
		CreatedAt:    v.CreatedAt,
		UserName:     v.UserName,
		UserMobile:   v.UserMobile,
	}
}

func toAPIParkingSpaceInfo(p *communityrpc.ParkingSpaceInfo) types.ParkingSpaceInfo {
	if p == nil {
		return types.ParkingSpaceInfo{}
	}
	return types.ParkingSpaceInfo{
		Id:         p.Id,
		ParkingNo:  p.ParkingNo,
		Status:     p.Status,
		UserId:     p.UserId,
		UserName:   p.UserName,
		UserMobile: p.UserMobile,
		CarPlate:   p.CarPlate,
		BindingId:  p.BindingId,
	}
}

func toAPIPropertyFeeInfo(f *communityrpc.PropertyFeeInfo) types.PropertyFeeInfo {
	if f == nil {
		return types.PropertyFeeInfo{}
	}
	return types.PropertyFeeInfo{
		Id:      f.Id,
		UserId:  f.UserId,
		Month:   f.Month,
		Amount:  f.Amount,
		Status:  f.Status,
		DueDate: f.DueDate,
		PaidAt:  f.PaidAt,
	}
}

func toAPIPropertyFeePaymentInfo(p *communityrpc.PropertyFeePaymentInfo) types.PropertyFeePaymentInfo {
	if p == nil {
		return types.PropertyFeePaymentInfo{}
	}
	return types.PropertyFeePaymentInfo{
		Id:                  p.Id,
		PropertyFeeId:       p.PropertyFeeId,
		UserId:              p.UserId,
		Amount:              p.Amount,
		WalletTransactionId: p.WalletTransactionId,
		IdempotencyKey:      p.IdempotencyKey,
		Status:              p.Status,
		PaidAt:              p.PaidAt,
	}
}

func toAPIMessageInfo(m *communityrpc.MessageInfo) types.MessageInfo {
	if m == nil {
		return types.MessageInfo{}
	}
	return types.MessageInfo{
		Id:        m.Id,
		UserId:    m.UserId,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		Username:  m.Username,
		Avatar:    m.Avatar,
	}
}

func toAPIWorkorderInfo(w *workorderrpc.WorkorderInfo) types.WorkorderInfo {
	if w == nil {
		return types.WorkorderInfo{}
	}
	return types.WorkorderInfo{
		Id:          w.Id,
		Type:        w.Type,
		UserId:      w.UserId,
		Category:    w.Category,
		Description: w.Description,
		Status:      w.Status,
		Result:      w.Result,
		ProcessorId: w.ProcessorId,
		ProcessedAt: w.ProcessedAt,
		CreatedAt:   w.CreatedAt,
		UserName:    w.UserName,
		UserMobile:  w.UserMobile,
	}
}

func toAPIWorkorderLogInfo(wl *workorderrpc.WorkorderLogInfo) types.WorkorderLogInfo {
	if wl == nil {
		return types.WorkorderLogInfo{}
	}
	return types.WorkorderLogInfo{
		Id:         wl.Id,
		TargetType: wl.TargetType,
		TargetId:   wl.TargetId,
		FromStatus: wl.FromStatus,
		ToStatus:   wl.ToStatus,
		OperatorId: wl.OperatorId,
		Action:     wl.Action,
		Remark:     wl.Remark,
		CreatedAt:  wl.CreatedAt,
	}
}

func toAPISalesRankInfo(sr *statsrpc.SalesRankInfo) types.SalesRankInfo {
	if sr == nil {
		return types.SalesRankInfo{}
	}
	return types.SalesRankInfo{
		ProductId:   sr.ProductId,
		ProductName: sr.ProductName,
		TotalSales:  sr.TotalSales,
		TotalAmount: sr.TotalAmount,
	}
}

func toAPIViewRankInfo(vr *statsrpc.ViewRankInfo) types.ViewRankInfo {
	if vr == nil {
		return types.ViewRankInfo{}
	}
	return types.ViewRankInfo{
		ProductId:   vr.ProductId,
		ProductName: vr.ProductName,
		ViewCount:   vr.ViewCount,
		UniqueUsers: vr.UniqueUsers,
	}
}

func toAPIOrderSummaryInfo(os *statsrpc.OrderSummaryInfo) types.OrderSummaryInfo {
	if os == nil {
		return types.OrderSummaryInfo{}
	}
	return types.OrderSummaryInfo{
		Status:      os.Status,
		Count:       os.Count,
		TotalAmount: os.TotalAmount,
	}
}

func toAPIOrderTrendInfo(ot *statsrpc.OrderTrendInfo) types.OrderTrendInfo {
	if ot == nil {
		return types.OrderTrendInfo{}
	}
	return types.OrderTrendInfo{
		Date:   ot.Date,
		Count:  ot.Count,
		Amount: ot.Amount,
	}
}

func toAPIWorkorderSummaryInfo(ws *statsrpc.WorkorderSummaryInfo) types.WorkorderSummaryInfo {
	if ws == nil {
		return types.WorkorderSummaryInfo{}
	}
	return types.WorkorderSummaryInfo{
		Type:   ws.Type,
		Status: ws.Status,
		Count:  ws.Count,
	}
}

func toAPIAIReportInfo(ar *statsrpc.AIReportInfo) types.AIReportInfo {
	if ar == nil {
		return types.AIReportInfo{}
	}
	return types.AIReportInfo{
		Id:                 ar.Id,
		RepairNewCount:     ar.RepairNewCount,
		RepairPendingCount: ar.RepairPendingCount,
		VisitorNewCount:    ar.VisitorNewCount,
		PropertyPaidCount:  ar.PropertyPaidCount,
		PropertyPaidAmount: ar.PropertyPaidAmount,
		ReportSummary:      ar.ReportSummary,
		ReportMarkdown:     ar.ReportMarkdown,
		GeneratedBy:        ar.GeneratedBy,
		CreatedAt:          ar.CreatedAt,
	}
}
