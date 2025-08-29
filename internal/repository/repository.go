package repository

import (
	"WBTechTestTask/internal/models"
	"WBTechTestTask/pkg/suberrors"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepositoryInterface interface {
	GetOrder(id string) (*models.Order, error)
	CreateOrder(order *models.Order) (string, error)
	GetAllOrders() ([]*models.Order, error)
}

type OrderRepository struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewOrderRepository(ctx context.Context, db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		ctx: ctx,
		db:  db,
	}
}

func (o OrderRepository) GetAllOrders() ([]*models.Order, error) {
	var orders []*models.Order

	rows, err := o.db.Query(o.ctx, "SELECT order_uid FROM orders")
	if err != nil {
		return nil, fmt.Errorf("error reading orders: %w", err)
	}

	for rows.Next() {
		var orderUid string
		if err := rows.Scan(&orderUid); err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		order, err := o.GetOrder(orderUid)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return nil, nil
}

func (o OrderRepository) GetOrder(id string) (*models.Order, error) {
	var order models.Order
	var deliveryJSON, paymentJSON, itemsJSON []byte

	err := o.db.QueryRow(o.ctx,
		"SELECT order_uid, track_number, entry, locale, internal_signature, "+
			"customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, "+
			"delivery, payment, items FROM orders WHERE order_uid = $1",
		id).Scan(
		&order.OrderUid,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.DateCreated,
		&order.OofShard,
		&deliveryJSON,
		&paymentJSON,
		&itemsJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, suberrors.ErrIdOrderNotFound
		}
		return nil, fmt.Errorf("error reading order: %w", err)
	}

	if err := json.Unmarshal(deliveryJSON, &order.Delivery); err != nil {
		return nil, fmt.Errorf("error unmarshaling delivery: %w", err)
	}
	if err := json.Unmarshal(paymentJSON, &order.Payment); err != nil {
		return nil, fmt.Errorf("error unmarshaling payment: %w", err)
	}
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("error unmarshaling items: %w", err)
	}

	return &order, nil
}

func (o OrderRepository) CreateOrder(order *models.Order) (string, error) {
	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		return "", fmt.Errorf("error marshaling delivery: %w", err)
	}

	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		return "", fmt.Errorf("error marshaling payment: %w", err)
	}

	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return "", fmt.Errorf("error marshaling items: %w", err)
	}

	_, err = o.db.Exec(o.ctx,
		"INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, "+
			"customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, "+
			"delivery, payment, items) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shardkey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
		deliveryJSON,
		paymentJSON,
		itemsJSON,
	)
	if err != nil {
		return "", fmt.Errorf("error creating order: %w", err)
	}
	return order.OrderUid, nil
}
