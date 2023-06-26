package lambda

import (
    "errors"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
    "time"
)

const (
    JwtEnvSecret = "AUTH_JWT_SECRET"
    JwtEnvIssuer = "AUTH_JWT_ISSUER"
    JwtEnvExpsec = "AUTH_JWT_EXPSEC"
)

// JWTSign 生成JWT字符串
func JWTSign(sub string, id string) (string, error) {
    expsec := RequiredIntEnv(JwtEnvExpsec, 7*24*3600)
    now := time.Now()
    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        ID:        id,
        Issuer:    RequiredEnv(JwtEnvIssuer),
        Subject:   sub,
        ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(expsec))),
        IssuedAt:  jwt.NewNumericDate(now),
    }).SignedString([]byte(RequiredEnv(JwtEnvSecret)))
    if err != nil {
        return "", fmt.Errorf("JWT-SIGN-TOKEN: %w", err)
    }
    return token, nil
}

// JWTParse 解析JWT字符串
func JWTParse(tokenStr string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(RequiredEnv(JwtEnvSecret)), nil
    }, jwt.WithIssuer(RequiredEnv(JwtEnvIssuer)))
    if token.Valid {
        return token, nil
    }
    if errors.Is(err, jwt.ErrTokenMalformed) {
        return nil, fmt.Errorf("JWT-MALFORMED: %w", err)
    } else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
        return nil, fmt.Errorf("JWT-SIGN-INVALID: %w", err)
    } else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
        return nil, fmt.Errorf("JWT-EXPIRED: %w", err)
    } else {
        return nil, fmt.Errorf("JWT-UNKNOWN: %w", err)
    }
}
