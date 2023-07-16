package lambda

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

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

func ParseBody(data string, isBase64 bool, outptr any) (err error) {
	var bytes []byte
	if isBase64 {
		bytes, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			return fmt.Errorf("decode base64 %w", err)
		}
	} else {
		bytes = []byte(data)
	}
	return json.Unmarshal(bytes, outptr)
}
