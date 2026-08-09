package core

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	validatorpkg "github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	emaileventscontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
)

type senderStub struct {
	err error
}

func (s senderStub) SendWelcomeEmail(context.Context, emaileventscontract.SendWelcomeEmailRequestDto) error {
	return s.err
}

func (s senderStub) SendValidationEmail(context.Context, emaileventscontract.SendValidationEmailRequestDto) error {
	return s.err
}

func (s senderStub) SendSecurityAlertEmail(context.Context, emaileventscontract.SendSecurityAlertEmailRequestDto) error {
	return s.err
}

func TestEmailRequestConsumerMapsLocalErrorClassification(t *testing.T) {
	cases := []struct {
		name      string
		retryable bool
		wantClass platformkafka.ErrorClassification
	}{
		{
			name:      "retryable delivery error",
			retryable: true,
			wantClass: platformkafka.ErrorClassification_Transient,
		},
		{
			name:      "non retryable configuration error",
			retryable: false,
			wantClass: platformkafka.ErrorClassification_PoisonMessage,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			requestId := uuid.New()
			request := emaileventscontract.SendWelcomeEmailRequestDto{
				RequestId:  requestId,
				Operation:  emailcontract.SendWelcomeEmailOperation,
				OccurredAt: time.Now().UTC(),
				To:         "user@example.com",
				UserName:   "Notezy User",
				Status:     "active",
			}
			data, err := json.Marshal(request)
			if err != nil {
				t.Fatalf("marshal request: %v", err)
			}

			stubException := exceptions.New("DeliveryFailed", "Email", "SendEmail", "Failed to deliver the email", 502)
			stubException.Retryable = test.retryable
			consumer := &EmailRequestConsumer{
				sender: senderStub{
					err: stubException,
				},
			}
			consumer.validator = validatorpkg.New()
			resultErr := consumer.consume(
				context.Background(),
				platformkafka.ConsumerRecord{},
				eventcontract.EventEnvelope[json.RawMessage]{
					SchemaVersion: eventcontract.Version,
					EventType:     emaileventscontract.EventType_EmailRequested,
					AggregateType: emaileventscontract.AggregateType_EmailRequest,
					AggregateId:   requestId,
					Data:          data,
				},
			)

			consumerError, ok := resultErr.(*platformkafka.ConsumerError)
			if !ok {
				t.Fatalf("error type = %T, want *platformkafka.ConsumerError", resultErr)
			}
			if consumerError.Classification != test.wantClass {
				t.Fatalf("classification = %q, want %q", consumerError.Classification, test.wantClass)
			}
		})
	}
}
