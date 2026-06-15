package service

import (
	"log/slog"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"

	"gorm.io/gorm"
)

// OrderTimeoutService handles automatic cancellation of expired pending orders.
//
// Two mechanisms:
//   - Primary: RabbitMQ delayed message (not yet implemented — needs DLX setup)
//   - Fallback: Periodic polling every 60 seconds
//
// Both converge to processExpiredOrder which conditionally cancels orders
// that are past their expire_at time.
type OrderTimeoutService struct {
	db               *gorm.DB
	orderRepo        *repository.OrderRepo
	storeProductRepo *repository.StoreProductRepo
	productRepo      *repository.ProductRepo
	eventBus         *EventBus
	log              *slog.Logger
	stopCh           chan struct{}
}

func NewOrderTimeoutService(
	db *gorm.DB,
	orderRepo *repository.OrderRepo,
	storeProductRepo *repository.StoreProductRepo,
	productRepo *repository.ProductRepo,
	eventBus *EventBus,
	log *slog.Logger,
) *OrderTimeoutService {
	return &OrderTimeoutService{
		db:               db,
		orderRepo:        orderRepo,
		storeProductRepo: storeProductRepo,
		productRepo:      productRepo,
		eventBus:         eventBus,
		log:              log,
		stopCh:           make(chan struct{}),
	}
}

// Start begins the service (MQ-trigger mode, no polling).
func (s *OrderTimeoutService) Start() {
	s.log.Info("order timeout service started (MQ-trigger mode)")
}

// Stop signals the service to exit.
func (s *OrderTimeoutService) Stop() {
	s.log.Info("order timeout service stopped")
}

// processExpiredOrder cancels a single expired order within a transaction.
// Restores stock but does NOT refund (payment never happened for pending orders).
func (s *OrderTimeoutService) processExpiredOrder(order *model.Order) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		orderRepo := s.orderRepo.WithTx(tx)
		storeProductRepo := s.storeProductRepo.WithTx(tx)
		productRepo := s.productRepo.WithTx(tx)

		// Conditional cancel: only if still pending_payment
		affected, err := orderRepo.MarkAsCancelled(tx, order.ID, consts.OrderStatusPendingPayment, "订单超时自动取消")
		if err != nil {
			return err
		}
		if affected == 0 {
			// Already cancelled or paid — skip
			return nil
		}

		// Restore store stock
		for _, item := range order.Items {
			if err := storeProductRepo.RestoreStock(tx, order.StoreID, item.ProductID, item.Quantity); err != nil {
				return err
			}
			if _, err := productRepo.RestoreStock(tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}

		s.log.Info("expired order cancelled",
			"order_id", order.ID,
			"order_no", order.OrderNo,
		)

		// Best-effort event
		if s.eventBus != nil {
			s.eventBus.PublishOrderCancelled(order, "订单超时自动取消")
		}

		return nil
	})
}

func (s *OrderTimeoutService) CancelExpiredOrder(orderID int64) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	return s.processExpiredOrder(order)
}
