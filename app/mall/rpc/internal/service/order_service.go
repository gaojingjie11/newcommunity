package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderService struct {
	db               *gorm.DB
	orderRepo        *repository.OrderRepo
	cartRepo         *repository.CartRepo
	productRepo      *repository.ProductRepo
	storeRepo        *repository.StoreRepo
	storeProductRepo *repository.StoreProductRepo
	walletRepo       *repository.WalletRepo
	eventBus         *EventBus
	rdb              *goredis.Client
}

func NewOrderService(
	db *gorm.DB,
	orderRepo *repository.OrderRepo,
	cartRepo *repository.CartRepo,
	productRepo *repository.ProductRepo,
	storeRepo *repository.StoreRepo,
	storeProductRepo *repository.StoreProductRepo,
	walletRepo *repository.WalletRepo,
	eventBus *EventBus,
	rdb *goredis.Client,
) *OrderService {
	return &OrderService{
		db:               db,
		orderRepo:        orderRepo,
		cartRepo:         cartRepo,
		productRepo:      productRepo,
		storeRepo:        storeRepo,
		storeProductRepo: storeProductRepo,
		walletRepo:       walletRepo,
		eventBus:         eventBus,
		rdb:              rdb,
	}
}

type CreateOrderRequest struct {
	CartIDs []int64 `json:"cart_ids" binding:"required"`
	StoreID int64   `json:"store_id" binding:"required"`
}

type pendingOrderCache struct {
	OrderID     int64     `json:"order_id"`
	OrderNo     string    `json:"order_no"`
	UserID      int64     `json:"user_id"`
	StoreID     int64     `json:"store_id"`
	CartIDs     []int64   `json:"cart_ids"`
	TotalAmount int64     `json:"total_amount"`
	ExpireAt    time.Time `json:"expire_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func pendingOrderRedisKey(orderID int64) string {
	return fmt.Sprintf("mall:order:pending:%d", orderID)
}

type AvailableStoreProduct struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	NeedQty     int     `json:"need_qty"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}

type AvailableStore struct {
	model.Store
	Products []AvailableStoreProduct `json:"products"`
}

func (s *OrderService) ListAvailableStores(userID int64, cartIDs []int64) ([]AvailableStore, error) {
	if len(cartIDs) == 0 {
		return nil, errors.New("购物车项不能为空")
	}
	carts, err := s.cartRepo.FindByIDs(cartIDs, userID)
	if err != nil {
		return nil, err
	}
	if len(carts) != len(cartIDs) {
		return nil, errors.New("部分购物车项不存在")
	}

	neededQty := make(map[int64]int, len(carts))
	productNames := make(map[int64]string, len(carts))
	productPrices := make(map[int64]float64, len(carts))
	productIDs := make([]int64, 0, len(carts))
	seenProducts := make(map[int64]struct{}, len(carts))
	for _, cart := range carts {
		neededQty[cart.ProductID] += cart.Quantity
		productNames[cart.ProductID] = cart.Product.Name
		productPrices[cart.ProductID] = float64(cart.Product.Price) / 100
		if _, ok := seenProducts[cart.ProductID]; !ok {
			productIDs = append(productIDs, cart.ProductID)
			seenProducts[cart.ProductID] = struct{}{}
		}
	}

	var common map[int64]struct{}
	for index, productID := range productIDs {
		storeIDs, err := s.storeProductRepo.ListAvailableStoreIDs(productID, neededQty[productID])
		if err != nil {
			return nil, err
		}
		current := make(map[int64]struct{}, len(storeIDs))
		for _, storeID := range storeIDs {
			current[storeID] = struct{}{}
		}
		if index == 0 {
			common = current
			continue
		}
		for storeID := range common {
			if _, ok := current[storeID]; !ok {
				delete(common, storeID)
			}
		}
	}

	storeIDs := make([]int64, 0, len(common))
	for storeID := range common {
		storeIDs = append(storeIDs, storeID)
	}
	stores, err := s.storeRepo.ListByIDs(storeIDs)
	if err != nil {
		return nil, err
	}

	details, err := s.storeProductRepo.ListAvailableDetails(storeIDs, productIDs)
	if err != nil {
		return nil, err
	}
	stockByStoreProduct := make(map[int64]map[int64]int, len(storeIDs))
	for _, detail := range details {
		if stockByStoreProduct[detail.StoreID] == nil {
			stockByStoreProduct[detail.StoreID] = make(map[int64]int)
		}
		stockByStoreProduct[detail.StoreID][detail.ProductID] = detail.Stock
	}

	result := make([]AvailableStore, 0, len(stores))
	for _, store := range stores {
		item := AvailableStore{Store: store}
		for _, productID := range productIDs {
			item.Products = append(item.Products, AvailableStoreProduct{
				ProductID:   productID,
				ProductName: productNames[productID],
				NeedQty:     neededQty[productID],
				Stock:       stockByStoreProduct[store.ID][productID],
				Price:       productPrices[productID],
			})
		}
		result = append(result, item)
	}
	return result, nil
}

