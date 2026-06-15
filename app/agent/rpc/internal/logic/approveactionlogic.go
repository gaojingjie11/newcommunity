package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"smartcommunity-microservices/app/agent/rpc/agent"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"
	"smartcommunity-microservices/app/workorder/rpc/workorderrpc"
	"strings"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ApproveActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveActionLogic {
	return &ApproveActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApproveActionLogic) ApproveAction(in *agent.ApproveActionReq) (*agent.ApproveActionResp, error) {
	// 1. Fetch approval record with pessimistic locking (FOR UPDATE) in a transaction
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		l.Errorf("failed to start transaction: %v", tx.Error)
		return &agent.ApproveActionResp{Code: 500, Message: "启动事务失败"}, nil
	}

	var approval model.AgentActionApproval
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("id = ? AND conversation_id = ? AND user_id = ?", in.ActionId, in.ConversationId, in.UserId).
		First(&approval).Error
	if err != nil {
		tx.Rollback()
		l.Errorf("approval action not found: %v", err)
		return &agent.ApproveActionResp{Code: 404, Message: "审批动作未找到"}, nil
	}

	// Idempotency: if already executed, return success immediately
	if approval.Status == "executed" || approval.Status == "approved" {
		tx.Rollback()
		return &agent.ApproveActionResp{Code: 0, Message: "审批已成功执行（幂等保护）"}, nil
	}

	if approval.Status != "pending" {
		tx.Rollback()
		return &agent.ApproveActionResp{Code: 400, Message: fmt.Sprintf("审批动作当前不可执行（状态: %s）", approval.Status)}, nil
	}

	// 2. Transition status to executing (lock status machine)
	approval.Status = "executing"
	approval.UpdatedAt = time.Now()
	if err := tx.Save(&approval).Error; err != nil {
		tx.Rollback()
		l.Errorf("failed to lock approval state: %v", err)
		return &agent.ApproveActionResp{Code: 500, Message: "锁定审批状态失败"}, nil
	}

	// Commit transaction to release lock but mark status as executing
	if err := tx.Commit().Error; err != nil {
		l.Errorf("failed to commit locking transaction: %v", err)
		return &agent.ApproveActionResp{Code: 500, Message: "提交锁定事务失败"}, nil
	}

	var resultMsg string
	var execErr error
	var payURL string
	var orderID int64

	// 3. Execute real RPC depending on action type
	switch approval.ActionType {
	case "create_order":
		var input CreateOrderInput
		if execErr = json.Unmarshal([]byte(approval.ActionPayload), &input); execErr == nil {
			var addedCartID int64 // Track for cleanup on failure

			// A. Add item to user's cart
			_, execErr = l.svcCtx.MallRpc.AddCartItem(l.ctx, &mall.AddCartItemReq{
				UserId:    in.UserId,
				ProductId: input.ProductID,
				Quantity:  input.Quantity,
			})
			if execErr != nil {
				break
			}

			// B. Fetch user's cart
			var cartResp *mall.CartResp
			cartResp, execErr = l.svcCtx.MallRpc.ListCart(l.ctx, &mall.UserIDReq{
				UserId: in.UserId,
			})
			if execErr != nil {
				break
			}

			var cartID int64
			for _, item := range cartResp.Items {
				if item.ProductId == input.ProductID {
					cartID = item.Id
					break
				}
			}

			if cartID == 0 {
				execErr = errors.New("未能成功在购物车中定位商品")
				break
			}
			addedCartID = cartID

			// C. Find available stores
			var storesResp *mall.StoreListResp
			storesResp, execErr = l.svcCtx.MallRpc.ListAvailableStores(l.ctx, &mall.ListAvailableStoresReq{
				UserId:  in.UserId,
				CartIds: []int64{cartID},
			})
			if execErr != nil {
				// Cleanup: remove the cart item we just added
				_, _ = l.svcCtx.MallRpc.RemoveCartItem(l.ctx, &mall.RemoveCartItemReq{UserId: in.UserId, Id: addedCartID})
				break
			}

			var storeID int64
			if len(storesResp.List) > 0 {
				storeID = storesResp.List[0].Id
			} else {
				storeID = 1
			}

			// D. Create order
			var orderInfo *mall.OrderInfo
			orderInfo, execErr = l.svcCtx.MallRpc.CreateOrder(l.ctx, &mall.CreateOrderReq{
				UserId:  in.UserId,
				CartIds: []int64{cartID},
				StoreId: storeID,
			})
			if execErr != nil {
				// Cleanup: remove the cart item we just added
				_, _ = l.svcCtx.MallRpc.RemoveCartItem(l.ctx, &mall.RemoveCartItemReq{UserId: in.UserId, Id: addedCartID})
				break
			}

			priceYuan := float64(orderInfo.TotalAmount) / 100.0
			resultMsg = fmt.Sprintf("已为您成功下单！\n- 订单ID: **%d**\n- 订单号: `%s`\n- 总金额: **￥%.2f**\n- 状态: **待支付**\n\n您可以开始付款支付此订单。",
				orderInfo.Id, orderInfo.OrderNo, priceYuan)

			orderID = orderInfo.Id
			resultBytes, _ := json.Marshal(map[string]interface{}{
				"order_id": orderInfo.Id,
				"order_no": orderInfo.OrderNo,
			})
			approval.ResultPayload = string(resultBytes)
		}

	case "pay_order":
		var input PayOrderInput
		if execErr = json.Unmarshal([]byte(approval.ActionPayload), &input); execErr == nil {
			payType := in.PayType
			if payType == "" {
				payType = "password"
				if in.FaceImageUrl != "" {
					payType = "face"
				} else if in.PaymentPassword == "" {
					execErr = errors.New("支付验证参数缺失，请输入密码或使用人脸扫描")
					break
				}
			}

			// Use approval.ID as idempotency key to prevent duplicate payments on retry
			idempotencyKey := fmt.Sprintf("agent-pay:%s:%d", approval.ID, input.OrderID)

			var payResp *mall.PayOrderResp
			payResp, execErr = l.svcCtx.MallRpc.PayOrder(l.ctx, &mall.PayOrderReq{
				UserId:         in.UserId,
				Id:             input.OrderID,
				PayType:        payType,
				Password:       in.PaymentPassword,
				FaceImageUrl:   in.FaceImageUrl,
				IdempotencyKey: idempotencyKey,
				ReturnUrl:      in.ReturnUrl,
			})
			if execErr != nil {
				break
			}

			if payType == "alipay" {
				payURL = payResp.PayUrl
				resultMsg = fmt.Sprintf("已成功为您生成支付宝支付链接，正在前往支付收银台付款。\n- 订单ID: **%d**\n- 订单号: `%s`\n- 状态: **待支付(跳转中)**", input.OrderID, payResp.OrderNo)
			} else {
				resultMsg = fmt.Sprintf("订单支付成功！\n- 订单号: `%s`\n- 状态: **已支付**\n已成功扣款，商品正在打包准备发货。", payResp.OrderNo)
			}
		}

	case "submit_repair":
		var input SubmitRepairInput
		if execErr = json.Unmarshal([]byte(approval.ActionPayload), &input); execErr == nil {
			_, execErr = l.svcCtx.WorkorderRpc.CreateWorkorder(l.ctx, &workorderrpc.CreateWorkorderReq{
				UserId:      in.UserId,
				Type:        input.Type,
				Category:    input.Category,
				Description: input.Description,
			})
			if execErr != nil {
				break
			}

			stName := "报修"
			if strings.ToLower(input.Type) == "complaint" {
				stName = "投诉"
			}
			resultMsg = fmt.Sprintf("您的物业**%s工单**已成功提交！\n- 分类: `%s`\n- 详情描述: %s\n\n我们将尽快为您指派师傅处理，请保持电话畅通。",
				stName, input.Category, input.Description)
		}

	default:
		execErr = fmt.Errorf("不支持的审批动作类型: %s", approval.ActionType)
	}

	// 4. Update status and save chat history
	if execErr != nil {
		l.Errorf("ApproveAction execution error: %v", execErr)
		// Revert status to pending so user can retry (simple atomic update, no manual tx needed)
		if err := l.svcCtx.DB.Model(&model.AgentActionApproval{}).
			Where("id = ?", approval.ID).
			Updates(map[string]interface{}{"status": "pending", "updated_at": time.Now()}).Error; err != nil {
			l.Errorf("failed to revert approval status: %v", err)
		}
		return &agent.ApproveActionResp{Code: 500, Message: fmt.Sprintf("执行失败: %v", execErr)}, nil
	}

	txFinish := l.svcCtx.DB.Begin()
	approval.Status = "executed"
	approval.UpdatedAt = time.Now()
	txFinish.Save(&approval)

	// Save executing result message as Assistant reply in the DB
	l.saveChatMessageTx(txFinish, in.ConversationId, in.UserId, "assistant", resultMsg)
	txFinish.Commit()

	return &agent.ApproveActionResp{Code: 0, Message: "审批执行成功", PayUrl: payURL, OrderId: orderID}, nil
}

func (l *ApproveActionLogic) saveChatMessageTx(tx *gorm.DB, convID string, userID int64, role, content string) {
	msg := &model.SysUserChatMessage{
		ID:             uuid.NewString(),
		UserID:         userID,
		ConversationID: convID,
		Role:           role,
		Content:        content,
		CreatedAt:      time.Now(),
	}
	tx.Create(msg)

	tx.Model(&model.SysUserConversation{}).
		Where("id = ? AND user_id = ?", convID, userID).
		Update("updated_at", time.Now())
}

