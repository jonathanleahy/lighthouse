package sns

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/logger"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/message"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"runtime/debug"
	"time"
)

type (
	SNSSender struct {
		sns   *sns.SNS
		topic string
	}

	Event struct {
		Domain        string            `json:"domain"`
		Type          string            `json:"event_type"`
		CID           string            `json:"cid"`
		Org           string            `json:"org_id"`
		Timestamp     string            `json:"timestamp"`
		Data          interface{}       `json:"data"`
		SchemaVersion int               `json:"schema_version"`
		MessageID     string            `json:"message_id,omitempty"`
		CustomHeaders map[string]string `json:"custom_headers,omitempty"`
	}
)

const (
	schemaVersion = 1
)

func NewAwsSNSClient(topic string, sess *session.Session) *SNSSender {
	return &SNSSender{
		sns:   sns.New(sess),
		topic: topic,
	}
}

func (s *SNSSender) SendEvent(ctx context.Context, domain string, eventType string, body interface{}) (err error) {
	start := time.Now()
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)

	sp, _ := tracer.GenerateChildSpanWithCtx(ctx, "send_sns_event_service")
	attributes := []attribute.KeyValue{
		attribute.String("Cid", rctx.Cid),
		attribute.String("OrgId", rctx.Tenant),
		attribute.String("AccountId", rctx.AccountID),
		attribute.String("sns.topic", s.topic),
		attribute.String("sns.domain", domain),
		attribute.String("sns.event_type", eventType),
	}

	defer func() {
		if err == nil {
			sp.SetStatus(codes.Ok, "success")
		} else {
			sp.SetStatus(codes.Error, err.Error())
			attributes = append(attributes, attribute.String("error.stacktrace", string(debug.Stack())))
		}

		sp.SetAttributes(attributes...)
		sp.End()
	}()

	now := time.Now().UTC().Format(time.RFC3339Nano)

	event := Event{
		Domain:        domain,
		Type:          eventType,
		CID:           rctx.Cid,
		Org:           rctx.Tenant,
		Timestamp:     now,
		Data:          body,
		SchemaVersion: schemaVersion,
		CustomHeaders: rctx.CustomHeaders,
	}

	b, err := json.Marshal(&event)
	if err != nil {
		s.logError(time.Since(start), rctx.Cid, rctx.Tenant, s.topic, domain, eventType, err)
		return err
	}

	bStr := string(b)
	attributes = append(attributes, attribute.String("sns.event", bStr))

	_, err = s.sns.Publish(&sns.PublishInput{
		Message:  &bStr,
		TopicArn: &s.topic,
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"domain": {
				DataType:    aws.String("String"),
				StringValue: aws.String(domain),
			},
			"type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(eventType),
			},
		},
	})

	if err != nil {
		s.logError(time.Since(start), rctx.Cid, rctx.Tenant, s.topic, domain, eventType, err)
		return err
	}

	return nil
}

func (s *SNSSender) logError(t time.Duration, cid string, orgId string, topic string, domain string, eventType string, err error) {
	logger.Error(message.SNSExecErrorMessage, cid, orgId, logger.Fields{
		"sns.duration_exec_ms": t.Milliseconds(),
		"error":                err.Error(),
		"sns.topic":            topic,
		"sns.domain":           domain,
		"sns.event_type":       eventType,
		"stacktrace":           string(debug.Stack()),
	})
}

