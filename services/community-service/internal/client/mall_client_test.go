package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDebitWalletParsesSharedSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Internal-Token"); got != "test-token" {
			t.Fatalf("X-Internal-Token = %q, want test-token", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"wallet_transaction_id": 42,
				"balance_before":        1000,
				"balance_after":         900,
			},
		})
	}))
	defer server.Close()

	client := NewMallClient(server.URL, "test-token")
	txID, err := client.DebitWallet(1, 100, "community-fee:9", "物业费缴纳", "password", "123456", "")
	if err != nil {
		t.Fatalf("DebitWallet returned error: %v", err)
	}
	if txID != 42 {
		t.Fatalf("wallet transaction id = %d, want 42", txID)
	}
}

func TestDebitWalletRejectsBusinessError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    400,
			"message": "余额不足",
			"data":    nil,
		})
	}))
	defer server.Close()

	client := NewMallClient(server.URL, "test-token")
	if _, err := client.DebitWallet(1, 100, "community-fee:9", "物业费缴纳", "password", "bad", ""); err == nil {
		t.Fatal("DebitWallet should reject non-zero response code")
	}
}
