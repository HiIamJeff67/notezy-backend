package core

import (
	"context"
	"encoding/json"
	"fmt"

	validatorpkg "github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	emaileventscontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
)

type EmailRequestConsumer struct {
	sender      SenderInterface
	validator   *validatorpkg.Validate
	kafkaConfig platformkafka.ConsumerConfig
}

func NewEmailRequestConsumer(
	sender SenderInterface,
	validator *validatorpkg.Validate,
	kafkaConfig platformkafka.ConsumerConfig,
) *EmailRequestConsumer {
	if validator == nil {
		validator = validatorpkg.New()
	}
	return &EmailRequestConsumer{sender: sender, validator: validator, kafkaConfig: kafkaConfig}
}

func (c *EmailRequestConsumer) Start(ctx context.Context) func() {
	consumer, err := platformkafka.NewConsumer(
		c.kafkaConfig,
		emaileventscontract.CoreEmailRequestTopic.String(),
	)
	if err != nil {
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to create Core email request consumer")
		}
		return func() {}
	}

	workerCtx, cancel := context.WithCancel(ctx)
	go func() {
		if err := consumer.Run(workerCtx, c.consume); err != nil && workerCtx.Err() == nil && logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(workerCtx, err, "Core email request consumer stopped")
		}
	}()

	return func() {
		cancel()
		consumer.Close()
	}
}

func (c *EmailRequestConsumer) consume(
	ctx context.Context,
	_ platformkafka.ConsumerRecord,
	event eventcontract.EventEnvelope[json.RawMessage],
) error {
	if event.EventType != emaileventscontract.EventType_EmailRequested ||
		event.AggregateType != emaileventscontract.AggregateType_EmailRequest {
		return nil
	}

	var metadata struct {
		RequestId uuid.UUID `json:"requestId"`
		Operation string    `json:"operation"`
	}
	if err := json.Unmarshal(event.Data, &metadata); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("decode Core email request: %w", err),
		}
	}
	if metadata.RequestId == uuid.Nil || metadata.RequestId != event.AggregateId {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("Core email request ID does not match the aggregate ID"),
		}
	}

	var exception error
	switch metadata.Operation {
	case emailcontract.SendWelcomeEmailOperation:
		var request emaileventscontract.SendWelcomeEmailRequestDto
		if err := json.Unmarshal(event.Data, &request); err != nil {
			return invalidEmailRequest(err.Error())
		}
		if request.RequestId != event.AggregateId || request.Operation != metadata.Operation {
			return invalidEmailRequest("welcome request metadata is invalid")
		}
		if err := c.validator.Struct(&request); err != nil {
			return invalidEmailRequest(err.Error())
		}
		exception = c.sender.SendWelcomeEmail(ctx, request)
	case emailcontract.SendValidationEmailOperation:
		var request emaileventscontract.SendValidationEmailRequestDto
		if err := json.Unmarshal(event.Data, &request); err != nil {
			return invalidEmailRequest(err.Error())
		}
		if request.RequestId != event.AggregateId || request.Operation != metadata.Operation {
			return invalidEmailRequest("validation request metadata is invalid")
		}
		if err := c.validator.Struct(&request); err != nil {
			return invalidEmailRequest(err.Error())
		}
		exception = c.sender.SendValidationEmail(ctx, request)
	case emailcontract.SendSecurityAlertEmailOperation:
		var request emaileventscontract.SendSecurityAlertEmailRequestDto
		if err := json.Unmarshal(event.Data, &request); err != nil {
			return invalidEmailRequest(err.Error())
		}
		if request.RequestId != event.AggregateId || request.Operation != metadata.Operation {
			return invalidEmailRequest("security alert request metadata is invalid")
		}
		if err := c.validator.Struct(&request); err != nil {
			return invalidEmailRequest(err.Error())
		}
		exception = c.sender.SendSecurityAlertEmail(ctx, request)
	default:
		return invalidEmailRequest("unsupported email operation")
	}
	if exception != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_Transient,
			Origin:         exception,
		}
	}

	return nil
}

func invalidEmailRequest(message string) error {
	return &platformkafka.ConsumerError{
		Classification: platformkafka.ErrorClassification_SchemaIncompatible,
		Origin:         fmt.Errorf("invalid Core email request: %s", message),
	}
}
