package create

import (
	"context"
	"log/slog"

	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type Service struct {
	repository   repository.PaymentRepository
	timeProvider provider.TimeProvider
}

func NewService(
	repository repository.PaymentRepository,
	timeProvider provider.TimeProvider,
) *Service {
	return &Service{
		repository:   repository,
		timeProvider: timeProvider,
	}
}

func (s *Service) Handle(ctx context.Context, request CreatePaymentDTO) (*payment_entity.Payment, error) {
	if err := request.Validate(ctx); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "checking if payment already exists", "payment_id", request.PaymentId)

	exists, err := s.repository.GetByID(ctx, request.PaymentId)
	if err != nil && err != custom_error.ErrPaymentNotFound {
		return nil, err
	}

	if exists.Exists() {
		slog.ErrorContext(ctx, "payment already exists", "payment_id", request.PaymentId)
		return nil, custom_error.ErrPaymentAlreadyExists
	}

	slog.InfoContext(ctx, "payment not found, creating new payment", "payment_id", request.PaymentId, "order_id", request.OrderId)

	items := make([]payment_entity.PaymentItem, len(request.Items))
	for i, item := range request.Items {
		items[i] = payment_entity.NewPaymentItem(item.Id, item.Name, item.Quantity)
	}

	payment := payment_entity.NewPayment(
		request.OrderId,
		request.PaymentId,
		items,
		request.TotalItems,
		request.Amount,
		s.timeProvider.GetTime(),
	)

	if err := s.repository.Create(ctx, &payment); err != nil {
		slog.ErrorContext(ctx, "error creating payment", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "payment created", "payment_id", payment.PaymentId)

	return &payment, nil
}
