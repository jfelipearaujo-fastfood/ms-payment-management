package repository

import (
	"context"

	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *payment_entity.Payment) error
	GetByID(ctx context.Context, paymentId string) (payment_entity.Payment, error)
	GetByOrderID(ctx context.Context, orderId string) ([]payment_entity.Payment, error)
	Update(ctx context.Context, payment *payment_entity.Payment) error
}