// CreateOrder creates an order with atomic stock deduction.
// All mutations run inside a single DB transaction using WithTx.
func (s *OrderService) CreateOrder(userID int64, req CreateOrderRequest) (*model.Order, error) {
	if len(req.CartIDs) == 0 {
		return nil, errors.New("购物车项不能为空")
	}

	carts, err := s.cartRepo.FindByIDs(req.CartIDs, userID)
	if err != nil {
		return nil, err
	}
	if len(carts) != len(req.CartIDs) {
		return nil, errors.New("部分购物车项不存在")
	}

	var order *model.Order
	err = s.db.Transaction(func(tx *gorm.DB) error {
		orderRepo := s.orderRepo.WithTx(tx)
		storeProductRepo := s.storeProductRepo.WithTx(tx)
		cartRepo := s.cartRepo.WithTx(tx)
		productRepo := s.productRepo.WithTx(tx)

		var totalAmount int64
		var items []model.OrderItem

		for _, cart := range carts {
			affected, err := storeProductRepo.DeductStock(tx, req.StoreID, cart.ProductID, cart.Quantity)
			if err != nil {
				return err
			}
			if affected == 0 {
				return fmt.Errorf("商品「%s」在该门店库存不足", cart.Product.Name)
			}

			affectedProd, err := productRepo.DeductStock(tx, cart.ProductID, cart.Quantity)
			if err != nil {
				return err
			}
			if affectedProd == 0 {
				return fmt.Errorf("商品「%s」商城总库存不足", cart.Product.Name)
			}

			totalAmount += cart.Product.Price * int64(cart.Quantity)
			items = append(items, model.OrderItem{
				StoreID:   req.StoreID,
				ProductID: cart.ProductID,
				Price:     cart.Product.Price,
				Quantity:  cart.Quantity,
			})
		}

		now := time.Now()
		expireAt := now.Add(time.Duration(consts.OrderExpireDuration) * time.Minute)
		orderNo := fmt.Sprintf("%d%d", now.UnixMilli(), userID)
		order = &model.Order{
			OrderNo:     orderNo,
			UserID:      userID,
			StoreID:     req.StoreID,
			TotalAmount: totalAmount,
			Status:      consts.OrderStatusPendingPayment,
			ExpireAt:    &expireAt,
			Items:       items,
		}

		if err := orderRepo.CreateOrder(tx, order); err != nil {
			return err
		}

		return cartRepo.DeleteByIDs(req.CartIDs, userID)
	})

	// Best-effort event publish after tx commits
	if err == nil && s.eventBus != nil {
		s.eventBus.PublishOrderCreated(order)
	}
	if err == nil {
		s.cachePendingOrder(order, req.CartIDs)
	}

	return order, err
}

func (s *OrderService) cachePendingOrder(order *model.Order, cartIDs []int64) {
	if s.rdb == nil || order == nil || order.ExpireAt == nil {
		return
	}
	ttl := time.Until(*order.ExpireAt)
	if ttl <= 0 {
		return
	}
	body, err := json.Marshal(pendingOrderCache{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		UserID:      order.UserID,
		StoreID:     order.StoreID,
		CartIDs:     cartIDs,
		TotalAmount: order.TotalAmount,
		ExpireAt:    *order.ExpireAt,
		CreatedAt:   order.CreatedAt,
	})
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_ = s.rdb.Set(ctx, pendingOrderRedisKey(order.ID), body, ttl).Err()
}

func (s *OrderService) deletePendingOrderCache(orderID int64) {
	if s.rdb == nil || orderID <= 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_ = s.rdb.Del(ctx, pendingOrderRedisKey(orderID)).Err()
}

func (s *OrderService) GetOrder(id, userID int64) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, errors.New("无权查看此订单")
	}
	return order, nil
}

func (s *OrderService) ListOrders(userID int64, page, size int, status *int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.orderRepo.ListByUser(userID, page, size, status)
}

func (s *OrderService) AdminListOrders(page, size int, status *int, keyword string) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.orderRepo.ListAll(page, size, status, keyword)
}

// ShipOrder transitions order from paid to shipped (status 1→2).
func (s *OrderService) ShipOrder(orderID int64) error {
	affected, err := s.orderRepo.UpdateStatus(nil, orderID, consts.OrderStatusPaid, consts.OrderStatusShipped)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("订单状态不允许发货")
	}
	return nil
}

// ReceiveOrder transitions order from shipped to completed (status 2→3).
func (s *OrderService) ReceiveOrder(orderID int64) error {
	affected, err := s.orderRepo.UpdateStatus(nil, orderID, consts.OrderStatusShipped, consts.OrderStatusCompleted)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("订单状态不允许确认收货")
	}
	return nil
}

// CancelOrder cancels a pending-payment order and restores locked store stock.
func (s *OrderService) CancelOrder(orderID int64, reason string) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.Status != consts.OrderStatusPendingPayment {
		return errors.New("当前订单状态不允许取消")
	}
	return s.cancelPendingOrder(order, reason)
}

func (s *OrderService) cancelPendingOrder(order *model.Order, reason string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		orderRepo := s.orderRepo.WithTx(tx)
		storeProductRepo := s.storeProductRepo.WithTx(tx)
		productRepo := s.productRepo.WithTx(tx)

		affected, err := orderRepo.MarkAsCancelled(tx, order.ID, consts.OrderStatusPendingPayment, reason)
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("订单状态已变更，请刷新后重试")
		}

		for _, item := range order.Items {
			if err := storeProductRepo.RestoreStock(tx, order.StoreID, item.ProductID, item.Quantity); err != nil {
				return err
			}
			if _, err := productRepo.RestoreStock(tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		return nil
	})

	if err == nil && s.eventBus != nil {
		s.eventBus.PublishOrderCancelled(order, reason)
	}
	if err == nil {
		s.deletePendingOrderCache(order.ID)
	}
	return err
}
