package server

import (
	"github.com/jfelipearaujo-org/ms-payment-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/create"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/update"
)

type Dependency struct {
	TimeProvider *time_provider.TimeProvider

	PaymentRepository repository.PaymentRepository

	CreatePaymentService service.CreatePaymentService[create.CreatePaymentDTO]

	UpdatePaymentService service.UpdatePaymentService[update.UpdatePaymentDTO]
}
