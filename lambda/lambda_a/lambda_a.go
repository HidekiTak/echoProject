package main

import (
	"context"
	"encoding/json"

	lmd "aws/lambda"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	h := handler{}
	lambda.Start(h.HandleRequest)
}

type handler struct {
}

type Parameter struct {
}

func s3EventConverter(bucket string, key string) (*Parameter, error) {
	return nil, nil
}
func bytesConverter(bytes []byte) (*Parameter, error) {
	return nil, nil
}

func (h *handler) HandleRequest(ctx context.Context, event json.RawMessage) error {
	converter := lmd.NewEventConverter[Parameter](s3EventConverter, bytesConverter)
	param, err := converter.Convert(event)
	if err != nil {
		return err
	}
	println(json.Marshal(param))
	return nil
}
