package models

import (
	"errors"
	"fmt"
	"time"
)

type Order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          *Delivery `json:"delivery"`
	Payment           *Payment  `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shard_key"`
	SmId              int       `json:"sm_id"`
	DateCreated       string    `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func (o *Order) Validate() error {
	if o.OrderUid == "" {
		return errors.New("order_uid is empty")
	}
	if o.TrackNumber == "" {
		return errors.New("track_number is empty")
	}
	if o.Entry == "" {
		return errors.New("entry is empty")
	}
	if o.Locale == "" {
		return errors.New("locale is empty")
	}
	if o.CustomerId == "" {
		return errors.New("customer_id is empty")
	}
	if o.DeliveryService == "" {
		return errors.New("delivery_service is empty")
	}
	if o.ShardKey == "" {
		return errors.New("shard_key is empty")
	}
	if o.OofShard == "" {
		return errors.New("oof_shard is empty")
	}
	if o.SmId <= 0 {
		return errors.New("sm_id must be positive")
	}
	if o.DateCreated != "" {
		_, err := time.Parse(time.RFC3339, o.DateCreated)
		if err != nil {
			return fmt.Errorf("invalid date_created format: %v", err)
		}
	}
	if o.Delivery == nil {
		return errors.New("delivery is nil")
	}
	if err := o.Delivery.Validate(); err != nil {
		return fmt.Errorf("delivery validation failed: %v", err)
	}

	if o.Payment == nil {
		return errors.New("payment is nil")
	}
	if err := o.Payment.Validate(); err != nil {
		return fmt.Errorf("payment validation failed: %v", err)
	}
	if len(o.Items) == 0 {
		return errors.New("items array is empty")
	}
	for i, item := range o.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item[%d] validation failed: %v", i, err)
		}
	}
	return nil
}
