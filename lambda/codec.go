package lambda

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
)

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
