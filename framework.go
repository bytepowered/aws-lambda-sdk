package lambda

import (
    "context"
    "github.com/aws/aws-lambda-go/events"
)

type HandleFunc func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error)

type MiddlewareFunc func(handler HandleFunc) HandleFunc
