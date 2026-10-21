package main

import "testing"

func TestApprovePayment(t *testing.T) {
	tests := []struct {
		name string
		e    PaymentEvent
		want bool
	}{
		{"ordinary usd", PaymentEvent{AmountCents: 125000, Currency: "USD", Beneficiary: "Acme"}, true},
		{"over limit", PaymentEvent{AmountCents: 5000001, Currency: "USD", Beneficiary: "Acme"}, false},
		{"missing beneficiary", PaymentEvent{AmountCents: 125000, Currency: "USD"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := approvePayment(tt.e); got != tt.want {
				t.Fatalf("approvePayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
