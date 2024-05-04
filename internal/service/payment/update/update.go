package update

import (
	"context"

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

func (s *Service) Handle(ctx context.Context, request UpdatePaymentDTO) (*payment_entity.Payment, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	payment, err := s.repository.GetByID(ctx, request.PaymentId)
	if err != nil {
		return nil, err
	}

	validStates := []payment_entity.PaymentState{
		payment_entity.Approved,
		payment_entity.Rejected,
	}

	if payment.IsInState(validStates...) {
		return nil, custom_error.ErrPaymentAlreadyInState
	}

	var state payment_entity.PaymentState

	if request.Approved {
		state = payment_entity.Approved
	} else {
		state = payment_entity.Rejected
	}

	payment.UpdateState(state, s.timeProvider.GetTime())

	if err := s.repository.Update(ctx, &payment); err != nil {
		return nil, err
	}

	return &payment, nil
}
