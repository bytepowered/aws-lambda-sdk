package lambda

import (
    "context"
    "github.com/aws/aws-lambda-go/events"
    "strings"
)

type HandleFunc func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error)

type MiddlewareFunc func(handler HandleFunc) HandleFunc

func HeaderAuthorization(request *events.APIGatewayV2HTTPRequest) string {
    bearer, ok := request.Headers["authorization"]
    if !ok {
        return ""
    }
    if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
        return bearer[7:]
    }
    return ""
}
