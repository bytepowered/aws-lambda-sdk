package lambda

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io"
	"net/http"
)

func CheckResponse(statusCode int, body io.ReadCloser) (int, error) {
	if !(200 <= statusCode && statusCode < 300) {
		data, _ := io.ReadAll(body)
		return statusCode, fmt.Errorf("http response status not ok, status: %d, body: %s", statusCode, string(data))
	}
	return statusCode, nil
}

func SerializeResponse(statusCode int, body io.ReadCloser, outptr any) (int, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("http response read error: %w", err)
	}
	err = json.Unmarshal(data, outptr)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("http response unmarshal error: %w, data: %s", err, string(data))
	}
	return statusCode, nil
}

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
	return SendERR("ERR-INVALID-TOKEN", 401)
}

func SendInvalidArgs(name string) (*events.APIGatewayV2HTTPResponse, error) {
	return SendERR("ERR-INVALID-ARG: "+name, 400)
}
