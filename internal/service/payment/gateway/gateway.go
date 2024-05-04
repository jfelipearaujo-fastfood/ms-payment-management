package gateway

import (
	"context"
	"log/slog"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Handle(ctx context.Context, request CreatePaymentGatewayDTO) error {
	if err := request.Validate(); err != nil {
		return err
	}

	// TODO: This is a mock for now and will be replaced by a real call to the gateway API in the future
	slog.InfoContext(ctx, "payment request sent to gateway", "payment_id", request.PaymentID, "amount", request.Amount)

	return nil
}
