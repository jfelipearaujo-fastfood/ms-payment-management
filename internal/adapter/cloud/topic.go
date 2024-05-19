package cloud

import (
	"context"
)

type TopicService interface {
	GetTopicName() string
	UpdateTopicArn(ctx context.Context) error
	PublishMessage(ctx context.Context, message interface{}) (*string, error)
}

type TopicNotification struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}
