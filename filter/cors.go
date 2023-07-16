package filter

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bytepowered/aws-lambda-sdk"
)

func CORSAllowHeaders(req events.APIGatewayV2HTTPRequest) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      lambda.HeaderLookup(req.Headers, "access-control-allow-methods", "*"),
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     lambda.HeaderLookup(req.Headers, "access-control-allow-methods", "GET,POST,PUT,DELETE,OPTIONS"),
		"Access-Control-Allow-Headers":     lambda.HeaderLookup(req.Headers, "access-control-allow-origin", "*"),
	}
}

func CORSFilter(next lambda.HandleFunc) lambda.HandleFunc {
	return UseCORS(next)
}

func UseCORS(next lambda.HandleFunc) lambda.HandleFunc {
	return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
		resp, err := next(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp.Headers == nil {
			resp.Headers = map[string]string{}
		}
		for k, v := range CORSAllowHeaders(req) {
			resp.Headers[k] = v
		}
		return resp, nil
	}
}
