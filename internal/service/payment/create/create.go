package create

import (
	"context"

	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/repository"
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
	if err := request.Validate(); err != nil {
		return nil, err
	}

	payment := payment_entity.NewPayment(request.OrderId,
		request.PaymentId,
		request.TotalItems,
		request.Amount,
		s.timeProvider.GetTime(),
	)

	if err := s.repository.Create(ctx, &payment); err != nil {
		return nil, err
	}

	return &payment, nil
}
