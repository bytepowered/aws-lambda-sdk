package lambda

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type HandleFunc func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error)

type MiddlewareFunc func(handler HandleFunc) HandleFunc

func HeaderAuthorization(request *events.APIGatewayV2HTTPRequest) string {
	bearer := HeaderLookup(request.Headers, "authorization", "")
	if bearer == "" {
		return bearer
	}
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func HeaderLookup(headers map[string]string, name string, defaults string) string {
	if headers == nil || len(headers) == 0 {
		return defaults
	}
	for k, v := range headers {
		if strings.ToLower(name) == strings.ToLower(k) {
			return v
		}
	}
	return defaults
}
