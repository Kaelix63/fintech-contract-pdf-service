package main

type PaymentEvent struct {
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	Beneficiary string `json:"beneficiary"`
}

type ContractRequest struct {
	ContractID string       `json:"contract_id"`
	Customer   string       `json:"customer"`
	Event      PaymentEvent `json:"payment"`
}

func approvePayment(e PaymentEvent) bool {
	return e.AmountCents > 0 && e.AmountCents <= 5000000 && e.Currency == "USD" && e.Beneficiary != ""
}
