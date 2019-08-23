package kafka_status

import (
	"context"
	"encoding/json"
	"github.com/ypapax/jsn"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check/pipes"
	"github.com/ypapax/status_check/status"
)

type kafkaStatusPipe struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaStatusPipe(writer *kafka.Writer, reader *kafka.Reader) pipes.StatusPipe {
	return &kafkaStatusPipe{
		writer: writer,
		reader: reader,
	}
}

func (ksp kafkaStatusPipe) Publish(parent context.Context, st status.Status) error {
	b, err := json.Marshal(st)
	if err != nil {
		logrus.Error(err)
		return err
	}
	message := kafka.Message{
		Value: b,
		Time:  time.Now(),
	}
	if err := ksp.writer.WriteMessages(parent, message); err != nil {
		logrus.Error(err)
		return err
	}
	return err
}

func (ksp kafkaStatusPipe) Listen(ctx context.Context, statusChan chan<- status.Status, errs chan<- error) {
	if ksp.reader == nil {
		panic("reader is nil")
	}
	for {
		logrus.Tracef("before reading a message: %+v", jsn.B(ksp.reader.Config()))
		m, err := ksp.reader.ReadMessage(ctx)
		if err != nil {
			logrus.Error(err)
			errs <- err
			continue
		}
		logrus.Tracef("partition: %+v, message at offset %d: %s = %s\n", m.Partition, m.Offset, string(m.Key), string(m.Value))
		var st status.Status
		if err := json.Unmarshal(m.Value, &st); err != nil {
			logrus.Error(err)
			errs <- err
			continue
		}
		statusChan <- st
	}
}
