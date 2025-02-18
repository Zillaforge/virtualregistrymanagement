package common

import "context"

type Data struct {
	Action   string            `json:"action"`
	Metadata map[string]string `json:"metadata"`
	Request  interface{}       `json:"request"`
	Response interface{}       `json:"response"`
}

type StartConsumerInput struct {
	Routers map[string]func(ctx context.Context, conmessage interface{})
}
