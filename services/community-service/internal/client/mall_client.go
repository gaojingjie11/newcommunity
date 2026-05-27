package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MallClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewMallClient(baseURL, token string) *MallClient {
	return &MallClient{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type debitWalletRequest struct {
	UserID         int64  `json:"user_id"`
	Amount         int64  `json:"amount"`
	BizType        string `json:"biz_type"`
	BizID          string `json:"biz_id"`
	IdempotencyKey string `json:"idempotency_key"`
	Remark         string `json:"remark"`
	PayType        string `json:"pay_type"`
	Password       string `json:"password"`
	FaceImageURL   string `json:"face_image_url"`
}

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type debitWalletData struct {
	WalletTransactionID int64 `json:"wallet_transaction_id"`
	BalanceBefore       int64 `json:"balance_before"`
	BalanceAfter        int64 `json:"balance_after"`
}

// DebitWallet calls mall-service internal API to debit from a user's wallet.
// Returns the wallet transaction ID on success.
func (c *MallClient) DebitWallet(userID, amount int64, idempotencyKey, remark, payType, password, faceImageURL string) (int64, error) {
	reqBody := debitWalletRequest{
		UserID:         userID,
		Amount:         amount,
		BizType:        "property_fee",
		BizID:          idempotencyKey,
		IdempotencyKey: idempotencyKey,
		Remark:         remark,
		PayType:        payType,
		Password:       password,
		FaceImageURL:   faceImageURL,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/api/internal/mall/wallet/debit"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("call mall-service: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read response: %w", err)
	}
	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			return 0, fmt.Errorf("wallet debit http status %d", resp.StatusCode)
		}
		return 0, fmt.Errorf("parse response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if apiResp.Message != "" {
			return 0, fmt.Errorf("%s", apiResp.Message)
		}
		return 0, fmt.Errorf("wallet debit http status %d", resp.StatusCode)
	}

	if apiResp.Code != 0 {
		return 0, fmt.Errorf("%s", apiResp.Message)
	}

	var data debitWalletData
	if err := json.Unmarshal(apiResp.Data, &data); err != nil {
		return 0, fmt.Errorf("parse debit data: %w", err)
	}

	return data.WalletTransactionID, nil
}
