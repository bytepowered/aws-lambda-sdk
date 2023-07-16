package lambda

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

func SendJSON(object any, statusCode int) (*events.APIGatewayV2HTTPResponse, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("response json, marshal error, %w", err)
	}
	return &events.APIGatewayV2HTTPResponse{
		StatusCode:      statusCode,
		Body:            base64.StdEncoding.EncodeToString(bytes),
		IsBase64Encoded: true,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func SendOK(data any) (*events.APIGatewayV2HTTPResponse, error) {
	return SendJSON(data, 200)
}

func SendERR(error string, statusCode int) (*events.APIGatewayV2HTTPResponse, error) {
	return SendJSON(map[string]string{"error": error}, statusCode)
}

func SendInvalidToken() (*events.APIGatewayV2HTTPResponse, error) {
	return SendERR("Verify: invalid token", 401)
}

func SendInvalidArgs(name string) (*events.APIGatewayV2HTTPResponse, error) {
	return SendERR("Verify: invalid argument: "+name, 400)
}
