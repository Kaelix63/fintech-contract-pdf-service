package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func contractHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in ContractRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if !approvePayment(in.Event) {
		http.Error(w, "payment requires review", http.StatusUnprocessableEntity)
		return
	}
	pdf, err := generatePDF(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"contract_id": in.ContractID, "status": "approved", "pdf": json.RawMessage(pdf)})
}

func main() {
	http.HandleFunc("/contracts", contractHandler)
	fmt.Println("contract service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
