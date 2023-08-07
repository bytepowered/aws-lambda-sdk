package lambda

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

const (
	JwtEnvSecret = "AUTH_JWT_SECRET"
	JwtEnvIssuer = "AUTH_JWT_ISSUER"
	JwtEnvExpsec = "AUTH_JWT_EXPSEC"
)

const (
	ctxKeyJwtClaims = "@internal.aws.lambda.sdk.jwt.claims"
)

// JWTSigned 生成JWT字符串
func JWTSigned(sub string, id string) (string, error) {
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

func JWTFilter(next HandleFunc) HandleFunc {
	return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
		str := HeaderAuthorization(&req)
		if !CheckNotEmpty(str) {
			log.Println("ERROR: jwt token not found in header")
			return SendInvalidToken()
		}
		token, err := JWTParse(str)
		if err != nil {
			log.Println("ERROR: jwt token invalid, error:", err, ", token:", str)
			return SendInvalidToken()
		}
		return next(context.WithValue(ctx, ctxKeyJwtClaims, token.Claims), req)
	}
}

func JWTLoadClaims(ctx context.Context) (claims jwt.Claims, ok bool) {
	value := ctx.Value(ctxKeyJwtClaims)
	if value == nil {
		return nil, false
	}
	claims, ok = value.(jwt.Claims)
	return
}

func JWTLoadSubject(ctx context.Context) (sub string, ok bool) {
	claims, ok := JWTLoadClaims(ctx)
	if !ok {
		return "", false
	}

	if v, err := claims.GetSubject(); err != nil {
		return "", false
	} else {
		return v, true
	}
}
