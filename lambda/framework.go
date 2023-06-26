package lambda

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/aws/aws-lambda-go/events"
    "io"
    "net/http"
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
