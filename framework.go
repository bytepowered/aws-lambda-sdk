package lambda

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	aws "github.com/aws/aws-lambda-go/lambda"
)

type HandleFunc func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error)

type MiddlewareFunc func(handler HandleFunc) HandleFunc

func Start(handler HandleFunc) {
	aws.Start(handler)
}

func StartCORS(handler HandleFunc) {
	aws.StartWithOptions(UseCORS(handler))
}

func StartWithJWT(handler HandleFunc) {
	aws.StartWithOptions(JWTFilter(handler))
}

func StartWithCORSJWT(handler HandleFunc) {
	aws.StartWithOptions(UseCORS(JWTFilter(handler)))
}
