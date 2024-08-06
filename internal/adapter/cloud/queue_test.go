package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/create"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/gateway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetQueueName(t *testing.T) {
	t.Run("Should return queue name", func(t *testing.T) {
		// Arrange
		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		service := NewQueueService("test-queue", aws.Config{}, createPayment, createPaymentGateway)

		// Act
		queueName := service.GetQueueName()

		// Assert
		assert.Equal(t, "test-queue", queueName)
	})
}

func TestUpdateQueueUrl(t *testing.T) {
	t.Run("Should return nil when queue is found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		// Act
		err := service.UpdateQueueUrl(ctx)

		// Assert
		assert.NoError(t, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Should return error when GetQueueUrl operation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		raiseErr := &testtools.StubError{Err: errors.New("ClientError")}

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Error:         raiseErr,
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		// Act
		err := service.UpdateQueueUrl(ctx)

		// Assert
		testtools.VerifyError(err, raiseErr, t)
		testtools.ExitTest(stubber, t)
	})
}

func TestStartConsuming(t *testing.T) {
	t.Run("Should start consuming messages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"be6293ff-4ec0-4ed8-95c9-b36ce99aa105\",\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		createPayment.On("Handle", mock.Anything, mock.Anything).
			Return(&payment_entity.Payment{}, nil).
			Once()

		createPaymentGateway.On("Handle", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		createPayment.AssertExpectations(t)
		createPaymentGateway.AssertExpectations(t)
	})

	t.Run("Should do nothing when receive message return an error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		raiseErr := &testtools.StubError{Err: errors.New("ClientError")}

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Error: raiseErr,
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		createPayment.AssertExpectations(t)
		createPaymentGateway.AssertExpectations(t)
	})

	t.Run("Should log when cannot unmarshal message", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : false,
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":123,\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		createPayment.AssertExpectations(t)
	})

	t.Run("Should log when message processor returns an error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"be6293ff-4ec0-4ed8-95c9-b36ce99aa105\",\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("123"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		createPayment.On("Handle", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).
			Once()

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		createPayment.AssertExpectations(t)
		createPaymentGateway.AssertExpectations(t)
	})

	t.Run("Should log when payment gateway returns an error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"be6293ff-4ec0-4ed8-95c9-b36ce99aa105\",\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("123"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		createPayment.On("Handle", mock.Anything, mock.Anything).
			Return(&payment_entity.Payment{}, nil).
			Once()

		createPaymentGateway.On("Handle", mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		createPayment.AssertExpectations(t)
		createPaymentGateway.AssertExpectations(t)
	})

	t.Run("Should not consume the message if the type is not Notification", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "AnotherType",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"be6293ff-4ec0-4ed8-95c9-b36ce99aa105\",\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567891"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		createPayment.AssertExpectations(t)
		createPaymentGateway.AssertExpectations(t)
	})

	t.Run("Should not consume the message if cannot unmarshal the message", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\": false,\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567891"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		createPayment := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)
		createPaymentGateway := mocks.NewMockCreatePaymentGatewayService[gateway.CreatePaymentGatewayDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, createPayment, createPaymentGateway)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		createPayment.AssertExpectations(t)
		createPaymentGateway.AssertExpectations(t)
	})
}
