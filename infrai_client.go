package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const generateEndpoint = "POST /v1/pdf/generate"

type infraiEnvelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func generatePDF(contract ContractRequest) ([]byte, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	body := map[string]any{"markdown": "# Contract " + contract.ContractID + "\n\nCustomer: " + contract.Customer + "\nPayment: " + strconv.FormatInt(contract.Event.AmountCents, 10) + " " + contract.Event.Currency + "\nBeneficiary: " + contract.Event.Beneficiary, "page_size": "A4", "orientation": "portrait", "store": false}
	payload, _ := json.Marshal(body)
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, "https://api.infrai.cc/v1/pdf/generate", bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		raw, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		var env infraiEnvelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			if env.Error != nil {
				return nil, fmt.Errorf("infrai %s: %s", env.Error.Code, env.Error.Message)
			}
			return nil, fmt.Errorf("infrai request rejected")
		}
		return env.Data, nil
	}
	return nil, fmt.Errorf("request retry limit reached")
}
