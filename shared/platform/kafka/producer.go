package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	franzkgo "github.com/twmb/franz-go/pkg/kgo"
	plain "github.com/twmb/franz-go/pkg/sasl/plain"
	scram "github.com/twmb/franz-go/pkg/sasl/scram"
)

type Producer struct {
	client *franzkgo.Client
}

func NewProducer(kafkaConfig ClientConfig) (*Producer, error) {
	options, err := newConnectionOptions(kafkaConfig)
	if err != nil {
		return nil, err
	}

	client, err := franzkgo.NewClient(options...)
	if err != nil {
		return nil, fmt.Errorf("create Kafka client: %w", err)
	}

	return &Producer{
		client: client,
	}, nil
}

func newConnectionOptions(kafkaConfig ClientConfig) ([]franzkgo.Opt, error) {
	if len(kafkaConfig.Brokers) == 0 {
		return nil, errors.New("Kafka brokers are required")
	}

	options := []franzkgo.Opt{
		franzkgo.SeedBrokers(kafkaConfig.Brokers...),
		franzkgo.ClientID(kafkaConfig.ClientId),
		franzkgo.DialTimeout(kafkaConfig.DialTimeout),
	}
	if kafkaConfig.TLS.Enabled {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: kafkaConfig.TLS.ServerName,
		}
		if kafkaConfig.TLS.CAFile != "" {
			certificateAuthority, err := os.ReadFile(kafkaConfig.TLS.CAFile)
			if err != nil {
				return nil, fmt.Errorf("read Kafka TLS certificate authority: %w", err)
			}
			certificateAuthorities := x509.NewCertPool()
			if !certificateAuthorities.AppendCertsFromPEM(certificateAuthority) {
				return nil, errors.New("parse Kafka TLS certificate authority")
			}
			tlsConfig.RootCAs = certificateAuthorities
		}
		if kafkaConfig.TLS.CertFile != "" || kafkaConfig.TLS.KeyFile != "" {
			if kafkaConfig.TLS.CertFile == "" || kafkaConfig.TLS.KeyFile == "" {
				return nil, errors.New("Kafka TLS certificate and key must be configured together")
			}
			certificate, err := tls.LoadX509KeyPair(
				kafkaConfig.TLS.CertFile,
				kafkaConfig.TLS.KeyFile,
			)
			if err != nil {
				return nil, fmt.Errorf("load Kafka TLS client certificate: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{
				certificate,
			}
		}
		options = append(options, franzkgo.DialTLSConfig(tlsConfig))
	}
	if kafkaConfig.SASL.Mechanism != "" {
		if kafkaConfig.SASL.Username == "" || kafkaConfig.SASL.Password == "" {
			return nil, errors.New("Kafka SASL username and password are required")
		}

		switch kafkaConfig.SASL.Mechanism {
		case "PLAIN":
			options = append(options, franzkgo.SASL(plain.Auth{
				User: kafkaConfig.SASL.Username,
				Pass: kafkaConfig.SASL.Password,
			}.AsMechanism()))
		case "SCRAM-SHA-256":
			options = append(options, franzkgo.SASL(scram.Auth{
				User: kafkaConfig.SASL.Username,
				Pass: kafkaConfig.SASL.Password,
			}.AsSha256Mechanism()))
		case "SCRAM-SHA-512":
			options = append(options, franzkgo.SASL(scram.Auth{
				User: kafkaConfig.SASL.Username,
				Pass: kafkaConfig.SASL.Password,
			}.AsSha512Mechanism()))
		default:
			return nil, fmt.Errorf("unsupported Kafka SASL mechanism %q", kafkaConfig.SASL.Mechanism)
		}
	}

	return options, nil
}

func (p *Producer) Ping(ctx context.Context) error {
	if p == nil || p.client == nil {
		return errors.New("Kafka producer is unavailable")
	}

	startedAt := time.Now()
	err := p.client.Ping(ctx)
	RecordBrokerPing(ctx, time.Since(startedAt), err)

	return err
}

func (p *Producer) Produce(
	ctx context.Context,
	topic string,
	key string,
	value []byte,
) error {
	if p == nil || p.client == nil {
		return errors.New("Kafka producer is unavailable")
	}

	startedAt := time.Now()
	result := p.client.ProduceSync(ctx, &franzkgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
	err := result.FirstErr()
	RecordPublish(ctx, topic, time.Since(startedAt), err)

	return err
}

func (p *Producer) Close() {
	if p == nil || p.client == nil {
		return
	}

	p.client.Close()
}
