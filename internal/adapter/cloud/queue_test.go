package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/service/payment/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetQueueName(t *testing.T) {
	t.Run("Should return queue name", func(t *testing.T) {
		// Arrange
		fakeProcessor := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)

		service := NewQueueService("test-queue", aws.Config{}, fakeProcessor)

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

		fakeProcessor := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

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

		fakeProcessor := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

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
			"order_id": "c3fdab1b-3c06-4db2-9edc-4760a2429460",
			"payment_id": "9dfa1386-2f52-4cca-b9aa-f9bd6887d447",
			"total_items": 1,
			"amount": 100.0
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("123"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
					{
						MessageId:     aws.String("456"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
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

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567891"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)

		fakeProcessor.On("Handle", ctx, mock.Anything).
			Return(nil, nil).
			Times(2)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
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
				WaitTimeSeconds:     5,
			},
			Error: raiseErr,
		})

		fakeProcessor := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
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
			"order_id": "c3fdab1b-3c06-4db2-9edc-4760a2429460",
			"payment_id": "9dfa1386-2f52-4cca-b9aa-f9bd6887d447",
			"total_items": "abc",
			"amount": 100.0
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
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

		fakeProcessor := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
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
			"order_id": "c3fdab1b-3c06-4db2-9edc-4760a2429460",
			"payment_id": "9dfa1386-2f52-4cca-b9aa-f9bd6887d447",
			"total_items": 1,
			"amount": 100.0
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
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

		fakeProcessor := mocks.NewMockCreatePaymentService[create.CreatePaymentDTO](t)

		fakeProcessor.On("Handle", ctx, mock.Anything).
			Return(nil, assert.AnError).
			Once()

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})
}
