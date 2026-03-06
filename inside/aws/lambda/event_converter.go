package lambda

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

type EventConverter[T any] interface {
	Convert(message json.RawMessage) (*T, error)
}

// NewEventConverter Create Lambda EventMessage Converter
func NewEventConverter[T any](
	// s3Event -> T Converter
	s3Event func(bucket string, key string) (*T, error),
	// []byte -> T Converter
	bytes func(message []byte) (*T, error),
) EventConverter[T] {
	var converters []EventConverter[T]
	if s3Event != nil {
		converters = append(converters, &s3EventConverter[T]{convert: s3Event})
	}
	if bytes != nil {
		converters = append(converters, &bytesConverter[T]{convert: bytes})
	}
	return &eventConverter[T]{
		converters: converters,
	}
}

type eventConverter[T any] struct {
	converters []EventConverter[T]
}

func (e eventConverter[T]) Convert(message json.RawMessage) (*T, error) {
	for _, converter := range e.converters {
		t, _ := converter.Convert(message)
		if t != nil {
			return t, nil
		}
	}
	return nil, errors.New("no converter found")
}

type s3EventConverter[T any] struct {
	convert func(bucket string, key string) (*T, error)
}

func (conv *s3EventConverter[T]) Convert(message json.RawMessage) (*T, error) {
	var s3event events.S3Event
	err := json.Unmarshal(message, &s3event)
	if err != nil {
		return nil, err
	}
	if len(s3event.Records) != 1 {
		return nil, errors.New("record must be 1")
	}
	record := s3event.Records[0].S3
	return conv.convert(record.Bucket.Name, record.Object.Key)
}

type bytesConverter[T any] struct {
	convert func(bytes []byte) (*T, error)
}

func (conv *bytesConverter[T]) Convert(message json.RawMessage) (*T, error) {
	return conv.convert(message)
}
