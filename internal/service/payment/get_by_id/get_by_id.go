package get_by_id

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

func (s *Service) Handle(ctx context.Context, request GetByIdDTO) (payment_entity.Payment, error) {
	if err := request.Validate(); err != nil {
		return payment_entity.Payment{}, err
	}

	payment, err := s.repository.GetByID(ctx, request.PaymentId)
	if err != nil {
		return payment_entity.Payment{}, err
	}

	return payment, nil
}
