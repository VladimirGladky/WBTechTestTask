package models

import "errors"

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

func (p *Payment) Validate() error {
	if p.Transaction == "" {
		return errors.New("transaction is empty")
	}
	if p.RequestId == "" {
		return errors.New("request_id is empty")
	}
	if p.Currency == "" {
		return errors.New("currency is empty")
	}
	if p.Provider == "" {
		return errors.New("provider is empty")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if p.PaymentDt <= 0 {
		return errors.New("payment_dt must be positive")
	}
	if p.Bank == "" {
		return errors.New("bank is empty")
	}
	if p.DeliveryCost <= 0 {
		return errors.New("delivery_cost must be positive")
	}
	if p.GoodsTotal <= 0 {
		return errors.New("goods_total must be positive")
	}
	if p.CustomFee < 0 {
		return errors.New("custom_fee must be non-negative")
	}
	return nil
}
