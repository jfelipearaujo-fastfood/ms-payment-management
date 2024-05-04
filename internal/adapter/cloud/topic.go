package cloud

import (
	"context"
)

type TopicService interface {
	GetTopicName() string
	UpdateTopicArn(ctx context.Context) error
	PublishMessage(ctx context.Context, message interface{}) (*string, error)
}
