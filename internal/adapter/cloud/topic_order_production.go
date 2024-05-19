package cloud

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type OrderProductionTopicService struct {
	TopicName string
	TopicArn  string
	Client    *sns.Client
}

func NewOrderProductionTopicService(topicName string, config aws.Config) TopicService {
	client := sns.NewFromConfig(config)

	return &OrderProductionTopicService{
		TopicName: topicName,
		Client:    client,
	}
}

func (s *OrderProductionTopicService) GetTopicName() string {
	return s.TopicName
}

func (s *OrderProductionTopicService) UpdateTopicArn(ctx context.Context) error {
	output, err := s.Client.ListTopics(ctx, &sns.ListTopicsInput{})
	if err != nil {
		return err
	}

	for _, topic := range output.Topics {
		if strings.Contains(*topic.TopicArn, s.TopicName) {
			s.TopicArn = *topic.TopicArn
			return nil
		}
	}

	return custom_error.ErrTopicNotFound
}

func (s *OrderProductionTopicService) PublishMessage(ctx context.Context, message interface{}) (*string, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req := &sns.PublishInput{
		TopicArn: aws.String(s.TopicArn),
		Message:  aws.String(string(body)),
	}

	out, err := s.Client.Publish(ctx, req)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "message published", "topic", s.TopicName, "message_id", *out.MessageId, "message", string(body))

	return out.MessageId, nil
}
