package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type MailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	UseTLS      bool
}

// SendMail sends an HTML email using SMTP over SSL (port 465) or standard SMTP.
func SendMail(cfg MailConfig, to, subject, bodyHtml string) error {
	if cfg.Host == "" || cfg.Port == 0 || cfg.Username == "" || cfg.Password == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", cfg.FromName, cfg.FromAddress)
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + bodyHtml

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// Since QQ Mail uses SSL on port 465, we must use tls.Dial
	if cfg.Port == 465 {
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         cfg.Host,
		})
		if err != nil {
			return fmt.Errorf("failed to dial SMTP over TLS: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, cfg.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}

		if err := client.Mail(cfg.FromAddress); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to open data writer: %w", err)
		}

		if _, err := w.Write([]byte(message)); err != nil {
			return fmt.Errorf("failed to write message body: %w", err)
		}

		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}

		return client.Quit()
	}

	// Fallback to normal SendMail for port 587/25
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	err := smtp.SendMail(addr, auth, cfg.FromAddress, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send normal mail: %w", err)
	}

	return nil
}

type OrderItemInfo struct {
	ProductName string
	PriceCents  int64
	Quantity    int
}

func GenerateOrderPaidHTML(username, orderNo string, totalAmountCents int64, usedPoints int, actualPaidCents int64, paymentMethod string, paidAt string, items []OrderItemInfo) string {
	amountYuan := fmt.Sprintf("%.2f", float64(totalAmountCents)/100.0)
	discountYuan := fmt.Sprintf("%.2f", float64(usedPoints*10)/100.0)
	actualPaidYuan := fmt.Sprintf("%.2f", float64(actualPaidCents)/100.0)

	itemsHtml := ""
	for _, item := range items {
		itemTotal := fmt.Sprintf("%.2f", float64(item.PriceCents*int64(item.Quantity))/100.0)
		itemPrice := fmt.Sprintf("%.2f", float64(item.PriceCents)/100.0)
		itemsHtml += fmt.Sprintf(`
			<tr>
				<td style="padding: 12px 10px; border-bottom: 1px solid #f1f5f9; color: #334155; font-size: 14px;">%s</td>
				<td style="padding: 12px 10px; border-bottom: 1px solid #f1f5f9; color: #64748b; font-size: 14px; text-align: center;">¥%s</td>
				<td style="padding: 12px 10px; border-bottom: 1px solid #f1f5f9; color: #64748b; font-size: 14px; text-align: center;">%d</td>
				<td style="padding: 12px 10px; border-bottom: 1px solid #f1f5f9; color: #0f172a; font-size: 14px; text-align: right; font-weight: 600;">¥%s</td>
			</tr>
		`, item.ProductName, itemPrice, item.Quantity, itemTotal)
	}

	pointsRow := ""
	if usedPoints > 0 {
		pointsRow = fmt.Sprintf(`
			<div class="details-row">
				<span class="label">积分抵扣</span>
				<span class="value" style="color: #ef4444;">-¥%s 元 (%d 积分)</span>
			</div>
		`, discountYuan, usedPoints)
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f4f6f8; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 30px auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 15px rgba(0,0,0,0.05); overflow: hidden; border: 1px solid #e1e8ed; }
        .header { background: linear-gradient(135deg, #4F46E5, #3B82F6); padding: 35px 20px; text-align: center; color: #ffffff; }
        .header h1 { margin: 0; font-size: 24px; font-weight: 600; letter-spacing: 0.5px; }
        .content { padding: 35px 25px; color: #334155; line-height: 1.6; }
        .content p { margin: 0 0 20px 0; font-size: 16px; }
        .details-box { background-color: #f8fafc; border-radius: 6px; padding: 20px; margin-bottom: 25px; border: 1px solid #f1f5f9; }
        .details-row { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 14px; }
        .details-row:last-child { margin-bottom: 0; }
        .label { color: #64748b; font-weight: 500; }
        .value { color: #0f172a; font-weight: 600; text-align: right; }
        .price { color: #4F46E5 !important; font-size: 18px; }
        .items-table { width: 100%%; border-collapse: collapse; margin-top: 15px; margin-bottom: 25px; }
        .items-table th { background-color: #f8fafc; color: #64748b; font-weight: 600; font-size: 13px; text-transform: uppercase; padding: 12px 10px; text-align: left; border-bottom: 2px solid #e2e8f0; }
        .footer { background-color: #f8fafc; padding: 20px; text-align: center; color: #94a3b8; font-size: 12px; border-top: 1px solid #f1f5f9; }
        .footer p { margin: 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 订单支付成功</h1>
        </div>
        <div class="content">
            <p>尊敬的 %s 用户：</p>
            <p>您的订单已成功支付并已开始处理，感谢您在智能社区商城的消费！</p>
            
            <div style="font-weight: 600; font-size: 16px; margin-bottom: 10px; color: #1e293b;">📦 订单商品明细</div>
            <table class="items-table">
                <thead>
                    <tr>
                        <th style="text-align: left;">商品名称</th>
                        <th style="text-align: center; width: 80px;">单价</th>
                        <th style="text-align: center; width: 60px;">数量</th>
                        <th style="text-align: right; width: 90px;">小计</th>
                    </tr>
                </thead>
                <tbody>
                    %s
                </tbody>
            </table>

            <div style="font-weight: 600; font-size: 16px; margin-bottom: 10px; color: #1e293b;">💳 支付账单</div>
            <div class="details-box">
                <div class="details-row">
                    <span class="label">订单编号</span>
                    <span class="value">%s</span>
                </div>
                <div class="details-row">
                    <span class="label">支付方式</span>
                    <span class="value" style="color: #059669;">%s</span>
                </div>
                <div class="details-row">
                    <span class="label">订单总额</span>
                    <span class="value">¥%s 元</span>
                </div>
                %s
                <div class="details-row" style="border-top: 1px dashed #e2e8f0; padding-top: 12px; margin-top: 12px;">
                    <span class="label" style="font-size: 15px; font-weight: 600; color: #1e293b;">实际支付金额</span>
                    <span class="value price" style="font-size: 20px; font-weight: 700; color: #4F46E5;">¥%s 元</span>
                </div>
                <div class="details-row">
                    <span class="label">支付时间</span>
                    <span class="value">%s</span>
                </div>
            </div>
            <p style="margin-top: 20px; font-size: 14px; color: #64748b;">如果您对本订单有任何疑问，请联系社区客服协助处理。</p>
        </div>
        <div class="footer">
            <p>© 2026 智能社区管理服务中心. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
	`, username, itemsHtml, orderNo, paymentMethod, amountYuan, pointsRow, actualPaidYuan, paidAt)
}

func GenerateWalletRechargedHTML(username string, amountCents int64, paymentMethod string, rechargedAt, idempotencyKey string) string {
	amountYuan := fmt.Sprintf("%.2f", float64(amountCents)/100.0)
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f4f6f8; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 30px auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 15px rgba(0,0,0,0.05); overflow: hidden; border: 1px solid #e1e8ed; }
        .header { background: linear-gradient(135deg, #10B981, #059669); padding: 35px 20px; text-align: center; color: #ffffff; }
        .header h1 { margin: 0; font-size: 24px; font-weight: 600; letter-spacing: 0.5px; }
        .content { padding: 35px 25px; color: #334155; line-height: 1.6; }
        .content p { margin: 0 0 20px 0; font-size: 16px; }
        .details-box { background-color: #f8fafc; border-radius: 6px; padding: 20px; margin-bottom: 25px; border: 1px solid #f1f5f9; }
        .details-row { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 14px; }
        .details-row:last-child { margin-bottom: 0; }
        .label { color: #64748b; font-weight: 500; }
        .value { color: #0f172a; font-weight: 600; text-align: right; }
        .price { color: #10b981 !important; font-size: 20px; font-weight: 700; }
        .footer { background-color: #f8fafc; padding: 20px; text-align: center; color: #94a3b8; font-size: 12px; border-top: 1px solid #f1f5f9; }
        .footer p { margin: 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💰 电子钱包充值成功</h1>
        </div>
        <div class="content">
            <p>尊敬的 %s 用户：</p>
            <p>您的个人电子钱包余额已充值成功，以下是您的充值及支付账单明细：</p>
            
            <div class="details-box">
                <div class="details-row">
                    <span class="label">充值到账金额</span>
                    <span class="value price">+¥%s 元</span>
                </div>
                <div class="details-row">
                    <span class="label">支付渠道</span>
                    <span class="value" style="color: #059669; font-weight: 600;">%s</span>
                </div>
                <div class="details-row" style="border-top: 1px dashed #e2e8f0; padding-top: 12px; margin-top: 12px;">
                    <span class="label" style="font-weight: 600; color: #1e293b;">支付款项金额</span>
                    <span class="value" style="font-size: 16px; font-weight: 700; color: #1e293b;">¥%s 元</span>
                </div>
                <div class="details-row">
                    <span class="label">流水单号</span>
                    <span class="value" style="font-family: monospace; font-size: 12px; color: #64748b;">%s</span>
                </div>
                <div class="details-row">
                    <span class="label">到账时间</span>
                    <span class="value">%s</span>
                </div>
            </div>
            
            <p style="font-size: 14px; color: #64748b;">充值款项已实时存入您的余额。如有任何疑问，请通过个人中心账单或咨询社区物业前台。</p>
        </div>
        <div class="footer">
            <p>© 2026 智能社区管理服务中心. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
	`, username, amountYuan, paymentMethod, amountYuan, idempotencyKey, rechargedAt)
}

func GeneratePropertyFeePaidHTML(username string, billingMonth string, amountCents int64, paymentMethod string, paidAt string, transactionNo string) string {
	amountYuan := fmt.Sprintf("%.2f", float64(amountCents)/100.0)
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f4f6f8; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 30px auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 15px rgba(0,0,0,0.05); overflow: hidden; border: 1px solid #e1e8ed; }
        .header { background: linear-gradient(135deg, #F59E0B, #D97706); padding: 35px 20px; text-align: center; color: #ffffff; }
        .header h1 { margin: 0; font-size: 24px; font-weight: 600; letter-spacing: 0.5px; }
        .content { padding: 35px 25px; color: #334155; line-height: 1.6; }
        .content p { margin: 0 0 20px 0; font-size: 16px; }
        .details-box { background-color: #f8fafc; border-radius: 6px; padding: 20px; margin-bottom: 25px; border: 1px solid #f1f5f9; }
        .details-row { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 14px; }
        .details-row:last-child { margin-bottom: 0; }
        .label { color: #64748b; font-weight: 500; }
        .value { color: #0f172a; font-weight: 600; text-align: right; }
        .price { color: #D97706 !important; font-size: 20px; font-weight: 700; }
        .footer { background-color: #f8fafc; padding: 20px; text-align: center; color: #94a3b8; font-size: 12px; border-top: 1px solid #f1f5f9; }
        .footer p { margin: 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏠 物业缴费成功通知</h1>
        </div>
        <div class="content">
            <p>尊敬的 %s 业主：</p>
            <p>您的社区物业服务费已缴费成功，以下是您的缴费账单明细：</p>
            
            <div class="details-box">
                <div class="details-row">
                    <span class="label">缴费项目</span>
                    <span class="value">物业服务费</span>
                </div>
                <div class="details-row">
                    <span class="label">账单月份</span>
                    <span class="value">%s</span>
                </div>
                <div class="details-row">
                    <span class="label">支付方式</span>
                    <span class="value" style="color: #d97706; font-weight: 600;">%s</span>
                </div>
                <div class="details-row" style="border-top: 1px dashed #e2e8f0; padding-top: 12px; margin-top: 12px;">
                    <span class="label" style="font-weight: 600; color: #1e293b;">已缴金额</span>
                    <span class="value price">¥%s 元</span>
                </div>
                <div class="details-row">
                    <span class="label">交易流水号</span>
                    <span class="value" style="font-family: monospace; font-size: 12px; color: #64748b;">%s</span>
                </div>
                <div class="details-row">
                    <span class="label">缴费时间</span>
                    <span class="value">%s</span>
                </div>
            </div>
            
            <p style="font-size: 14px; color: #64748b;">感谢您对社区建设的大力支持！如有任何疑问，请联系社区物业服务中心前台协助解决。</p>
        </div>
        <div class="footer">
            <p>© 2026 智能社区管理服务中心. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
	`, username, billingMonth, paymentMethod, amountYuan, transactionNo, paidAt)
}

