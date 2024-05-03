package get_by_order_id

import (
	"context"

	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/repository"
)

type Service struct {
	repository repository.PaymentRepository
}

func NewService(repository repository.PaymentRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Handle(ctx context.Context, request GetByOrderIdDTO) ([]payment_entity.Payment, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	payments, err := s.repository.GetByOrderID(ctx, request.OrderId)
	if err != nil {
		return nil, err
	}

	return payments, nil
}
