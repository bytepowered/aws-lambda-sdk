package lambda

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

func CORSAllowHeaders(req events.APIGatewayV2HTTPRequest) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      HeaderLookup(req.Headers, "access-control-allow-methods", "*"),
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     HeaderLookup(req.Headers, "access-control-allow-methods", "GET,POST,PUT,DELETE,OPTIONS"),
		"Access-Control-Allow-Headers":     HeaderLookup(req.Headers, "access-control-allow-origin", "*"),
	}
}

func UseCORS(next HandleFunc) HandleFunc {
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
