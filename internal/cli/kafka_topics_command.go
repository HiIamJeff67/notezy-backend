package main

import (
	"fmt"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	kafkatopics "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka/topics"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(newKafkaCommand())
}

func newKafkaCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "kafka",
		Short: "Manage Kafka development resources.",
	}
	command.AddCommand(newEnsureKafkaTopicsCommand())
	return command
}

func newEnsureKafkaTopicsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "topics ensure",
		Short: "Create all versioned Notezy Kafka topics and their dead-letter topics.",
		RunE: func(command *cobra.Command, _ []string) error {
			connectionConfig, err := platformkafka.LoadConnectionConfig()
			if err != nil {
				return err
			}

			provisioner, err := platformkafka.NewTopicProvisioner(platformkafka.ClientConfig{
				ConnectionConfig: connectionConfig,
				ClientId:         "notezy-kafka-topic-bootstrap",
			})
			if err != nil {
				return err
			}
			defer provisioner.Close()

			if err := provisioner.EnsureTopics(command.Context(), kafkatopics.All()); err != nil {
				return fmt.Errorf("ensure Notezy Kafka topics: %w", err)
			}

			return nil
		},
	}
}
