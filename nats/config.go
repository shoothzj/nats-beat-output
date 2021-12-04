package nats

import "github.com/elastic/beats/v7/libbeat/outputs/codec"

type natConfig struct {
	// Workers Number of worker goroutines publishing log events
	Workers int `config:"workers" validate:"min=1"` // Number of worker goroutines publishing log events
	// BatchSize Max number of events in a batch to send to a single client
	BatchSize int `config:"batch_size" validate:"min=1"`
	// RetryLimit Max number of retries for single batch of events
	RetryLimit int `config:"retry_limit"`
	// Url Nat url
	Url string `config:"url"`
	// Subj Nat subj
	Subj string `config:"subj"`

	Codec codec.Config `config:"codec"`
}
