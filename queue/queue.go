package queue

import (
	"context"
	"io"

	"github.com/segmentio/kafka-go"
)

//Producer consumes in and produces to the topic in config
func Producer(in chan []byte, config *ProducerConfig) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  config.Brokers,
		Topic:    config.Topic,
		Balancer: &kafka.Hash{},
	})

	defer w.Close()

	for v := range in {
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   nil,
				Value: v,
			},
		)
		if err != nil && config.Errs != nil {
			config.Errs <- err
		}
	}
}

// Consumer consumes a topic and puts the messages into out
func Consumer(ctx context.Context, out chan []byte, config *ConsumerConfig) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Brokers,
		GroupID: config.GroupID,
		Topic:   config.Topic,
	})

	defer r.Close()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF {
				close(out)
				return
			}
			if config.Errs != nil {
				config.Errs <- err
			}
		} else {
			out <- m.Value
		}
	}
}
