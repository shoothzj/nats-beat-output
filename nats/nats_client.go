package nats

import (
	"context"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"github.com/nats-io/nats.go"
)

type natClient struct {
	config *natConfig

	info  beat.Info
	codec codec.Codec
	stats outputs.Observer
	conn  *nats.Conn
}

func (n *natClient) String() string {
	return "websocket"
}

func (n *natClient) Connect() error {
	var err error
	n.conn, err = nats.Connect(n.config.Url)
	if err != nil {
		return err
	}
	n.codec, err = codec.CreateEncoder(n.info, n.config.Codec)
	return err
}

func (n *natClient) Close() error {
	n.conn.Close()
	return nil
}

func (n *natClient) Publish(_ context.Context, batch publisher.Batch) error {
	events := batch.Events()
	// record this batch
	n.stats.NewBatch(len(events))
	failEvents, err := n.PublishEvents(events)
	if err != nil {
		// send success ack
		batch.ACK()
	} else {
		// send fail retry. Limited to RetryLimit
		batch.RetryEvents(failEvents)
	}
	return err
}

func (n *natClient) PublishEvents(events []publisher.Event) ([]publisher.Event, error) {
	for i, event := range events {
		err := n.publishEvent(&event)
		if err != nil {
			// if send one failure, mark the rest msg to retry
			return events[i:], err
		}
	}
	return nil, nil
}

func (n *natClient) publishEvent(event *publisher.Event) error {
	serializedEvent, err := n.codec.Encode(n.info.Beat, &event.Content)
	if err != nil {
		return err
	}
	buf := make([]byte, len(serializedEvent))
	copy(buf, serializedEvent)
	return n.conn.Publish(n.config.Subj, buf)
}
