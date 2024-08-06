package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/adapter/cloud"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/adapter/database"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/environment"
	get_by_order_id_handler "github.com/jfelipearaujo-org/ms-payment-management/internal/handler/get_by_order_id"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/handler/health"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/handler/payment_hook"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/repository/payment"
	token "github.com/jfelipearaujo-org/ms-payment-management/internal/server/middlewares"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/create"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/gateway"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/get_by_id"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/get_by_order_id"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/update"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config          *environment.Config
	DatabaseService database.DatabaseService
	QueueService    cloud.QueueService

	UpdateOrderTopicService     cloud.TopicService
	OrderProductionTopicService cloud.TopicService

	Dependency Dependency
}

func NewServer(config *environment.Config) *Server {
	ctx := context.Background()

	cloudConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	if config.CloudConfig.IsBaseEndpointSet() {
		cloudConfig.BaseEndpoint = aws.String(config.CloudConfig.BaseEndpoint)
	}

	databaseService := database.NewDatabase(config)

	timeProvider := time_provider.NewTimeProvider(time.Now)
	paymentRepository := payment.NewPaymentRepository(databaseService.GetInstance())
	createPaymentService := create.NewService(paymentRepository, timeProvider)
	createPaymentGatewayService := gateway.NewService()

	updateOrderTopicService := cloud.NewUpdateOrderTopicService(config.CloudConfig.UpdateOrderTopic, cloudConfig)
	orderProductionTopicService := cloud.NewOrderProductionTopicService(config.CloudConfig.OrderProductionTopic, cloudConfig)

	return &Server{
		Config:          config,
		DatabaseService: databaseService,
		QueueService: cloud.NewQueueService(
			config.CloudConfig.OrderPaymentQueue,
			cloudConfig,
			createPaymentService,
			createPaymentGatewayService,
		),

		UpdateOrderTopicService:     updateOrderTopicService,
		OrderProductionTopicService: orderProductionTopicService,

		Dependency: Dependency{
			TimeProvider: timeProvider,

			PaymentRepository: paymentRepository,

			CreatePaymentService: createPaymentService,
			UpdatePaymentService: update.NewService(paymentRepository, timeProvider),

			UpdateOrderTopicService:     updateOrderTopicService,
			OrderProductionTopicService: orderProductionTopicService,

			GetPaymentByOrderIdService: get_by_order_id.NewService(paymentRepository),
			GetPaymentByIDService:      get_by_id.NewService(paymentRepository),
		},
	}
}

func (s *Server) GetHttpServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.ApiConfig.Port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(logger.Middleware())
	e.Use(middleware.Recover())

	s.registerHealthCheck(e)

	group := e.Group(fmt.Sprintf("/api/%s", s.Config.ApiConfig.ApiVersion))

	s.registerPaymentHandlers(group)

	return e
}

func (server *Server) registerHealthCheck(e *echo.Echo) {
	healthHandler := health.NewHandler(server.DatabaseService)

	e.GET("/health", healthHandler.Handle)
}

func (s *Server) registerPaymentHandlers(e *echo.Group) {
	updatePaymentHandler := payment_hook.NewHandler(
		s.Dependency.GetPaymentByIDService,
		s.Dependency.UpdatePaymentService,
		s.Dependency.OrderProductionTopicService,
		s.Dependency.UpdateOrderTopicService,
	)

	getPaymentByOrderIdHandler := get_by_order_id_handler.NewHandler(s.Dependency.GetPaymentByOrderIdService)

	e.Use(token.Middleware())
	e.PATCH("/payments/webhook/:payment_id", updatePaymentHandler.Handle)
	e.GET("/payments/order/:order_id", getPaymentByOrderIdHandler.Handle)
}
