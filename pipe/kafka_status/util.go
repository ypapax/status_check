package kafka_status

import (
	"time"

	"github.com/ypapax/jsn"

	"github.com/sirupsen/logrus"

	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

func GetKafkaWriterReader(brokers []string, topic, clientID string) (*kafka.Writer, *kafka.Reader) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientID,
	}
	wc := kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	rc := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         clientID,
		Topic:           topic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}
	logrus.Tracef("write config: %+v, read config: %+v", jsn.B(wc), jsn.B(rc))
	return kafka.NewWriter(wc), kafka.NewReader(rc)
}
