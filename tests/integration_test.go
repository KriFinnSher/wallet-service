package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func post(t *testing.T, path string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	return resp
}

func get(t *testing.T, path string) *http.Response {
	t.Helper()
	resp, err := http.Get(baseURL + path)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	return resp
}

func decode(t *testing.T, resp *http.Response, dest any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		t.Fatalf("decode error: %v", err)
	}
}

func TestUpdateBalanceAndGet(t *testing.T) {
	walletId := "00000000-0000-0000-0000-000000000001"

	// Пополнение на 100
	resp := post(t, "/api/v1/wallet", map[string]any{
		"walletId":      walletId,
		"operationType": "DEPOSIT",
		"amount":        100,
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on deposit, got %d", resp.StatusCode)
	}
	var depositResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Balance uint64 `json:"balance"`
	}
	decode(t, resp, &depositResp)

	// Снятие 50
	resp = post(t, "/api/v1/wallet", map[string]any{
		"walletId":      walletId,
		"operationType": "WITHDRAW",
		"amount":        50,
	})
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on withdraw, got %d", resp.StatusCode)
	}
	var withdrawResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Balance uint64 `json:"balance"`
	}
	decode(t, resp, &withdrawResp)

	// GET для проверки текущего баланса
	resp = get(t, "/api/v1/wallets/"+walletId)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 on get wallet, got %d", resp.StatusCode)
	}
}
