package cloud

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/create"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/gateway"
)

type CtxKey string

const MessageId CtxKey = "message_id"

type QueueService interface {
	GetQueueName() string
	UpdateQueueUrl(ctx context.Context) error
	ConsumeMessages(ctx context.Context)
}

type AwsSqsService struct {
	queueName string
	queueUrl  string
	client    *sqs.Client

	createPayment        service.CreatePaymentService[create.CreatePaymentDTO]
	createPaymentGateway service.CreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO]

	chanMessage chan types.Message

	mutex     sync.Mutex
	waitGroup sync.WaitGroup
}

func NewQueueService(
	queueName string,
	config aws.Config,
	createPayment service.CreatePaymentService[create.CreatePaymentDTO],
	createPaymentGateway service.CreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO],
) QueueService {
	client := sqs.NewFromConfig(config)

	return &AwsSqsService{
		queueName: queueName,
		client:    client,

		createPayment:        createPayment,
		createPaymentGateway: createPaymentGateway,

		chanMessage: make(chan types.Message, 10),

		mutex:     sync.Mutex{},
		waitGroup: sync.WaitGroup{},
	}
}

func (s *AwsSqsService) GetQueueName() string {
	return s.queueName
}

func (s *AwsSqsService) UpdateQueueUrl(ctx context.Context) error {
	output, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &s.queueName,
	})
	if err != nil {
		return err
	}

	s.queueUrl = *output.QueueUrl

	return nil
}

func (s *AwsSqsService) ConsumeMessages(ctx context.Context) {
	output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &s.queueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		slog.ErrorContext(ctx, "error receiving message from queue", "queue_url", s.queueUrl, "error", err)
		return
	}

	s.waitGroup.Add(len(output.Messages))

	for _, message := range output.Messages {
		go s.processMessage(ctx, message)
	}

	s.waitGroup.Wait()
}

func (s *AwsSqsService) processMessage(ctx context.Context, message types.Message) {
	defer s.waitGroup.Done()
	s.mutex.Lock()

	ctx = context.WithValue(ctx, MessageId, *message.MessageId)

	slog.InfoContext(ctx, "message received")

	var notification TopicNotification

	err := json.Unmarshal([]byte(*message.Body), &notification)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshalling message", "error", err)
	} else {
		if notification.Type != "Notification" {
			slog.ErrorContext(ctx, "invalid notification type", "type", notification.Type)
		} else {
			var request create.CreatePaymentDTO

			err = json.Unmarshal([]byte(notification.Message), &request)
			if err != nil {
				slog.ErrorContext(ctx, "error unmarshalling message", "error", err)
			}

			if err == nil {
				slog.InfoContext(ctx, "message unmarshalled", "request", request)
				payment, err := s.createPayment.Handle(ctx, request)
				if err != nil {
					slog.ErrorContext(ctx, "error create payment", "error", err)
				}

				if payment != nil {
					gatewayReq := gateway.CreatePaymentGatewayDTO{
						PaymentID: payment.PaymentId,
						Amount:    payment.Amount,
					}

					if err := s.createPaymentGateway.Handle(ctx, gatewayReq); err != nil {
						slog.ErrorContext(ctx, "error create payment gateway", "error", err)
					}
				}
			}
		}
	}

	if err := s.deleteMessage(ctx, message); err != nil {
		slog.ErrorContext(ctx, "error deleting message", "error", err)
	}

	s.mutex.Unlock()
}

func (s *AwsSqsService) deleteMessage(ctx context.Context, message types.Message) error {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &s.queueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return err
	}

	return nil
}
