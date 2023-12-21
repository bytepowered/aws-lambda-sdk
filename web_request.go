package lambda

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

func GetHeaderAuthorization(request *events.APIGatewayV2HTTPRequest) string {
	bearer := GetKV(request.Headers, "authorization", "")
	if bearer == "" {
		return bearer
	}
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func GetHeaderParam(request *events.APIGatewayV2HTTPRequest, name string, defaultValue string) string {
	return GetKV(request.Headers, name, defaultValue)
}

func GetQueryParam(request *events.APIGatewayV2HTTPRequest, name string, defaultValue string) string {
	return GetKV(request.QueryStringParameters, name, defaultValue)
}

func GetPathParam(request *events.APIGatewayV2HTTPRequest, name string, defaultValue string) string {
	return GetKV(request.PathParameters, name, defaultValue)
}

func GetKV(kvMap map[string]string, name string, defaultValue string) string {
	if kvMap == nil || len(kvMap) == 0 {
		return defaultValue
	}
	lowerName := strings.ToLower(name)
	for k, v := range kvMap {
		if lowerName == strings.ToLower(k) {
			return v
		}
	}
	return defaultValue
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
